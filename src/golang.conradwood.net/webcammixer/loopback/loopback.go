package loopback

import (
	"flag"
	"fmt"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.conradwood.net/webcammixer/defaults"
	sys "golang.org/x/sys/unix"
	"sync"
	"time"
)

var (
	debug = flag.Bool("debug_loopback", false, "debug loopback stuff")
)

type LoopBackDevice struct {
	timingChan        chan bool
	DeviceName        string
	Device            *device.Device
	displayIdleFrame  bool
	writingLock       sync.Mutex
	lastImageWritten  time.Time
	AutoIdleDuration  time.Duration
	frameWidth        uint32
	frameHeight       uint32
	frameProviders    map[string]FrameProvider
	frameProviderLock sync.Mutex
	bufsize           uint32
	timingSources     []TimingSource
}

type TimingSource interface {
	GetID() string
	SetTimerTarget(chan bool) error
}
type FrameProvider interface {
	GetID() string
	GetFrame() ([]byte, error)
}

func Open(name string, height, width uint32) (*LoopBackDevice, error) {
	res := &LoopBackDevice{
		timingChan:       make(chan bool),
		frameHeight:      height,
		frameWidth:       width,
		DeviceName:       name,
		displayIdleFrame: true,
		AutoIdleDuration: time.Duration(5) * time.Second,
		lastImageWritten: time.Now(),
	}
	fm, factor := defaults.GetPreferredFormat()
	res.bufsize = res.frameHeight * res.frameWidth * factor
	fmt.Printf("Opening loopback device %s size %dx%d\n", res.DeviceName, res.frameWidth, res.frameHeight)
	dev, err := device.Open(res.DeviceName, device.WithIOType(v4l2.IOTypeMMAP))
	if err != nil {
		return nil, err
	}
	res.Device = dev
	fdes, err := v4l2.GetAllFormatDescriptions(res.Device.GetFD(), v4l2.BufTypeVideoOutput)
	if err != nil {
		return nil, err
	}
	if len(fdes) == 0 {
		return nil, fmt.Errorf("Got 0 format descriptions from device %s", res.DeviceName)
	}
	if *debug {
		for _, fd := range fdes {
			fmt.Printf("FDES: %s\n", fd.Description)
		}
	}
	fd := res.Device.GetFD()
	pixformat, err := v4l2.GetPixFormat(fd, v4l2.BufTypeVideoOutput)
	if err != nil {
		return nil, err
	}
	if *debug {
		fmt.Printf("PIX Format: %#v\n", pixformat)
	}

	pixformat.Width = res.frameWidth
	pixformat.Height = res.frameHeight
	pixformat.BytesPerLine = pixformat.Width * 3 // rgb
	//	pixformat.PixelFormat = v4l2.PixelFmtYUV420
	pixformat.PixelFormat = fm //v4l2.PixelFmtRGB24
	//	pixformat.SizeImage = 1382400
	pixformat.SizeImage = res.bufsize
	pixformat.Field = v4l2.FieldNone

	fmt.Printf("Loopback device %s: setting format to %dx%d (bytesperline=%d,bufsize=%d)\n", res.DeviceName, pixformat.Width, pixformat.Height, pixformat.BytesPerLine, pixformat.SizeImage)
	npf, err := v4l2.SetPixFormat2(fd, pixformat, v4l2.BufTypeVideoOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to set format on device %s: %s", res.DeviceName, err)
	}

	//	v4l2.SetPixFormat(res.Device.GetFD(), fdes[0], v4l2.BufTypeVideoOutput)
	fmt.Printf("Opened loopback device %s size %dx%d (bytesperline=%d,bufsize=%d,field=%v) \n", res.DeviceName, npf.Width, npf.Height, npf.BytesPerLine, npf.SizeImage, npf.Field)
	return res, nil
}

// display the idleframe...
func (l *LoopBackDevice) StartWriter() error {
	for {
		_ = <-l.timingChan
		err := l.NewFrame()
		if err != nil {
			fmt.Printf("failed to write frame: %s\n", err)
		}

	}
}

func (l *LoopBackDevice) TimeOfLastFrame() time.Time {
	return l.lastImageWritten
}

// stop the idleframe and instead display this image
func (l *LoopBackDevice) WriteImage(buf []byte) error {
	l.displayIdleFrame = false
	l.writingLock.Lock()
	defer l.writingLock.Unlock()
	sys.Write(int(l.Device.GetFD()), buf)
	l.lastImageWritten = time.Now()
	return nil
}

// return height,width
func (l *LoopBackDevice) GetDimensions() (uint32, uint32) {
	return l.frameHeight, l.frameWidth
}

func (l *LoopBackDevice) SetTimingSource(ts TimingSource) {
	l.writingLock.Lock()
	l.timingSources = append(l.timingSources, ts)
	for _, ots := range l.timingSources {
		ots.SetTimerTarget(nil)
	}
	l.writingLock.Unlock()
	ts.SetTimerTarget(l.timingChan)
}
