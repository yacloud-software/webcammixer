package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/defaults"
	"golang.conradwood.net/webcammixer/mixerapp"
	"golang.org/x/image/draw"
	"google.golang.org/grpc"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"time"
)

var (
	debug      = flag.Bool("debug", false, "debug mode")
	pr         = utils.ProgressReporter{}
	port       = flag.Int("port", 4190, "The grpc server port")
	frame_chan = make(chan []byte, 5)
)

type echoServer struct {
}

func main() {
	var err error
	flag.Parse()
	fmt.Printf("Starting WebCamMixerServer...\n")
	go cache_webcam_devices()
	sd := server.NewServerDef()
	sd.SetPort(*port)
	sd.SetRegister(server.Register(
		func(server *grpc.Server) error {
			e := new(echoServer)
			pb.RegisterWebCamMixerServer(server, e)
			return nil
		},
	))
	//	go test()
	go func() {
		utils.Bail("failed to start app", mixerapp.Start())
	}()
	//go frame_worker()
	err = server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)
}

/************************************
* grpc functions
************************************/
func (e *echoServer) SendFromCaptureDevice(ctx context.Context, req *pb.CaptureDevice) (*common.Void, error) {
	if req.Type == 0 {
		vdd := &pb.VideoDeviceDef{
			VideoDeviceName: req.Device,
		}
		return e.SendVideoDevice(ctx, vdd)
	}
	if req.Type == 1 {
		return e.SwitchToLiveImages(ctx, &pb.URL{URL: req.Device})
	}
	return nil, fmt.Errorf("cannot yet set type %d (%s)", req.Type, req.Device)

}
func (e *echoServer) SendVideoDevice(ctx context.Context, req *pb.VideoDeviceDef) (*common.Void, error) {
	loopdev := mixerapp.GetLoopDev()
	if loopdev == nil {
		return nil, fmt.Errorf("not ready yet - try again later")
	}
	fmt.Printf("Sending from device %s\n", req.VideoDeviceName)
	if req.VideoDeviceName == loopdev.DeviceName {
		return nil, errors.InvalidArgs(ctx, "loopback device cannot be source", "loopback device cannot be source")
	}
	h, w := defaults.GetDimensions()
	sav, err := SourceActivateVideoDef(req.VideoDeviceName, h, w) // starts a thread reading from this video device
	if err != nil {
		return nil, err
	}
	loopdev.SetProvider(sav)
	loopdev.SetTimingSource(sav)
	return &common.Void{}, nil
}
func (e *echoServer) SwitchToIdle(ctx context.Context, req *common.Void) (*common.Void, error) {
	mfp := mixerapp.NewManualFrameProvider()
	loopdev := mixerapp.GetLoopDev()
	if loopdev == nil {
		return nil, fmt.Errorf("not ready yet - try again later")
	}
	loopdev.SetProvider(mfp)
	loopdev.SetTimingSource(mfp)
	return &common.Void{}, nil
}
func (e *echoServer) SwitchToLiveImages(ctx context.Context, req *pb.URL) (*common.Void, error) {
	mfp := NewLiveImageProvider(req.URL, mixerapp.GetLoopDev())

	loopdev := mixerapp.GetLoopDev()
	if loopdev == nil {
		return nil, fmt.Errorf("not ready yet - try again later")
	}
	loopdev.SetProvider(mfp)
	loopdev.SetTimingSource(mfp)
	return &common.Void{}, nil
}
func (e *echoServer) SendFrames(srv pb.WebCamMixer_SendFramesServer) error {
	lastimage := false
	var framedata []byte
	mfp := mixerapp.NewManualFrameProvider()
	loopdev := mixerapp.GetLoopDev()
	if loopdev == nil {
		return fmt.Errorf("not ready yet - try again later")
	}
	loopdev.SetProvider(mfp)
	loopdev.SetTimingSource(mfp)
	for {
		imgdata, err := srv.Recv()
		//		fmt.Printf("Received %d bytes\n", len(imgdata.Data))
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error receiving frames: %s\n", err)
			return err
		}
		if len(framedata) == 0 {
			lastimage = imgdata.NextImage
		}
		if lastimage != imgdata.NextImage {
			mfp.NewFrame(framedata)
			framedata = make([]byte, 0)
			lastimage = imgdata.NextImage
		}
		framedata = append(framedata, imgdata.Data...)
	}

	// once all images were sent, revert back to idle
	/*
		ifp := mixerapp.NewIdleFrameProvider(loopdev.GetDimensions())
		ifp.Run()
		loopdev.SetProvider(ifp)
		loopdev.SetTimingSource(ifp)
	*/
	return nil
}

