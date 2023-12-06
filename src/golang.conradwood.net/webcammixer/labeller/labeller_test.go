package labeller

import (
	"bytes"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"image/color"
	"image/png"
	"testing"
)

const (
	FILEPREFIX = "/tmp/testpix/img"
)

var (
	i            = 0
	std_labeldef = &LabelDef{colour: &LabelColour{Red: 0xff, Green: 255, Blue: 0xff}, xpos: 20, ypos: 20, text: "foobar"}
)

func TestCreateLabel(t *testing.T) {
	err := create_image(t, 200, 200, "foobar", 20, 20)
	if err != nil {
		t.Logf("failed: %s\n", err)
	}
	create_image(t, 200, 200, "X", 0, 0)
	create_image(t, 200, 200, "X", 200, 0)
	create_image(t, 200, 200, "X", 0, 200)
	create_image(t, 200, 200, "X", 200, 200)
}
func create_image(t *testing.T, x, y int, text string, xpos, ypos int) error {
	l := NewLabellerForBlankCanvas(x, y, color.RGBA{0, 0, 0, 255})
	std_labeldef.text = text
	std_labeldef.xpos = xpos
	std_labeldef.ypos = ypos
	err := l.PaintLabels([]*LabelDef{std_labeldef})
	if err != nil {
		t.Fatalf("paint failed: %s", err)
		return err
	}
	i++
	fc := i
	filename := fmt.Sprintf("%s_%d.png", FILEPREFIX, fc)
	w := &bytes.Buffer{}
	err = png.Encode(w, l.GetImage())
	if err != nil {
		t.Fatalf("encode failed: %s", err)
		return err
	}
	b := w.Bytes()
	err = utils.WriteFile(filename, b)
	if err != nil {
		t.Fatalf("write failed: %s", err)
		return err
	}
	t.Logf("Written to %s\n", filename)
	return nil
}
