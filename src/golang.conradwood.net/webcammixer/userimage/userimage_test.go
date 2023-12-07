package userimage

import (
	"bytes"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	//	"golang.conradwood.net/webcammixer/converters"
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
	uip := NewUserImageProvider(WIDTH, HEIGHT)
	go uip.Run()
	started := time.Now()
	i := 0
	for time.Since(started) < time.Duration(5)*time.Second {
		time.Sleep(time.Duration(100) * time.Millisecond)
		i++
		img := uip.GetImage()
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
