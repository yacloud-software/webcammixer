package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/mixerapp"
	"golang.org/x/term"
	"os"
	"strconv"
	"time"
)

var (
	app           = flag.Bool("app", false, "start app")
	repeat        = flag.Bool("repeat", false, "if true, repeat video selection")
	idle_text     = flag.String("text", "", "(L1) set idle text")
	countdown     = flag.Duration("countdown", 0, "set a countdown")
	delay         = flag.Duration("delay", time.Duration(500)*time.Millisecond, "`delay` between images")
	send_images   = flag.String("send_images", "", "(L1) send all pix in this `directory`")
	send_frames   = flag.String("send_frames", "", "(L1) send all frames in this `directory`")
	videocam      = flag.String("videodev", "", "if set, connect loopback to this `/dev/videoX`")
	stopvideo     = flag.Bool("idle", false, "switch to idle source")
	start_app     = flag.Bool("start_app", false, "start app in userspace")
	cam           = flag.String("camera", "", "if non nil uses images server to stream cameras, e.g. cam://espcam1")
	echoClient    pb.WebCamMixerClient
	dyntext       = flag.String("dyntext", "", "(L2) if set, set a dynamic text")
	status        = flag.Bool("status", false, "if true get status")
	overlay_text  = flag.String("overlay_text", "", "(L3) display text for a few seconds. also see -overlay")
	overlay_img   = flag.String("overlay_img", "", "(L3) display image for a few seconds. also see -overlay")
	overlay_emoji = flag.String("overlay_emoji", "", "(L3) display emoji")
	none          = flag.String("help", "", "NOTE: (Ln) denotes the level of abstraction. use the highest possible for your usecase")
)

func main() {
	flag.Parse()
	var err error
	if *overlay_text != "" || *overlay_img != "" || *overlay_emoji != "" {
		utils.Bail("failed to overlay", Overlay())
		goto end
	}
	if *status {
		utils.Bail("failed to do status", Status())
		goto end
	}
	if *dyntext != "" {
		utils.Bail("failed to set dyntext", DynText())
		goto end
	}
	if *countdown != 0 {
		utils.Bail("Failed to set camera", SetCountdown())
		goto end
	}

	if *cam != "" {
		utils.Bail("Failed to set camera", SetCam())
		goto end
	}
	if *idle_text != "" {
		utils.Bail("failed to set text", SetText())
		goto end
	}
	if *stopvideo {
		utils.Bail("failed to stop video", StopVideoCam())
		goto end
	}
	if *videocam != "" {
		utils.Bail("failed to set videocam", SetVideoCam())
		goto end
	}
	if *send_images != "" {
		utils.Bail("failed to send images", sendImages())
		goto end
	}
	if *send_frames != "" {
		utils.Bail("failed to send frames", sendFrames())
		goto end
	}
	if *app {
		ma := &mixerapp.MixerApp{}
		err = ma.Start()
		utils.Bail("failed to start", err)
		ma.WaitUntilDone()
		os.Exit(0)
	}
	utils.Bail("failed to get capture devices", detect())
end:
	fmt.Printf("Done.\n")
	os.Exit(0)
}
func detect() error {
	echoClient = pb.GetWebCamMixerClient()
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	utils.Bail("failed to set term", err)
	term.Restore(int(os.Stdin.Fd()), oldState)
repeat_with_print:
	wl, err := echoClient.GetCaptureDevices(authremote.Context(), &common.Void{})
	if err != nil {
		return err
	}
	devs := make(map[int]*pb.CaptureDevice)
	t := &utils.Table{}
	t.AddHeaders("#", "Name", "Device")
	for i, d := range wl.Devices {
		devs[i+1] = d
		t.AddInt(i + 1)
		t.AddString(d.Name)
		t.AddString(d.Device)
		t.NewRow()
	}
	fmt.Println("0 - idle text")
	fmt.Println(t.ToPrettyString())
	_, err = term.MakeRaw(int(os.Stdin.Fd()))
	utils.Bail("failed to set term", err)
repeat:
	_, err = term.MakeRaw(int(os.Stdin.Fd()))
	utils.Bail("failed to set term", err)
	fmt.Printf("Press number on keyboard to select video device -> ")
	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
		return err
	}
	if b[0] == 27 || b[0] == 'q' {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Printf("Aborted.\n")
		return nil
	}
	devnum, err := strconv.Atoi(string(b))
	if err != nil {
		fmt.Printf("Not a valid number: \"%s\"\n", err)
		goto repeat
	}
	if devnum == 0 {
		term.Restore(int(os.Stdin.Fd()), oldState)
		ctx := authremote.Context()

		_, err = pb.GetWebCamMixerClient().SwitchToIdle(ctx, &common.Void{})
		if err != nil {
			return err
		}

		_, err = pb.GetWebCamMixerClient().SetIdleText(ctx, &pb.IdleTextRequest{Text: "standby..."})
		if err != nil {
			return err
		}
		goto repeat_with_print
	}
	d, found := devs[devnum]
	if !found {
		fmt.Printf("Device does not exist (%d)\n", devnum)
		goto repeat
	}

	term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Printf("Setting device %s\n", d.Device)
	ctx := authremote.Context()
	//	vr := &pb.VideoDeviceDef{VideoDeviceName: d.Device}
	//	_, err = pb.GetWebCamMixerClient().SendVideoDevice(ctx, vr)
	_, err = pb.GetWebCamMixerClient().SendFromCaptureDevice(ctx, d)
	if err != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Printf("Failed to set device: %s\n", d.Device)
		return err
	}
	goto repeat_with_print

}
func SetVideoCam() error {
	ctx := authremote.Context()
	vr := &pb.VideoDeviceDef{VideoDeviceName: *videocam}
	_, err := pb.GetWebCamMixerClient().SendVideoDevice(ctx, vr)
	return err
}
func StopVideoCam() error {
	ctx := authremote.Context()
	_, err := pb.GetWebCamMixerClient().SwitchToIdle(ctx, &common.Void{})
	return err
}
func SetText() error {
	ctx := authremote.Context()
	_, err := pb.GetWebCamMixerClient().SetIdleText(ctx, &pb.IdleTextRequest{Text: *idle_text})
	return err
}

