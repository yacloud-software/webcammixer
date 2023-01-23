package webcam

import (
	"context"
	"fmt"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/defaults"
	"sort"
)

type resolution struct {
	height uint32
	width  uint32
}

func Open(devicename string, height, width uint32) (*WebCamInfo, error) {
	if !utils.FileExists(devicename) {
		return nil, fmt.Errorf("device \"%s\" does not exist", devicename)
	}
	//	needFormat := v4l2.PixelFmtRGB24
	//	bufsize := height * width * 3

	/*
		needFormat := v4l2.PixelFmtYUYV
		bl := width * 2
		bufsize := height * width * 2
	*/

	needFormat, factor := defaults.GetPreferredFormat()
	bl := width * factor
	bufsize := height * width * factor

	fmt.Printf("Opening device %s, size %dx%d\n", devicename, width, height)

	pf := v4l2.PixFormat{PixelFormat: needFormat,
		Width:        uint32(width),
		Height:       uint32(height),
		Field:        v4l2.FieldNone, // interlaced, or not
		Colorspace:   v4l2.ColorspaceSRGB,
		BytesPerLine: bl,
		SizeImage:    bufsize,
	}

	dev, err := device.Open(devicename,
		device.WithIOType(v4l2.IOTypeMMAP),
		//	device.WithPixFormat(pf),
		device.WithFPS(30),
	//	device.WithBufferSize(bufsize),
	)
	if err != nil {
		return nil, err
	}

	//getinfo
	fdes, err := v4l2.GetAllFormatDescriptions(dev.GetFD(), v4l2.BufTypeVideoCapture)
	if err != nil {
		dev.Close()
		return nil, err
	}
	if len(fdes) == 0 {
		dev.Close()
		return nil, fmt.Errorf("Got 0 format descriptions from device %s", devicename)
	}
	for _, fd := range fdes {
		fmt.Printf("(%s) FDES: %s\n", devicename, fd.Description)
	}
	fse, err := v4l2.GetAllFormatFrameSizes(dev.GetFD())
	if err != nil {
		dev.Close()
		return nil, err
	}
	var bigger_pixel_formats []*resolution
	for _, fs := range fse {
		h := fs.Size.MaxHeight
		w := fs.Size.MaxWidth
		if (h > height) && (w > width) && (fs.PixelFormat == needFormat) {
			bigger_pixel_formats = append(bigger_pixel_formats, &resolution{height: h, width: w})
		}
		fmt.Printf("(%s) FSE: %dx%d (%v) (%s)\n", devicename, w, h, fs.PixelFormat, v4l2.PixelFormats[fs.PixelFormat])
	}
	//fmt.Printf("BEFORE: %s\n", pf.String())
	sort.Slice(bigger_pixel_formats, func(i, j int) bool {
		if bigger_pixel_formats[i].height < bigger_pixel_formats[j].height {
			return true
		}
		if bigger_pixel_formats[i].width < bigger_pixel_formats[j].width {
			return true
		}
		return false
	})
	for i, bpf := range bigger_pixel_formats {
		fmt.Printf("Pixelformats %d: %dx%d\n", i, bpf.width, bpf.height)
	}
	next_try := 0
retry:
	npf, err := v4l2.SetPixFormat2(dev.GetFD(), pf, v4l2.BufTypeVideoCapture)
	//fmt.Printf("AFTER: %s\n", npf.String())
	if err != nil {
		dev.Close()
		return nil, err
	}
	if pf.Field != npf.Field {
		dev.Close()
		return nil, fmt.Errorf("%s - fields do not match. Wanted %d, got %d\n", devicename, pf.Field, npf.Field)
	}
	if needFormat != npf.PixelFormat {
		dev.Close()
		nf := v4l2.PixelFormats[needFormat]
		af := v4l2.PixelFormats[npf.PixelFormat]
		return nil, fmt.Errorf("%s: Wanted pixelformat %v (%s), but got %v (%s)", devicename, needFormat, nf, npf.PixelFormat, af)
	}

	if (npf.Height < height || npf.Width < width) && next_try < len(bigger_pixel_formats) {
		fmt.Printf("Next_try: %d, len=%d\n", next_try, len(bigger_pixel_formats))
		b := bigger_pixel_formats[next_try]
		pf = v4l2.PixFormat{PixelFormat: needFormat,
			Width:        b.width,
			Height:       b.height,
			Field:        v4l2.FieldNone, // interlaced, or not
			Colorspace:   v4l2.ColorspaceSRGB,
			BytesPerLine: bl,
			SizeImage:    bufsize,
		}
		goto retry
	}

	fmt.Printf("%s Actual pixelformat: %dx%d %v (wanted %v)\n", devicename, npf.Width, npf.Height, npf.PixelFormat, needFormat)
	wci := GetWebCamInfo(devicename)
	wci.device = dev
	wci.actualPixelFormat = npf
	fmt.Printf("Opened device \"%s\" @ %dx%d (bufsize=%d)\n", devicename, width, height, bufsize)
	return wci, nil
}
func (w *WebCamInfo) GetActualPixelFormat() v4l2.PixFormat {
	return w.actualPixelFormat
}
func (w *WebCamInfo) Close() error {
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	w.device.Stop()
	w.device.Close()
	return nil
}
func (w *WebCamInfo) Start() (<-chan []byte, error) { // definition of a read-only chan. didn't know that until now ;)
	ctx, cancel := context.WithCancel(context.TODO())
	if cancel != nil {
		w.cancel = cancel
	}
	err := w.device.Start(ctx)
	if err != nil {
		return nil, err
	}
	frameChan := w.device.GetOutput()
	return frameChan, nil
}
