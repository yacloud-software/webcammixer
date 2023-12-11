package userimage

import (
	"bytes"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"image"
	"image/png"
)

func SavePNG(filename string, img image.Image) error {
	w := &bytes.Buffer{}
	err := png.Encode(w, img)
	if err != nil {
		fmt.Printf("Failed to encode file %s: %s\n", filename, err)
		return err
	}
	b := w.Bytes()
	err = utils.WriteFile(filename, b)
	if err != nil {
		fmt.Printf("Failed to write file %s: %s\n", filename, err)
		return err
	}
	return nil
}
