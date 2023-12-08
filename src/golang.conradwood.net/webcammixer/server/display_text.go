package main

import (
	"context"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/webcammixer/loopback"
	"strings"
	"time"
)

var (
	pre_userframe_state *display_state
	display_until       time.Time
)

func init() {
	go restore_state_loop()
}

type display_state struct {
	userframe bool
	provider  loopback.FrameProvider
}

func restore_state_loop() {
	for {
		time.Sleep(time.Duration(1) * time.Second)
		if pre_userframe_state == nil {
			continue
		}
		if time.Now().Before(display_until) {
			continue
		}
		err := restore_state(pre_userframe_state)
		if err != nil {
			fmt.Printf("failed to restore state: %s\n", err)
		} else {
			pre_userframe_state = nil
		}
	}
}
func restore_state(ds *display_state) error {
	fmt.Printf("Restoring state to %#v\n", ds)
	return nil
}
func (e *echoServer) DisplayText(ctx context.Context, req *pb.DisplayTextRequest) (*common.Void, error) {
	if !loopback.IsWatched() {
		fmt.Printf("Displaytext ignored, nobody watching...\n")
		return &common.Void{}, nil
	}
	if req.MaxSeconds == 0 {
		req.MaxSeconds = 1
	}

	// is the userframeprovider active atm?
	if !switcher_impl.IsUserFrameActive() {
		// make up a config now
		ds := &display_state{
			userframe: false,
			provider:  mixapp.GetCurrentProvider(),
		}
		ifp := switcher_impl.GetCurrentUserImageProvider()
		err := ifp.SetConfig(config_from_request(req))
		if err != nil {
			return nil, err
		}
		switcher_impl.ActivateUserFrames()
		pre_userframe_state = ds
		display_until = time.Now().Add(time.Duration(req.MaxSeconds) * time.Second)
	}

	return &common.Void{}, nil
}

func config_from_request(req *pb.DisplayTextRequest) *pb.UserImageRequest {
	uir := &pb.UserImageRequest{
		Converters: []*pb.UserImageConverter{
			&pb.UserImageConverter{Type: pb.ConverterType_LABEL, Text: req.Text},
		},
	}
	if req.ContinueCamera {
		w := mixapp.GetCurrentProvider().GetID()
		fmt.Printf("Continuing camera (from current \"%s\")...\n", w)
		if strings.HasPrefix(w, "webcam://") {
			dev := "webcam://"
			cd := &pb.CaptureDevice{Device: strings.TrimPrefix(w, dev)}
			uic := &pb.UserImageConverter{Type: pb.ConverterType_WEBCAM, Device: cd}
			uir.Converters = append([]*pb.UserImageConverter{uic}, uir.Converters...)
			fmt.Printf("Added device \"%s\"\n", dev)
		}
	}
	return uir
}
