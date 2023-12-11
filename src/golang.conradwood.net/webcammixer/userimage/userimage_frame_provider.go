package userimage

import (
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/interfaces"
	"golang.conradwood.net/webcammixer/labeller"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"time"
)

const (
	TEXT_MINX = 100
	TEXT_MINY = 100
)

var (
	idle_sleep = flag.Duration("user_sleep", time.Duration(50)*time.Millisecond, "time to sleep between sending idle image updates")
)

type UserImageProvider struct {
	timer_source       chan bool
	sourceMixer        interfaces.SourceMixer
	stop_requested     bool
	pos_going_right    bool
	pos_going_down     bool
	TextXPOSAdjust     int
	TextYPOSAdjust     int
	colour_going_down  bool
	last_colour_update time.Time
	cur_rgba           []uint8
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
func NewUserImageProvider(sm interfaces.SourceMixer, h, w uint32) *UserImageProvider {
	ifp := &UserImageProvider{
		sourceMixer: sm,
		cur_rgba:    []uint8{200, 100, 0, 255},
		width:       w,
		height:      h,
	}
	return ifp
}
func (ifp *UserImageProvider) getSourceMixer() interfaces.SourceMixer {
	return ifp.sourceMixer
}

// if target is non-null will send message to channel each time there is a new frame
func (ifp *UserImageProvider) SetTimerTarget(target chan bool) error {
	ifp.notify = target
	return nil
}

func (ifp *UserImageProvider) GetDimensions() (uint32, uint32) {
	return ifp.height, ifp.width
}

func (ifp *UserImageProvider) Run() error {
	if ifp == nil {
		return fmt.Errorf("not running empty user image provider")
	}
	h, w := ifp.GetDimensions()
	fmt.Printf("Starting userframe provider with dimensions %dx%d...\n", w, h)
	p := utils.ProgressReporter{}
	for {
		if ifp.stop_requested {
			break
		}
		var err error
		err = ifp.createImage()
		if err != nil {
			fmt.Printf("failed to render idle text: %s\n", err)
			time.Sleep(time.Duration(1500) * time.Millisecond)
			continue
		}
		//		sys.Write(int(l.Device.GetFD()), b)
		c := ifp.notify
		if c != nil {
			//			fmt.Printf("notifying dev about new frame\n")
			// notify loopdev that we have a new frame
			c <- true
		}
		timesource := ifp.timer_source
		if timesource == nil {
			time.Sleep(*idle_sleep)
		} else {
			select {
			case <-time.After(*idle_sleep):
				//
			case <-timesource:
			}
		}
		p.Add(1)
		p.Print()
	}
	fmt.Printf("UserImageProvider stopped\n")
	return nil
}
func (ifp *UserImageProvider) Stop() {
	ifp.stop_requested = true
}
func (ifp *UserImageProvider) createIdleImageText() error {
	ifp.CalcColour()
	col := color.RGBA{ifp.cur_rgba[0], ifp.cur_rgba[1], ifp.cur_rgba[2], ifp.cur_rgba[3]}
	xsize := 640
	ysize := 480
	h, w := ifp.GetDimensions()
	xsize = int(w)
	ysize = int(h)
	txt := "."
	txt = "no-test"

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
	ifp.idleImage = l.GetImage()
	rawimage, err := converters.ConvertToRaw(ifp.idleImage)
	if err != nil {
		return err
	}
	ifp.idleImageRaw = rawimage.DefaultBytes()
	ifp.idleImageLastFile = ""
	return nil
}

func (ifp *UserImageProvider) GetID() string {
	return "userimageframe"
}

func (ifp *UserImageProvider) GetImage() image.Image {
	return ifp.idleImage
}

// called by looopback
func (ifp *UserImageProvider) GetFrame() ([]byte, error) {
	//	ifp.CalcColour()
	//	fmt.Printf("[%v] Getting frame...\n", time.Now())
	return ifp.idleImageRaw, nil
}

func (ifp *UserImageProvider) CalcColour() {
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

func (ifp *UserImageProvider) createImage() error {
	cc := current_config
	if cc == nil {
		return ifp.createIdleImageText()
	}

	h, w := ifp.GetDimensions()
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	//	draw.Draw(img, img.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
	dc := gg.NewContextForImage(img)
	for _, conv := range cc.converters {
		err := conv.Modify(dc)
		if err != nil {
			return err
		}
	}
	ifp.idleImage = dc.Image()
	rawimage, err := converters.ConvertToRaw(ifp.idleImage)
	if err != nil {
		return err
	}
	ifp.idleImageRaw = rawimage.DefaultBytes()
	ifp.idleImageLastFile = ""

	return nil
}
