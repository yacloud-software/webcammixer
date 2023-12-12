package userimage

import (
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/interfaces"
	"golang.conradwood.net/webcammixer/rates"
	"image"
	//	"image/draw"
	//	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"time"
)

var (
	mod_sleep  = flag.Duration("user_modify_sleep", time.Duration(00)*time.Millisecond, "minimum time the modifier loop should take")
	idle_sleep = flag.Duration("user_sleep", time.Duration(200)*time.Millisecond, "time to sleep between sending idle image updates")
)

type UserImageProvider struct {
	imageSource              ImageSource
	sourceMixer              interfaces.SourceMixer
	stop_requested           bool
	width                    uint32
	height                   uint32
	frame                    []byte      // last frame from source
	conv_frame               []byte      // last frame coming out of converter-chain (same as image)
	conv_image               image.Image // last iamge coming out of converter-chain (same as frame)
	notify                   chan bool
	merge_chan               chan bool
	modify_chan              chan bool // fed by video source, consumed by modify loop
	converters_had_no_impact bool      // set to true if no converter had any impact (produced an image of sort)
	last_printed_modify      time.Time
}

// blocks and provides "idle" frame
func NewUserImageProvider(sm interfaces.SourceMixer, h, w uint32) *UserImageProvider {
	ifp := &UserImageProvider{
		sourceMixer: sm,
		width:       w,
		height:      h,
		merge_chan:  make(chan bool, 30),
		modify_chan: make(chan bool, 30),
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

// return height,width
func (ifp *UserImageProvider) GetDimensions() (uint32, uint32) {
	return ifp.height, ifp.width
}

func (ifp *UserImageProvider) Run() error {
	if ifp == nil {
		return fmt.Errorf("not running empty user image provider")
	}
	h, w := ifp.GetDimensions()
	fmt.Printf("Starting userframe provider with dimensions %dx%d...\n", w, h)
	go ifp.modifier_loop()
	go ifp.merge_loop()
	for !ifp.stop_requested {
		src := ifp.imageSource
		if src == nil {
			time.Sleep(*idle_sleep)
			continue
		}

		select {
		case <-time.After(*idle_sleep):
			rates.Inc("videosource_timeout")
		case <-src.GetTimingChannel():
			rates.Inc("videosource_newframe")
		}
		if len(ifp.modify_chan) < 10 {
			ifp.modify_chan <- true
		}
		if len(ifp.merge_chan) < 10 {
			ifp.merge_chan <- true
		}
	}
	fmt.Printf("UserImageProvider stopped\n")
	if ifp.imageSource != nil {
		ifp.imageSource.Close()
		ifp.imageSource = nil
	}
	return nil
}

// this runs asynchronous to merge the update(s) with the source frames and send them to loop back device
func (ifp *UserImageProvider) merge_loop() {
	for !ifp.stop_requested {
		<-ifp.merge_chan
		src := ifp.imageSource
		if src == nil {
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}

		srcframe, err := src.GetFrame()
		if err != nil {
			fmt.Printf("Error getting frame: %s\n", err)
			continue
		}
		if len(srcframe) == 0 {
			continue
		}
		if ifp.converters_had_no_impact {
			ifp.frame = srcframe
		} else {
			ifp.merge_frame_with_conv(srcframe)
		}
		c := ifp.notify
		if c != nil {
			//			fmt.Printf("notifying dev about new frame\n")
			// notify loopdev that we have a new frame
			select {
			case c <- true:
				// sent
			default:
				// not sent
			}
		}
	}
}

// this runs asynchronous to generate the update(s), e.g. text/image
// it _does_ however use a channel to not run any faster than the video source
func (ifp *UserImageProvider) modifier_loop() {
	var last_run time.Time
	for !ifp.stop_requested {
		<-ifp.modify_chan

		min_wait := *mod_sleep
		if min_wait > 0 {
			diff := time.Since(last_run)
			if diff < min_wait {
				time.Sleep(min_wait - diff)
			}
		}

		var err error
		err = ifp.createImage()
		last_run = time.Now()
		if err != nil {
			fmt.Printf("failed to render idle text: %s\n", err)
			time.Sleep(time.Duration(1500) * time.Millisecond)
			continue
		}
	}

}

func (ifp *UserImageProvider) NewSource(is ImageSource) {
	if ifp.imageSource != nil {
		ifp.imageSource.Close()
	}
	ifp.frame = nil
	ifp.conv_frame = nil
	ifp.conv_image = nil
	ifp.imageSource = is
	is.Activate()
}

func (ifp *UserImageProvider) Stop() {
	ifp.stop_requested = true
	ifp.modify_chan <- true
	ifp.merge_chan <- true
}

func (ifp *UserImageProvider) GetID() string {
	return "userimageframe"
}

// called by looopback
func (ifp *UserImageProvider) GetFrame() ([]byte, error) {
	//	ifp.CalcColour()
	//	fmt.Printf("[%v] Getting frame...\n", time.Now())
	return ifp.frame, nil
}

func (ifp *UserImageProvider) createImage() error {
	cc := current_config
	if cc == nil {
		ifp.converters_had_no_impact = true
		return nil
	}

	h, w := ifp.GetDimensions()
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	dc := gg.NewContextForImage(img)
	has_changed := false
	if len(ifp.conv_frame) == 0 && ifp.conv_image == nil {
		has_changed = true // we don't have an image yet, do run it
	}
	for _, conv := range cc.converters {
		if conv.HasChanged() {
			has_changed = true
			break
		}
	}
	if !has_changed {
		ifp.printModify("modifiers have not changed (no_impact=%v)\n", ifp.converters_had_no_impact)
		return nil
	}
	had_impact := false
	for _, conv := range cc.converters {
		dc.ResetClip()
		b, err := conv.Modify(dc)
		if err != nil {
			return err
		}
		if b {
			had_impact = true
		}
	}
	if !had_impact {
		ifp.converters_had_no_impact = true
		ifp.printModify("modifiers have no impact\n")
		return nil
	}
	ifp.converters_had_no_impact = false
	new_img := dc.Image()
	rawimage, err := converters.ConvertToRaw(new_img)
	if err != nil {
		return err
	}
	ifp.conv_frame = rawimage.DefaultBytes()
	ifp.conv_image = new_img
	rates.Inc("modify")
	return nil
}
func (ifp *UserImageProvider) printModify(format string, args ...interface{}) {
	if time.Since(ifp.last_printed_modify) < time.Duration(2)*time.Second {
		return
	}
	ifp.last_printed_modify = time.Now()
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[modify]%s", s)
}
func (ifp *UserImageProvider) merge_frame_with_conv(srcframe []byte) {
	if len(srcframe) == 0 {
		return
	}
	if len(ifp.conv_frame) == 0 && ifp.conv_image == nil {
		ifp.frame = srcframe
		// no conversion (yet)
		return
	}
	conv_frame := ifp.conv_frame
	conv_image := ifp.conv_image
	if len(conv_frame) != 0 {
		// merge frame2frame
	}
	if conv_image == nil {
		ifp.frame = srcframe
		return
	}
	if conv_frame == nil {
		ifp.frame = srcframe
		return
	}
	h, w := ifp.GetDimensions()
	img := converters.ConvertYUV422ToImage(srcframe, int(h), int(w))
	if img == nil {
		ifp.frame = srcframe
		return
	}
	gctx := gg.NewContextForImage(img)
	gctx.DrawImage(conv_image, 0, 0)
	img = gctx.Image()
	rawimage, err := converters.ConvertToRaw(img)
	if err != nil {
		ifp.frame = srcframe
		return
	}
	ifp.frame = rawimage.DefaultBytes()
	rates.Inc("merge")
}
