package labeller

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	// "strings"
	"flag"
	//	pb "golang.conradwood.net/apis/webcammixer"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	debug = flag.Bool("debug_labeller", false, "debug mode for labeller")
)

/*
typical workflow:
1) create labeller (either blank canvas or from an image)
2) call PaintLabels
3) call GetImage()

Note: at the end of PaintLabels an image is converted. Prefer to set all labels in one call to PaintLabels instead of calling it once for each Label
*/

const (
	// apt install texlive-fonts-extra-links
	DEFAULT_FONTNAME = "dejavu/DejaVuSans.ttf"
	DEFAULT_FONTSIZE = 25
)

var (
	FONT_DIRS = []string{
		"/usr/share/texlive/texmf-dist/fonts/truetype/public",
		"/usr/share/fonts/truetype/",
	}
)

type LabelDef struct {
	xpos     int
	ypos     int
	colour   *LabelColour
	fontname string
	fontsize uint32
	text     string
}
type LabelColour struct {
	Red   int
	Green int
	Blue  int
}
type Labeller struct {
	img      image.Image
	fontname string
	fontsize uint32
	ctx      *gg.Context
	height   int
	width    int
}

func NewLabellerForCtx(ctx *gg.Context) *Labeller {
	res := &Labeller{ctx: ctx}
	res.set_defaults()
	return res
}
func NewLabellerForImage(img image.Image) *Labeller {
	dc := gg.NewContextForImage(img)
	res := NewLabellerForCtx(dc)
	res.img = img
	return res
}
func NewLabellerForBlankCanvas(xsize, ysize int, col color.RGBA) *Labeller {
	img := image.NewRGBA(image.Rect(0, 0, xsize, ysize))
	draw.Draw(img, img.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
	return NewLabellerForImage(img)
}
func (l *Labeller) set_defaults() {
	if l.fontname == "" {
		l.fontname = DEFAULT_FONTNAME
	}
	if l.fontsize == 0 {
		l.fontsize = DEFAULT_FONTSIZE
	}
}
func (l *Labeller) SetFontSize(fontsize uint32) {
	l.fontsize = fontsize
}
func (l *Labeller) NewLabel(x, y int, col color.RGBA, text string) *LabelDef {
	res := &LabelDef{
		xpos:   x,
		ypos:   y,
		colour: newLabelColourFromRGBA(col),
		text:   text,
	}
	return res
}

func (l *Labeller) GetImage() image.Image {
	return l.img
}

// height,width, available AFTER paintlabels was called
func (l *Labeller) GetMaxDimensions() (int, int) {
	return l.height, l.width
}

// paint single label on to current image
func (l *Labeller) PaintLabel(lt *LabelDef) error {
	return l.PaintLabels([]*LabelDef{lt})
}

// paint multiple labels onto current image (more efficient than single label)
func (l *Labeller) PaintLabels(lt []*LabelDef) error {
	dc := l.ctx
	var w, h float64
	for _, lab := range lt {

		fs := lab.fontsize
		if fs == 0 {
			fs = l.fontsize
		}
		fn := lab.fontname
		if fn == "" {
			fn = l.fontname
		}
		err := LoadFont(dc, fn, fs)
		if err != nil {
			return err
		}
		line := lab.text
		x := float64(lab.xpos)
		y := float64(lab.ypos)
		if *debug {
			fmt.Printf("Setting text to \"%s\" @ %0.1fx%0.1f\n", line, x, y)
		}
		col := lab.colour
		if col == nil {
			dc.SetRGB255(0, 0, 0)
		} else {
			dc.SetRGB255(col.Red, col.Green, col.Blue)
		}
		//                      dc.SetRGB255(255, 255, 255)
		dc.DrawString(line, x, y)
		cw, ch := dc.MeasureString(line)
		if cw > w {
			w = cw
		}
		if ch > h {
			h = ch
		}

	}
	l.width = int(w)
	l.height = int(h)
	l.img = dc.Image()

	return nil
}

func LoadFont(dc *gg.Context, fontname string, fontsize uint32) error {
	var err error
	if fontname == "" {
		fontname = DEFAULT_FONTNAME
	}
	if fontsize == 0 {
		fontsize = DEFAULT_FONTSIZE
	}
	for _, dir := range FONT_DIRS {
		err = dc.LoadFontFace(dir+"/"+fontname, float64(fontsize))
		if err == nil {
			return nil
		}
	}
	if err != nil {
		return err
	}
	return nil
}
func newLabelColourFromRGBA(col color.RGBA) *LabelColour {
	r, g, b, _ := col.RGBA()
	res := &LabelColour{
		Red:   int(r),
		Green: int(g),
		Blue:  int(b),
	}
	return res
}
