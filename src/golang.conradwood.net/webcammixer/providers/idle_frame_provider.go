package providers

import (
	"bytes"
	"flag"
	"fmt"
	_ "github.com/fogleman/gg"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/labeller"
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

const (
	TEXT_MINX = 100
	TEXT_MINY = 100
)

var (
	idle_sleep_send   = flag.Duration("idle_sleep_send", time.Duration(1000)*time.Millisecond, "time between sending idle image updates")
	idle_sleep_create = flag.Duration("idle_create_image", time.Duration(1000)*time.Millisecond, "time between recreating idle images")

	idle_text      = flag.String("idle_text", "idle", "if set, display an idle text instead of idle image")
	loopback_image = flag.String("idle_image", "/tmp/idle_image.jpeg", "set this to the image to be served whilst no camera is attached")
)

type IdleFrameProvider struct {
	wake_idle          chan bool
	pos_going_right    bool
	pos_going_down     bool
	TextXPOSAdjust     int
	TextYPOSAdjust     int
	colour_going_down  bool
	last_colour_update time.Time
	cur_rgba           []uint8
	idle_text          string
	width              uint32
	height             uint32
	idleImage          image.Image
	idleImageRaw       []byte
	idleImageLastMod   time.Time
	idleImageLastFile  string
	lastText           string
	notify             chan bool
}

// blocks and provides "idle" frame
func NewIdleFrameProvider(h, w uint32) *IdleFrameProvider {
	ifp := &IdleFrameProvider{
		cur_rgba:  []uint8{200, 100, 0, 255},
		width:     w,
		height:    h,
		idle_text: *idle_text,
		wake_idle: make(chan bool, 20),
	}

	return ifp
}
func (ifp *IdleFrameProvider) TriggerFrameNow() {
	ifp.wake_idle <- true
}
func (ifp *IdleFrameProvider) SetIdleText(s string) {
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
	var last_created time.Time
	for {
		var err error
		if time.Since(last_created) >= *idle_sleep_create || ifp.idle_text != ifp.lastText {
			if ifp.idle_text != "" {
				err = ifp.createIdleImageText()
			} else {
				err = ifp.loadIdleImage()
			}
			if err != nil {
				fmt.Printf("failed to load idleimage: %s\n", err)
				time.Sleep(time.Duration(1500) * time.Millisecond)
				continue
			}
			last_created = time.Now()
		}
		//              sys.Write(int(l.Device.GetFD()), b)
		c := ifp.notify
		if c != nil {
			//                      fmt.Printf("notifying dev about new frame\n")
			// notify loopdev that we have a new frame
			c <- true
		}
		select {
		case <-time.After(*idle_sleep_send):
			//
		case <-ifp.wake_idle:
			//
		}
		//		time.Sleep(*idle_sleep_send)
	}
}

func (ifp *IdleFrameProvider) createIdleImageText() error {
	ifp.CalcColour()
	col := color.RGBA{ifp.cur_rgba[0], ifp.cur_rgba[1], ifp.cur_rgba[2], ifp.cur_rgba[3]}
	xsize := 640
	ysize := 480
	h, w := ifp.GetDimensions()
	xsize = int(w)
	ysize = int(h)
	txt := ifp.idle_text
	l := labeller.NewLabellerForBlankCanvas(xsize, ysize, color.RGBA{0, 0, 0, 255})
	l.SetFontSize(80)

	x := 0
	y := 0
	x = x + ifp.TextXPOSAdjust
	y = y + ifp.TextYPOSAdjust

	ld := l.NewLabel(x, y, col, txt)
	err := l.PaintLabel(ld)
	if err != nil {
		return err
	}
	rawimage, err := converters.ConvertToRaw(l.GetImage())
	if err != nil {
		return err
	}
	ifp.idleImageRaw = rawimage.DefaultBytes()
	ifp.idleImageLastFile = ""
	return nil

}
func (ifp *IdleFrameProvider) OLD_createIdleImageText() error {
	if ifp.idle_text == "" {
		fmt.Printf("No idle-text, idletext called\n")
		return fmt.Errorf("no idle text")
	}
	label := ifp.idle_text
	//	if label == ifp.lastText {
	//		return nil
	//	}
	ifp.CalcColour()
	fmt.Printf("Rendering text \"%s\" (colour:%v)\n", label, ifp.cur_rgba)
	h, w := uint32(200), uint32(200) //ifp.GetDimensions()
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	//	col := color.RGBA{200, 100, 0, 255}
	col := color.RGBA{ifp.cur_rgba[0], ifp.cur_rgba[1], ifp.cur_rgba[2], ifp.cur_rgba[3]}
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
	y = y + ifp.TextYPOSAdjust
	wpos = wpos + ifp.TextXPOSAdjust
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
	//	ifp.CalcColour()
	//	fmt.Printf("[%v] Getting frame...\n", time.Now())
	return ifp.idleImageRaw, nil
}

func (ifp *IdleFrameProvider) CalcColour() {
	diff := time.Since(ifp.last_colour_update)
	if diff < time.Duration(100)*time.Millisecond {
		return
	}
	adjust := uint8(0)
	if diff > 20 {
		adjust = 20
	} else {
		adjust = uint8(diff)
	}

	ifp.last_colour_update = time.Now()

	if ifp.cur_rgba[0] <= 100 {
		ifp.colour_going_down = false
	}
	if ifp.cur_rgba[0] >= 250 {
		ifp.colour_going_down = true
	}
	if ifp.colour_going_down {
		ifp.cur_rgba[0] = ifp.cur_rgba[0] - adjust
	} else {
		ifp.cur_rgba[0] = ifp.cur_rgba[0] + adjust
	}

	if ifp.TextXPOSAdjust > 500 {
		ifp.pos_going_right = false
	}
	if ifp.TextXPOSAdjust < TEXT_MINX {
		ifp.TextXPOSAdjust = TEXT_MINX
		ifp.pos_going_right = true
	}

	if ifp.TextYPOSAdjust > 500 {
		ifp.pos_going_down = false
	}
	if ifp.TextYPOSAdjust < TEXT_MINY {
		ifp.TextYPOSAdjust = TEXT_MINY
		ifp.pos_going_down = true
	}

	pos_adjust := 1
	if ifp.pos_going_right {
		ifp.TextXPOSAdjust = ifp.TextXPOSAdjust + pos_adjust
	} else {
		ifp.TextXPOSAdjust = ifp.TextXPOSAdjust - pos_adjust
	}
	pos_adjust = 1
	if ifp.pos_going_down {
		ifp.TextYPOSAdjust = ifp.TextYPOSAdjust + pos_adjust
	} else {
		ifp.TextYPOSAdjust = ifp.TextYPOSAdjust - pos_adjust
	}
}