func SetCam() error {
	ctx := authremote.Context()
	ur := &pb.URL{URL: *cam}
	_, err := pb.GetWebCamMixerClient().SwitchToLiveImages(ctx, ur)
	return err
}

func SetCountdown() error {
	ctx := authremote.Context()
	dur := *countdown
	req := &pb.CountdownRequest{
		Text:    "Countdown",
		Seconds: uint32(dur.Seconds()),
	}
	_, err := pb.GetWebCamMixerClient().SetCountdown(ctx, req)
	fmt.Printf("Started countdown of %d seconds\n", req.Seconds)
	return err
}

func DynText() error {
	ctx := authremote.Context()
	uir := &pb.UserImageRequest{
		Converters: []*pb.UserImageConverter{
			&pb.UserImageConverter{Type: pb.ConverterType_LABEL, Text: *dyntext},
		},
	}
	if *videocam != "" {
		cd := &pb.CaptureDevice{Device: *videocam}
		uic := &pb.UserImageConverter{Type: pb.ConverterType_WEBCAM, Device: cd}
		uir.Converters = append([]*pb.UserImageConverter{uic}, uir.Converters...)
	}
	_, err := pb.GetWebCamMixerClient().SetUserImage(ctx, uir)
	if err != nil {
		return err
	}
	return nil
}
func Status() error {
	ctx := authremote.Context()
	fp, err := pb.GetWebCamMixerClient().GetCurrentProvider(ctx, &common.Void{})
	if err != nil {
		return err
	}
	fmt.Printf("Current Provider: \"%s\"\n", fp.HumanReadableDesc)
	return nil
}

func Overlay() error {
	ctx := authremote.Context()

	uir := &pb.UserImageRequest{
		ImageSource: &pb.ImageSource{FillColour: &pb.Colour{Red: 0, Green: 0, Blue: 255}},
	}
	if *videocam != "" {
		cd := &pb.CaptureDevice{Device: *videocam}
		uir.ImageSource = &pb.ImageSource{Device: cd}
	}

	if *overlay_text != "" {
		uir.Converters = append(uir.Converters, &pb.UserImageConverter{Reference: "text", Type: pb.ConverterType_LABEL, Text: *overlay_text})
	}
	if *overlay_emoji != "" {
		str := *overlay_emoji
		uir.Converters = append(uir.Converters, &pb.UserImageConverter{Reference: "emoji", Type: pb.ConverterType_EMOJI, Emoji: str})
	}
	if *overlay_img != "" {
		b, err := utils.ReadFile(*overlay_img)
		if err != nil {
			return err
		}
		uir.Converters = append(uir.Converters, &pb.UserImageConverter{Reference: "image", Type: pb.ConverterType_OVERLAY_IMAGE, OverlayImage: &pb.OverlayImageRequest{XPos: 100, YPos: 100, Image: b}})
	}

	_, err := pb.GetWebCamMixerClient().SetUserImage(ctx, uir)
	if err != nil {
		return err
	}
	return nil

}
