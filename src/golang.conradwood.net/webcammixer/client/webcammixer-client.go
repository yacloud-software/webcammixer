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
	app         = flag.Bool("app", false, "start app")
	repeat      = flag.Bool("repeat", false, "if true, repeat video selection")
	idle_text   = flag.String("text", "", "set idle text")
	delay       = flag.Duration("delay", time.Duration(500)*time.Millisecond, "`delay` between images")
	send_images = flag.String("send_images", "", "send all pix in this `directory`")
	send_frames = flag.String("send_frames", "", "send all frames in this `directory`")
	videocam    = flag.String("videodev", "", "if set, connect loopback to this `/dev/videoX`")
	stopvideo   = flag.Bool("idle", false, "switch to idle source")
	start_app   = flag.Bool("start_app", false, "start app in userspace")
	cam         = flag.String("camera", "", "if non nil uses images server to stream cameras, e.g. cam://espcam1")
	echoClient  pb.WebCamMixerClient
)

func main() {
	flag.Parse()
	var err error
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
		err = mixerapp.Start()
		utils.Bail("failed to start", err)
		mixerapp.WaitUntilDone()
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
	d, found := devs[devnum]
	if !found {
		fmt.Printf("Device does not exist (%d)\n", devnum)
		goto repeat
	}

	term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Printf("Setting device %s\n", d.Device)
	ctx := authremote.Context()
	vr := &pb.VideoDeviceDef{VideoDeviceName: d.Device}
	_, err = pb.GetWebCamMixerClient().SendVideoDevice(ctx, vr)
	if err != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
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
