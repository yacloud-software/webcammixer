package main

import (
	//	"context"
	"fmt"
	//	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/images"
	//	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/loopback"
	"io"
	"time"
)

const (
	live_img_idle = time.Duration(120) * time.Second
)

type liveimgprovider struct {
	url                   string
	time_target           chan bool
	last_frame            []byte
	last_time_had_channel time.Time
	loopdev               *loopback.LoopBackDevice
}

func NewLiveImageProvider(url string, loopdev *loopback.LoopBackDevice) *liveimgprovider {
	res := &liveimgprovider{url: url, loopdev: loopdev}
	go res.worker()
	return res
}
func (l *liveimgprovider) GetFrame() ([]byte, error) {
	return l.last_frame, nil
}
func (l *liveimgprovider) GetID() string {
	return "livimg"
}
func (l *liveimgprovider) SetTimerTarget(c chan bool) error {
	l.time_target = c
	return nil
}
func (l *liveimgprovider) worker() {
	ctx := authremote.ContextWithTimeout(time.Duration(2) * time.Hour)
	url := &images.URL{URL: l.url}
	srv, err := images.GetImagesClient().ReadFromHTTPStream(ctx, url)
	if err != nil {
		l.Printf("failed to start stream to \"%s\": %s\n", l.url, err)
		return
	}
	h, w := l.loopdev.GetDimensions()
	for {
		img, err := srv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.Printf("Failed to recv image: %s\n", err)
			break
		}
		//l.Printf("New image, scaling to %dx%d\n", h, w)
		b, err := converters.ConvertImageDataToFrame(img.Data, w, h)
		if err != nil {
			l.Printf("failed to convert image %s\n", err)
			break
		}
		l.last_frame = b
		c := l.time_target
		if c != nil {
			l.last_time_had_channel = time.Now()
			c <- true
		} else {
			l.Printf("got image, but no channel...\n")
			if time.Since(l.last_time_had_channel) > live_img_idle {
				break
			}

		}
	}
	l.Printf("Exited\n")
}

func (l *liveimgprovider) Printf(format string, args ...interface{}) {
	x := fmt.Sprintf(format, args...)
	fmt.Printf("[liveimgprovider %s] %s", l.url, x)
}
