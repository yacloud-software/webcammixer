package userimage

import (
	"github.com/fogleman/gg"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/interfaces"
	"golang.conradwood.net/webcammixer/labeller"
	//	"golang.conradwood.net/webcammixer/webcam"
	//	"image/color"
	"fmt"
	"image"
	//	"image/draw"
)

// the converter modifies the image. a chain of converters modifies the image at each step
type converter struct {
	cfg     *config
	convdef *pb.UserImageConverter
	typ     pb.ConverterType
	//	lab     *labeller.Labeller // implements ext_converter
	text               func() string
	tmv                *text_mover
	vcs                interfaces.VideoCamSource
	extbin             ExtBinary
	image              image.Image // for image modifier
	image_x            uint32
	image_y            uint32
	input_has_changed  bool // reset by renderer and set if text/image are changed
	last_text_rendered *last_text
}

// what was last text like? (to compute input_has_changed)
type last_text struct {
	text string
	X    int
	Y    int
}

// returns true if the image was _actually_ modified
func (c *converter) Modify(gctx *gg.Context) (bool, error) {
	if c.typ == pb.ConverterType_WEBCAM {
		return c.webcam(gctx)
	} else if c.typ == pb.ConverterType_LABEL {
		return c.modify_text(gctx)
	} else if c.typ == pb.ConverterType_EXT_BINARY {
		return c.modify_through_extbin(gctx)
	} else if c.typ == pb.ConverterType_OVERLAY_IMAGE {
		return c.modify_image(gctx)
	} else {
		return false, fmt.Errorf("cannot modify \"%v\" yet", c.typ)
	}
}

func (c *converter) HasChanged() bool {
	if c.typ == pb.ConverterType_WEBCAM {
		return true
	} else if c.typ == pb.ConverterType_LABEL {
		return c.has_changed_text()
	} else if c.typ == pb.ConverterType_EXT_BINARY {
		return true
	} else if c.typ == pb.ConverterType_OVERLAY_IMAGE {
		return c.has_changed_image()
	} else {
		return false
	}
}

func (c *converter) modify_through_extbin(gctx *gg.Context) (bool, error) {
	//TODO: figure out an input format for ext_binary
	h, w := c.cfg.ifp.GetDimensions()
	c.extbin.Call_Modify(nil, h, w)
	return true, nil
}

func (c *converter) webcam(gctx *gg.Context) (bool, error) {
	frame, err := c.vcs.GetFrame()
	if err != nil {
		return false, err
	}
	if len(frame) == 0 {
		return false, nil
	}
	h, w := c.cfg.ifp.GetDimensions()
	img := converters.ConvertYUV422ToImage(frame, int(h), int(w))
	gctx.DrawImage(img, 0, 0)

	return true, nil
}
func (c *converter) has_changed_text() bool {
	c.tmv.Step()
	lt := c.last_text_rendered
	if lt == nil {
		return true
	}
	if lt.X == c.tmv.X() && lt.Y == c.tmv.Y() && c.text() != "" && c.text() == lt.text {
		return false
	}
	return c.input_has_changed
}
func (c *converter) has_changed_image() bool {
	return c.input_has_changed
}
func (c *converter) modify_image(gctx *gg.Context) (bool, error) {
	img := c.image
	if img == nil {
		return false, nil
	}
	gctx.DrawImage(img, int(c.image_x), int(c.image_y))
	//	fmt.Printf("Drawing image at %dx%d\n", c.image_x, c.image_y)
	c.input_has_changed = false
	return true, nil
}
func (c *converter) modify_text(gctx *gg.Context) (bool, error) {
	if c.text == nil {
		return false, nil
	}
	lt := c.last_text_rendered
	if lt == nil {
		lt = &last_text{}
		c.last_text_rendered = lt
	}

	//	ifp := c.cfg.ifp
	col := c.tmv.RGBA()
	txt := c.text()
	if txt == "" {
		return false, nil
	}
	x := c.tmv.X()
	y := c.tmv.Y()

	l := labeller.NewLabellerForCtx(gctx)
	l.SetFontSize(60)

	ld := l.NewLabel(x, y, col, txt)
	err := l.PaintLabel(ld)
	if err != nil {
		return false, err
	}
	h, w := l.GetMaxDimensions()
	c.tmv.text_height = h
	c.tmv.text_width = w
	//	c.input_has_changed = false
	//	fmt.Printf("Label \"%s\" added at %dx%d\n", txt, x, y)
	lt.X = x
	lt.Y = y
	lt.text = txt
	return true, nil
}
