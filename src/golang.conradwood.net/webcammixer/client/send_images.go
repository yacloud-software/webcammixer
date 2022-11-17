package main

import (
	"fmt"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"io"
	"time"
)

func sendImages() error {
	s := &ImageSender{}
	for {
		err := utils.DirWalk(*send_images, func(root string, rel string) error {
			return s.sendImage(root + rel)
		})
		if err != nil {
			return err
		}
	}
}

type ImageSender struct {
	srv    pb.WebCamMixer_SendImagesClient
	toggle bool
}

func (s *ImageSender) sendImage(filename string) error {
	if s.srv == nil {
		ctx := authremote.ContextWithTimeout(time.Duration(10000) * time.Hour)
		srv, err := pb.GetWebCamMixerClient().SendImages(ctx)
		if err != nil {
			return err
		}
		s.srv = srv
	}
	buf, err := utils.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Printf("Sending image: %s (%v) (%d bytes)\n", filename, s.toggle, len(buf))
	offset := 0
	repeat := true
	for repeat {
		l := 8192
		if l+offset > len(buf) {
			l = len(buf) - offset
			repeat = false
		}
		data := buf[offset : offset+l]
		ist := &pb.ImageStream{
			NextImage: s.toggle,
			Data:      data,
		}
		err = s.srv.Send(ist)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		offset = offset + l
	}

	s.toggle = !s.toggle
	time.Sleep(*delay)
	return nil
}
