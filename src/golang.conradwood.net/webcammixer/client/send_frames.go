package main

import (
	"bufio"
	"fmt"
	"golang.conradwood.net/apis/common"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/v1/converters"
	"golang.org/x/image/draw"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	WIDTH  = uint32(800)
	HEIGHT = uint32(600)
)

func sendFrames() error {
	ctx := authremote.Context()
	li, err := pb.GetWebCamMixerClient().GetLoopbackInfo(ctx, &common.Void{})
	if err != nil {
		return err
	}
	WIDTH = li.Width
	HEIGHT = li.Height
	s := &FrameSender{}
	for {
		err := utils.DirWalk(*send_frames, func(root string, rel string) error {
			return s.send(root + rel)
		})
		if err != nil {
			fmt.Printf("An error has occured: %s\n", err)
		}
	}
}

type FrameSender struct {
	srv    pb.WebCamMixer_SendFramesClient
	toggle bool
}

func (s *FrameSender) send(filename string) error {
	var err error
	if s.srv == nil {
		ctx := authremote.ContextWithTimeout(time.Duration(10000) * time.Hour)
		s.srv, err = pb.GetWebCamMixerClient().SendFrames(ctx)
		if err != nil {
			fmt.Printf("err:%s\n", err)
			return err
		}
	}
	var buf []byte

	if !strings.HasSuffix(filename, "rgb") {
		buf, err = read_image_to_frame(filename)
	} else {
		buf, err = utils.ReadFile(filename)
	}
	if err != nil {
		return err
	}
	//time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("Sending frame: %s (%v) (%d bytes)\n", filename, s.toggle, len(buf))
	offset := 0
	repeat := true
	for repeat {
		l := 8192
		if l+offset > len(buf) {
			l = len(buf) - offset
			repeat = false
		}
		data := buf[offset : offset+l]
		ist := &pb.FrameStream{
			NextImage: s.toggle,
			Data:      data,
		}
		err = s.srv.Send(ist)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("err:%s\n", err)
			return err
		}
		offset = offset + l
	}

	s.toggle = !s.toggle
	time.Sleep(*delay)
	return nil
}

func read_image_to_frame(filename string) ([]byte, error) {
	var err error
	var img image.Image
	if strings.Contains(filename, "yuv422") {
		img, err = read_yuv422_image(filename)
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		r := bufio.NewReader(f)
		img, _, err = image.Decode(r)
		f.Close()
	}
	if err != nil {
		return nil, err
	}

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, int(WIDTH), int(HEIGHT)))
	// Resize:
	//draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	//	draw.BiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	draw.CatmullRom.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	//	draw.ApproxBiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	fmt.Printf("%s: Scaled to %dx%d\n", filename, dst.Bounds().Max.X, dst.Bounds().Max.Y)
	ri, err := converters.ConvertToRaw(dst)
	if err != nil {
		return nil, err
	}
	return ri.DefaultBytes(), nil

}

func read_yuv422_image(filename string) (image.Image, error) {
	i := strings.LastIndex(filename, "_")
	if i == -1 {
		return nil, fmt.Errorf("missing _ in filename. (expected filename.yuv422_640x480 or whatever WidthxHeight")
	}
	whs := filename[i+1:]
	sx := strings.Split(whs, "x")
	if len(sx) != 2 {
		panic(fmt.Sprintf("dimensions not in format WIDTHxHEIGHT (%s)", whs))
	}
	w, err := strconv.Atoi(sx[0])
	if err != nil {
		panic(fmt.Sprintf("invalid width specified: %s", err))
	}
	h, err := strconv.Atoi(sx[1])
	if err != nil {
		panic(fmt.Sprintf("invalid height specified: %s", err))
	}

	raw, err := utils.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	res := converters.ConvertYUV422ToImage(raw, h, w)
	return res, nil
}
