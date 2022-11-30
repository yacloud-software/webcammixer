package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/mixerapp"
	"os"
	"time"
)

var (
	idle_text   = flag.String("text", "", "set idle text")
	delay       = flag.Duration("delay", time.Duration(500)*time.Millisecond, "`delay` between images")
	send_images = flag.String("send_images", "", "send all pix in this `directory`")
	send_frames = flag.String("send_frames", "", "send all frames in this `directory`")
	videocam    = flag.String("videodev", "", "if set, connect loopback to this `/dev/videoX`")
	stopvideo   = flag.Bool("idle", false, "switch to idle source")
	echoClient  pb.WebCamMixerClient
)

func main() {
	flag.Parse()
	var err error
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
	err = mixerapp.Start()
	utils.Bail("failed to start", err)
	mixerapp.WaitUntilDone()
	os.Exit(0)
	echoClient = pb.GetWebCamMixerClient()

	// a context with authentication
	//	ctx := authremote.Context()

	//	empty := &common.Void{}
	//	response, err := echoClient.Ping(ctx, empty)
	//	utils.Bail("Failed to ping server", err)
	//	fmt.Printf("Response to ping: %v\n", response)
end:
	fmt.Printf("Done.\n")
	os.Exit(0)
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
