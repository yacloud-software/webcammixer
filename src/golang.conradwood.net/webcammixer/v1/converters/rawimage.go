package converters

import (
	"fmt"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.conradwood.net/webcammixer/v1/defaults"
	"image"
	"image/color"
)

type RawImage struct {
	orig   image.Image
	height int
	width  int
	data   []byte
}

// get in format whatever default is configured
func (r *RawImage) DefaultBytes() []byte {
	format, _ := defaults.GetPreferredFormat()
	if format == v4l2.PixelFmtYUYV {
		return r.YUV422()
	} else if format == v4l2.PixelFmtRGB24 {
		return r.RGB24()
	} else {
		panic(fmt.Sprintf("unsupported format: %v", format))
	}

}

/*
In this format each four bytes is two pixels. Each four bytes is two Y's, a Cb and a Cr. Each Y goes to one of the pixels, and the Cb and Cr belong to both pixels. As you can see, the Cr and Cb components have half the horizontal resolution of the Y component. V4L2_PIX_FMT_YUYV is known in the Windows environment as YUY2.

Example 2.19. V4L2_PIX_FMT_YUYV 4 × 4 pixel image

Byte Order. Each cell is one byte.

start + 0:	Y'00	Cb00	Y'01	Cr00	Y'02	Cb01	Y'03	Cr01
start + 8:	Y'10	Cb10	Y'11	Cr10	Y'12	Cb11	Y'13	Cr11
start + 16:	Y'20	Cb20	Y'21	Cr20	Y'22	Cb21	Y'23	Cr21
start + 24:	Y'30	Cb30	Y'31	Cr30	Y'32	Cb31	Y'33	Cr31
Color Sample Location.

 	0	 	1	 	2	 	3
0	Y	C	Y	 	Y	C	Y
1	Y	C	Y	 	Y	C	Y
2	Y	C	Y	 	Y	C	Y
3	Y	C	Y	 	Y	C	Y

*/

func (r *RawImage) YUV422() []byte { // actualy YUV422
	original := r.orig
	bounds := original.Bounds()
	w := (bounds.Max.X - bounds.Min.X)
	h := (bounds.Max.Y - bounds.Min.Y)
	raw := make([]uint8, w*h*2)
	//	converted := image.NewYCbCr(image.Rect(0, 0, int(w), int(h)), image.YCbCrSubsampleRatio422)
	maxrow := bounds.Max.Y
	maxcol := bounds.Max.X
	for row := 0; row < maxrow-1; row++ {
		for col := 0; col < maxcol; col++ {
			r, g, b, _ := original.At(col, row).RGBA()
			y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))
			/*
				converted.Y[converted.YOffset(col, row)] = y
				converted.Cb[converted.COffset(col, row)] = cb
				converted.Cr[converted.COffset(col, row)] = cr
			*/ //raw[converted.YOffset(col, row)] = y

			pos := (row*maxcol + col)
			raw[pos*2] = y
			if (pos & 1) > 0 {
				raw[pos*2+1] = cr
			} else {
				raw[pos*2+1] = cb
			}
			//			raw[(row*maxcol+col)*2+2] = 0xff
			//			m:=maxrow*maxcol
			//			raw[m+(row*col/2+col)] = cb
			//			raw[m+(maxrow*maxcol/2)+(row*col/2)] = cr

		}
	}
	return raw
}

func (r *RawImage) RGB24() []byte {
	img := r.orig
	sz := img.Bounds()
	raw := make([]uint8, (sz.Max.X-sz.Min.X)*(sz.Max.Y-sz.Min.Y)*3)
	idx := 0
	for y := sz.Min.Y; y < sz.Max.Y; y++ {
		for x := sz.Min.X; x < sz.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			raw[idx], raw[idx+1], raw[idx+2] = uint8(r), uint8(g), uint8(b)
			idx += 3
		}
	}

	return raw
}
