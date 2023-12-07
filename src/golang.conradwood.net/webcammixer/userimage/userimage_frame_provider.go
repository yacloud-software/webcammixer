package userimage

import (
	"flag"
	"fmt"
	_ "github.com/fogleman/gg"
	"golang.conradwood.net/webcammixer/converters"
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
	idle_sleep = flag.Duration("user_sleep", time.Duration(100)*time.Millisecond, "time to sleep between sending idle image updates")
)

type UserImageProvider struct {
	stop_requested     bool
	pos_going_right    bool
	pos_going_down     bool
	TextXPOSAdjust     int
	TextYPOSAdjust     int
	colour_going_down  bool
	last_colour_update time.Time
	cur_rgba           []uint8
	idle_text          func() string
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
func NewUserImageProvider(h, w uint32) *UserImageProvider {
	ifp := &UserImageProvider{
		cur_rgba: []uint8{200, 100, 0, 255},
		width:    w,
		height:   h,
	}
	return ifp
}
func (ifp *UserImageProvider) SetIdleText(s func() string) {
	ifp.idle_text = s
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
	for {
		if ifp.stop_requested {
			break
		}
		var err error
		err = ifp.createIdleImageText()
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
		time.Sleep(*idle_sleep)
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
	if ifp.idle_text != nil {
		txt = ifp.idle_text()
	}
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

func (ifp *UserImageProvider) GetID() string {
	return "userimageframe"
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
