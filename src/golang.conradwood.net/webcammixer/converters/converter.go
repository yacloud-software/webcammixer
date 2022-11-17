package converters

import (
	"image"
)

type Converter interface {
	Convert(*RawImage) (*RawImage, error)
}

func ConvertYUV422ToImage(frame []byte, height, width int) image.Image {
	res := image.NewYCbCr(image.Rect(0, 0, width, height), image.YCbCrSubsampleRatio422)
	maxrow := height
	maxcol := width
	for row := 0; row < maxrow-1; row++ {
		for col := 0; col < maxcol; col++ {
			yi := res.YOffset(col, row)
			ci := res.COffset(col, row)
			pos := (row*maxcol + col)
			res.Y[yi] = frame[pos*2]
			if (pos & 1) > 0 {
				res.Cr[ci] = frame[pos*2+1]
			} else {
				res.Cb[ci] = frame[pos*2+1]
			}
		}
	}
	return res
}
func ConvertToRaw(img image.Image) (*RawImage, error) {
	return &RawImage{orig: img}, nil
}

/*
	sz := img.Bounds()
	raw := make([]uint8, (sz.Max.X-sz.Min.X)*(sz.Max.Y-sz.Min.Y)*4)
	idx := 0
	for y := sz.Min.Y; y < sz.Max.Y; y++ {
		for x := sz.Min.X; x < sz.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			raw[idx], raw[idx+1], raw[idx+2], raw[idx+3] = uint8(r), uint8(g), uint8(b), uint8(a)
			idx += 4
		}
	}
	ri := &RawImage{
		height: sz.Max.X - sz.Min.X,
		width:  sz.Max.Y - sz.Min.Y,
		data:   raw,
	}
	return ri, nil
*/
