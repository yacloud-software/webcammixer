package userimage

import (
	"github.com/fogleman/gg"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/webcammixer/converters"
	"golang.conradwood.net/webcammixer/labeller"
	"image/color"
)

func (c *converter) Modify(gctx *gg.Context) error {
	if c.typ == pb.ConverterType_LABEL {
		return c.modify_text(gctx)
	}
	return nil
}

func (c *converter) modify_text(gctx *gg.Context) error {
	ifp := c.cfg.ifp
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
	ifp.idleImage = l.GetImage()
	rawimage, err := converters.ConvertToRaw(ifp.idleImage)
	if err != nil {
		return err
	}
	ifp.idleImageRaw = rawimage.DefaultBytes()
	ifp.idleImageLastFile = ""
	return nil
}
