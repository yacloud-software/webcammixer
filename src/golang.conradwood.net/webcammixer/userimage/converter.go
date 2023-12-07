package userimage

import (
	"github.com/fogleman/gg"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/labeller"
	"image/color"
)

type converter struct {
	cfg     *config
	convdef *pb.UserImageConverter
	typ     pb.ConverterType
	lab     *labeller.Labeller // implements ext_converter
	text    func() string
	tmv     *text_mover
}

func (c *converter) Modify(gctx *gg.Context) error {
	if c.typ == pb.ConverterType_LABEL {
		return c.modify_text(gctx)
	}
	return nil
}

func (c *converter) modify_text(gctx *gg.Context) error {
	c.tmv.Step()
	ifp := c.cfg.ifp
	col := c.tmv.RGBA()
	xsize := 640
	ysize := 480
	h, w := ifp.GetDimensions()
	xsize = int(w)
	ysize = int(h)
	txt := "."
	if c.text != nil {
		txt = c.text()
	}
	l := labeller.NewLabellerForBlankCanvas(xsize, ysize, color.RGBA{0, 0, 0, 255})
	l.SetFontSize(80)

	x := c.tmv.X()
	y := c.tmv.Y()
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
