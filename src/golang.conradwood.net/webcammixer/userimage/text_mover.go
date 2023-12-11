package userimage

import (
	"image/color"
)

const ()

type text_mover struct {
	min_xpos    int
	max_xpos    int
	min_ypos    int
	max_ypos    int
	red         int
	green       int
	blue        int
	cur_red     int
	cur_green   int
	cur_blue    int
	xpos        int
	ypos        int
	col_up      bool
	xpos_up     bool
	ypos_up     bool
	text_height int
	text_width  int
}

func NewTextMover(h, w uint32) *text_mover {
	hi := int(h)
	wi := int(w)
	w_diff := int(float32(w) * 0.2)
	h_diff := int(float32(h) * 0.2)
	t := &text_mover{
		min_xpos: w_diff,
		max_xpos: wi - w_diff,
		min_ypos: h_diff,
		max_ypos: hi - h_diff,
	}

	return t
}
func (t *text_mover) Step() {
	if t.col_up {
		if t.cur_red >= 254 || t.cur_blue >= 254 || t.cur_green >= 254 {
			t.col_up = false
		} else {
			t.cur_red++
			t.cur_blue++
			t.cur_green++
		}
	} else {
		if t.cur_red <= t.red || t.cur_blue <= t.red || t.cur_green <= t.red {
			t.col_up = true
			t.cur_red = t.red
			t.cur_blue = t.blue
			t.cur_green = t.green
		} else {
			t.cur_red--
			t.cur_blue--
			t.cur_green--
		}
	}

	max_xpos := t.max_xpos - t.text_width
	max_ypos := t.max_ypos - t.text_height
	if t.xpos <= t.min_xpos {
		t.xpos_up = true
		t.xpos = t.min_xpos
	}
	if t.xpos >= max_xpos {
		t.xpos_up = false
		t.xpos = max_xpos
	}
	if t.ypos <= t.min_ypos {
		t.ypos_up = true
		t.ypos = t.min_ypos
	}
	if t.ypos >= max_ypos {
		t.ypos_up = false
		t.ypos = max_ypos
	}
	if t.xpos_up {
		t.xpos++
	} else {
		t.xpos--
	}
	if t.ypos_up {
		t.ypos++
	} else {
		t.ypos--
	}
}
func (t *text_mover) RGBA() color.RGBA {
	return color.RGBA{uint8(t.cur_red), uint8(t.cur_green), uint8(t.cur_blue), 255}
}
func (t *text_mover) X() int {
	return t.xpos
}
func (t *text_mover) Y() int {
	return t.ypos
}
