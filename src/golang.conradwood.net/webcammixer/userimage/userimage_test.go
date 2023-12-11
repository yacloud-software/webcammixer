package userimage

import (
	"bytes"
	"fmt"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/webcam"
	"image"
	"image/png"
	"testing"
	"time"
)

const (
	WIDTH  = 640
	HEIGHT = 480
)

func TestLabeller(t *testing.T) {
	//img, err := get_emoji("test \U0001f600")
	//	img, err := get_emoji("\U0001f600")
	img, err := get_emoji("\U0001f3f4\U000e0067\U000e0062\U000e0065\U000e006e\U000e0067\U000e007f")
	if err != nil {
		t.Fatalf("failed to get well-known emoji: %s", err)
		return
	}
	save_image(0, img)
	uip := NewUserImageProvider(webcam.NewSourceMixer(), WIDTH, HEIGHT)
	uip.SetConfig(&pb.UserImageRequest{
		Converters: []*pb.UserImageConverter{
			&pb.UserImageConverter{Type: pb.ConverterType_LABEL, Text: "foobar"},
		},
	})
	go uip.Run()
	started := time.Now()
	i := 0
	for time.Since(started) < time.Duration(15)*time.Second {
		time.Sleep(time.Duration(100) * time.Millisecond)
		i++
		img := uip.conv_image
		err := save_image(i, img)
		if err != nil {
			t.Fatalf("failed to save frame: %s", err)
			return
		}
	}
}

func save_image(number int, img image.Image) error {
	if img == nil {
		return nil
	}
	w := &bytes.Buffer{}
	err := png.Encode(w, img)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("/tmp/testpix/test_%d.png", number)
	b := w.Bytes()
	//	err = utils.WriteFile(filename, b)
	if err != nil {
		return err
	}

	filename = fmt.Sprintf("/tmp/testpix/test.png")
	err = utils.WriteFile(filename, b)
	if err != nil {
		return err
	}
	return nil
}