func (e *echoServer) SendImages(srv pb.WebCamMixer_SendImagesServer) error {
	lastimage := false
	var framedata []byte
	for {
		imgdata, err := srv.Recv()
		//		fmt.Printf("Received %d bytes\n", len(imgdata.Data))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(framedata) == 0 {
			lastimage = imgdata.NextImage
		}
		if lastimage != imgdata.NextImage {
			err := newImage(framedata)
			if err != nil {
				fmt.Printf("image trouble: %s\n", err)
			}
			framedata = framedata[:0]
			lastimage = imgdata.NextImage
		}
		framedata = append(framedata, imgdata.Data...)
	}
	return nil
}

func frame_worker() {
	for {
		framedata := <-frame_chan
		newImage(framedata)
	}
}
func submitNewImage(framedata []byte) {
	frame_chan <- framedata
}
func newImage(framedata []byte) error {

	pr.Add(1)
	r := bytes.NewReader(framedata)
	img, name, err := image.Decode(r)
	if err != nil {
		return err
	}
	loopdev := mixerapp.GetLoopDev()
	h, w := loopdev.GetDimensions()

	scale := false
	var dst image.Image
	if scale {
		// Set the expected size that you want:
		rd := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

		// Resize:
		draw.NearestNeighbor.Scale(rd, rd.Rect, img, img.Bounds(), draw.Over, nil)
		dst = rd
	} else {
		dst = img
	}
	if *debug {
		fmt.Printf("New frame, %d bytes, name:%s,image bounds: %#v\n", len(framedata), name, dst.Bounds())
	}

	rawimage, err := converters.ConvertToRaw(dst)
	if err != nil {
		return err
	}
	//	return loopdev.WriteImage(rawimage.YUV420())

	fmt.Printf("NOT IMPLEMENTED -  streaming images (%d bytes)\n", len(rawimage.YUV422()))
	return nil
}
func test() {
	send := "/tmp/x/scaled_pngs"
	time.Sleep(time.Duration(3) * time.Second)
	for {
		pr.Print()
		err := utils.DirWalk(send, func(root string, rel string) error {
			b, err := utils.ReadFile(root + "/" + rel)
			if err != nil {
				return err
			}
			submitNewImage(b)
			return nil
		})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
func (e *echoServer) GetLoopbackInfo(ctx context.Context, req *common.Void) (*pb.LoopbackInfo, error) {
	loopdev := mixerapp.GetLoopDev()
	if loopdev == nil {
		return nil, fmt.Errorf("not ready yet - try again later")
	}
	res := &pb.LoopbackInfo{}
	res.Height, res.Width = loopdev.GetDimensions()
	return res, nil
}
func (e *echoServer) SetIdleText(ctx context.Context, req *pb.IdleTextRequest) (*common.Void, error) {
	ifp := mixerapp.DefaultIdleFrameProvider()
	if ifp == nil {
		return nil, fmt.Errorf("not ready yet - try again later")
	}

	ifp.SetIdleText(func() string { return req.Text })
	return &common.Void{}, nil
}

func (e *echoServer) GetCaptureDevices(ctx context.Context, req *common.Void) (*pb.CaptureDeviceList, error) {
	return getCaptureDevices()
}

func (e *echoServer) SetCountdown(ctx context.Context, req *pb.CountdownRequest) (*common.Void, error) {
	ifp := mixerapp.DefaultIdleFrameProvider()
	if ifp == nil {
		return nil, fmt.Errorf("not ready yet - try again later")
	}
	ct := &countdowner{started: time.Now(), duration: time.Duration(req.Seconds) * time.Second}
	ifp.SetIdleText(ct.getText)
	return &common.Void{}, nil
}

type countdowner struct {
	started  time.Time
	duration time.Duration
}

func (c *countdowner) getText() string {
	diff := time.Since(c.started)
	if diff > c.duration {
		return "idle"
	}
	remain := c.duration - diff
	return RenderDuration(remain)
}

func RenderDuration(t time.Duration) string {
	s := uint32(t.Seconds())
	t = time.Duration(s) * time.Second
	return fmt.Sprintf("%v", t)
}
