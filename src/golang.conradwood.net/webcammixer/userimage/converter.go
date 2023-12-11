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
	text    func() string
	tmv     *text_mover
	vcs     interfaces.VideoCamSource
	extbin  ExtBinary
	image   image.Image // for image modifier
	image_x uint32
	image_y uint32
}

func (c *converter) Modify(gctx *gg.Context) error {
	if c.typ == pb.ConverterType_WEBCAM {
		return c.webcam(gctx)
	} else if c.typ == pb.ConverterType_LABEL {
		return c.modify_text(gctx)
	} else if c.typ == pb.ConverterType_EXT_BINARY {
		return c.modify_through_extbin(gctx)
	} else if c.typ == pb.ConverterType_OVERLAY_IMAGE {
		return c.modify_image(gctx)
	} else {
		return fmt.Errorf("cannot modify \"%v\" yet", c.typ)
	}
}

func (c *converter) modify_through_extbin(gctx *gg.Context) error {
	//TODO: figure out an input format for ext_binary
	h, w := c.cfg.ifp.GetDimensions()
	c.extbin.Call_Modify(nil, h, w)
	return nil
}

func (c *converter) webcam(gctx *gg.Context) error {
	frame, err := c.vcs.GetFrame()
	if err != nil {
		return err
	}
	if len(frame) == 0 {
		return nil
	}
	h, w := c.cfg.ifp.GetDimensions()
	img := converters.ConvertYUV422ToImage(frame, int(h), int(w))
	gctx.DrawImage(img, 0, 0)

	return nil
}

func (c *converter) modify_image(gctx *gg.Context) error {
	img := c.image
	if img == nil {
		return nil
	}
	gctx.DrawImage(img, int(c.image_x), int(c.image_y))
	//	fmt.Printf("Drawing image at %dx%d\n", c.image_x, c.image_y)
	return nil
}
func (c *converter) modify_text(gctx *gg.Context) error {
	c.tmv.Step()
	//	ifp := c.cfg.ifp
	col := c.tmv.RGBA()
	txt := "."
	if c.text != nil {
		txt = c.text()
	}
	l := labeller.NewLabellerForCtx(gctx)
	l.SetFontSize(80)

	x := c.tmv.X()
	y := c.tmv.Y()
	ld := l.NewLabel(x, y, col, txt)
	err := l.PaintLabel(ld)
	if err != nil {
		return err
	}
	return nil
}
