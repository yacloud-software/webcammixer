package main

import (
	"flag"
	"fmt"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/webcam"
	"golang.org/x/image/draw"
	"image"
	"sync"
	"time"
)

var (
	webcam_idle_timeout = flag.Duration("webcam_idle_timeout", time.Duration(10)*time.Second, "duration after which a webcam /dev/video device will be closed if not in use")
	webcam_lock         sync.Mutex
	webcam_sources      = make(map[string]*VideoCamSource) // key:videodevicename
)

type VideoCamSource struct {
	videoDeviceName string
	lastFrame       *VideoCamFrame
	frameChan       <-chan []byte
	isRunning       bool
	lock            sync.Mutex
	cam             string
	wci             *webcam.WebCamInfo
	threadRunning   bool
	onNewFrame      chan bool // if non-nil send a bool each time a new frame is received. intented as a timing source, not a frame queueing thing. it's expected that the called function triggers a channel, and the new thread calls GetLastImage() for all active sources
	lastFrameUsed   time.Time // last time a frame was actually used (e.g. we had a channel)
	height          uint32    // height of frames we must produce
	width           uint32    // width of frames we must produce
}

type VideoCamFrame struct {
	created time.Time
	format  v4l2.PixFormat
	data    []byte
}

func SourceDectivateAll() error {
	webcam_lock.Lock()
	var ws []*VideoCamSource
	for _, v := range webcam_sources {
		ws = append(ws, v)
	}
	webcam_lock.Unlock()
	for _, v := range ws {
		v.Deactivate()
	}
	return nil
}

// called by loopback device
func (v *VideoCamSource) GetID() string {
	return "webcam-" + v.videoDeviceName
}

// called by loopback device
func (v *VideoCamSource) GetFrame() ([]byte, error) {
	lf := v.lastFrame
	if lf == nil {
		return nil, nil
	}
	return lf.data, nil
}

// NOTE: if we use more than height & width, consider making it a struct
func SourceActivateVideoDef(devicename string, height, width uint32) (*VideoCamSource, error) {
	if !utils.FileExists(devicename) {
		return nil, fmt.Errorf("device \"%s\" does not exist", devicename)
	}
	vs := GetOrCreateWebcamStream(devicename)
	err := vs.Activate(height, width)
	if err != nil {
		return nil, err
	}
	vs.height = height
	vs.width = width
	return vs, nil
}

func GetOrCreateWebcamStream(devicename string) *VideoCamSource {
	webcam_lock.Lock()
	defer webcam_lock.Unlock()
	vs, exists := webcam_sources[devicename]
	if !exists {
		vs = &VideoCamSource{videoDeviceName: devicename}
		webcam_sources[devicename] = vs
	}
	return vs
}
func (v *VideoCamSource) Activate(height, width uint32) error {
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.isRunning {
		fmt.Printf("Device \"%s\" already running\n", v.videoDeviceName)
		return nil
	}
	fmt.Printf("Opening device \"%s\"...\n", v.videoDeviceName)
	wci, err := webcam.Open(v.videoDeviceName, height, width)
	if err != nil {
		v.isRunning = false
		return err
	}
	v.wci = wci
	v.isRunning = true

	if v.frameChan == nil {
		v.frameChan, err = v.wci.Start()
		if err != nil {
			fmt.Printf("Failed to start stream: %s\n", err)
			wci.Close()
			v.isRunning = false
			return err
		}
	}

	go v.readerThread()
	fmt.Printf("Opened device \"%s\"...\n", v.videoDeviceName)
	return nil
}
func (v *VideoCamSource) Deactivate() error {
	v.lock.Lock()
	defer v.lock.Unlock()
	if !v.isRunning {
		return nil
	}
	if v.wci == nil {
		return nil
	}
	//	v.wci.Device.Close()
	//	v.wci = nil
	//	v.isRunning = false
	panic("cannot deactivate yet")
}

func (v *VideoCamSource) GetMostRecentFrame() (*VideoCamFrame, error) {
	return v.lastFrame, nil
}

func (v *VideoCamSource) readerThread() {
	v.threadRunning = true
	defer func() {
		v.threadRunning = false
	}()

	must_convert := false
	pf := v.wci.GetActualPixelFormat()
	if pf.Height != v.height || pf.Width != v.width {
		must_convert = true
		fmt.Printf("Warning - device %s must scale\n", v.videoDeviceName)
	}
	for {
		frame := <-v.frameChan
		vcf := &VideoCamFrame{
			created: time.Now(),
			data:    frame,
			format:  pf,
		}
		if must_convert {
			vcf = v.convert_video_frame(vcf)
		}
		v.lastFrame = vcf
		if len(frame) == 0 {
			time.Sleep(time.Duration(400) * time.Millisecond)
		}
		//fmt.Printf("Got frame from device %s.\n", v.videoDeviceName)
		c := v.onNewFrame
		if c != nil {
			v.lastFrameUsed = time.Now()
			c <- true
		}
		if time.Since(v.lastFrameUsed) >= *webcam_idle_timeout {
			// no longer in use (no listener for the frames...)
			break
		}
	}
	v.wci.Close()
	fmt.Printf("Webcam %s stopped (thread exit)\n", v.videoDeviceName)
	v.threadRunning = false
	v.frameChan = nil
	v.isRunning = false // RACE CONDITION
}
func (v *VideoCamSource) SetTimerTarget(c chan bool) error {
	v.onNewFrame = c
	return nil
}
func (v *VideoCamSource) convert_video_frame(vcf *VideoCamFrame) *VideoCamFrame {
	//	fname := fmt.Sprintf("/tmp/webcam-frame.yuv422_%dx%d", vcf.format.Width, vcf.format.Height)
	//	utils.WriteFile(fname, vcf.data)
	img := converters.ConvertYUV422ToImage(vcf.data, int(vcf.format.Height), int(vcf.format.Width))
	dst := image.NewRGBA(image.Rect(0, 0, int(v.width), int(v.height)))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	rawimage, err := converters.ConvertToRaw(dst)
	if err != nil {
		fmt.Printf("ERROR Converting: %s\n", err)
		return nil
	}
	vcf.data = rawimage.DefaultBytes()
	return vcf
}
