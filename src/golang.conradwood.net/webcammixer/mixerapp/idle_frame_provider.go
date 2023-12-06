package mixerapp

import (
	"bytes"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/converters"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"
)

var (
	idle_sleep     = flag.Duration("idle_sleep", time.Duration(1)*time.Second, "time to sleep between sending idle image updates")
	idle_text      = flag.String("idle_text", "idle", "if set, display an idle text instead of idle image")
	loopback_image = flag.String("idle_image", "/tmp/idle_image.jpeg", "set this to the image to be served whilst no camera is attached")
)

type IdleFrameProvider struct {
	idle_text         func() string
	width             uint32
	height            uint32
	idleImage         image.Image
	idleImageRaw      []byte
	idleImageLastMod  time.Time
	idleImageLastFile string
	lastText          string
	notify            chan bool
}

// blocks and provides "idle" frame
func NewIdleFrameProvider(h, w uint32) *IdleFrameProvider {
	ifp := &IdleFrameProvider{
		width:  w,
		height: h,
	}
	if *idle_text != "" {
		ifp.idle_text = func() string {
			return *idle_text
		}
	}
	return ifp
}
func (ifp *IdleFrameProvider) SetIdleText(s func() string) {
	ifp.idle_text = s
}

// if target is non-null will send message to channel each time there is a new frame
func (ifp *IdleFrameProvider) SetTimerTarget(target chan bool) error {
	ifp.notify = target
	return nil
}

func (ifp *IdleFrameProvider) GetDimensions() (uint32, uint32) {
	return ifp.height, ifp.width
}

func (ifp *IdleFrameProvider) Run() error {
	fmt.Printf("Starting idleframe provider...\n")
	for {
		var err error
		if ifp.idle_text != nil {
			err = ifp.createIdleImageText()
		} else {
			err = ifp.loadIdleImage()
		}
		if err != nil {
			fmt.Printf("failed to load idleimage: %s\n", err)
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		//		sys.Write(int(l.Device.GetFD()), b)
		c := ifp.notify
		if c != nil {
			//			fmt.Printf("notifying dev about new frame\n")
			// notify loopdev that we have a new frame
			c <- true
		}
		time.Sleep(*idle_sleep)
	}
}
func (ifp *IdleFrameProvider) createIdleImageText() error {
	if ifp.idle_text == nil {
		fmt.Printf("No idle-text, idletext called\n")
		return fmt.Errorf("no idle text")
	}
	label := ifp.idle_text()
	if label == ifp.lastText {
		return nil
	}
	fmt.Printf("Rendering text \"%s\"\n", label)
	h, w := uint32(200), uint32(200) //ifp.GetDimensions()
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	col := color.RGBA{200, 100, 0, 255}
	//x := int(w / 2)
	y := int(h / 2)
	point := fixed.Point26_6{X: fixed.I(0), Y: fixed.I(0)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	b, _ := d.BoundString(label)
	xsize := int(b.Max.X >> 6)
	wpos := int(w/2) - xsize/2
	point = fixed.Point26_6{X: fixed.I(wpos), Y: fixed.I(y)}
	d.Dot = point
	//fmt.Printf("b=%d,xsize:%d, Wpos: %d\n", b, xsize, wpos)
	d.DrawString(label)
	h, w = ifp.GetDimensions()
	rd := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	draw.NearestNeighbor.Scale(rd, rd.Rect, img, img.Bounds(), draw.Over, nil)

	rawimage, err := converters.ConvertToRaw(rd)
	if err != nil {
		return err
	}
	ifp.idleImageRaw = rawimage.DefaultBytes()
	ifp.lastText = label
	ifp.idleImageLastFile = ""
	return nil
}
func (ifp *IdleFrameProvider) loadIdleImage() error {
	reread := false
	if *loopback_image != ifp.idleImageLastFile {
		reread = true
	}
	if !reread {
		file, err := os.Stat(*loopback_image)
		if err != nil {
			return err
		}
		if !file.ModTime().Equal(ifp.idleImageLastMod) {
			reread = true
		}
	}
	if reread {
		fmt.Printf("reading idleimage \"%s\"\n", *loopback_image)
	}
	b, err := utils.ReadFile(*loopback_image)
	if err != nil {
		return fmt.Errorf("failed to read idleframe: %s", err)
	}
	r := bytes.NewReader(b)
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	ifp.idleImage = img
	h, w := ifp.GetDimensions()
	// Set the expected size that you want:
	rd := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	// Resize:
	draw.NearestNeighbor.Scale(rd, rd.Rect, img, img.Bounds(), draw.Over, nil)
	dst := rd
	rawimage, err := converters.ConvertToRaw(dst)
	if err != nil {
		return err
	}
	ifp.idleImageRaw = rawimage.DefaultBytes()
	file, err := os.Stat(*loopback_image)
	if err != nil {
		return err
	}
	ifp.idleImageLastMod = file.ModTime()
	ifp.idleImageLastFile = *loopback_image
	ifp.lastText = ""
	return nil
}

func (ifp *IdleFrameProvider) GetID() string {
	return "idleframe"
}

// called by looopback
func (ifp *IdleFrameProvider) GetFrame() ([]byte, error) {
	return ifp.idleImageRaw, nil
}
