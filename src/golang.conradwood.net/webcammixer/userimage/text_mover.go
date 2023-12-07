package userimage

import (
	"image/color"
)

const (
	MIN_XPOS = 100
	MAX_XPOS = 500
	MIN_YPOS = 100
	MAX_YPOS = 500
)

type text_mover struct {
	red       int
	green     int
	blue      int
	cur_red   int
	cur_green int
	cur_blue  int
	xpos      int
	ypos      int
	col_up    bool
	xpos_up   bool
	ypos_up   bool
}

func (t *text_mover) Step() {
	if t.col_up {
		if t.cur_red >= t.red || t.cur_blue >= t.blue || t.cur_green >= t.green {
			t.col_up = false
		}
		t.cur_red++
		t.cur_blue++
		t.cur_green++
	} else {
		if t.cur_red == 0 || t.cur_blue == 0 || t.cur_green == 0 {
			t.col_up = true
		}
		t.cur_red--
		t.cur_blue--
		t.cur_green--
	}

	if t.xpos <= MIN_XPOS {
		t.xpos_up = true
		t.xpos = MIN_XPOS
	}
	if t.xpos >= MAX_XPOS {
		t.xpos_up = false
		t.xpos = MAX_XPOS
	}
	if t.ypos <= MIN_YPOS {
		t.ypos_up = true
		t.ypos = MIN_YPOS
	}
	if t.ypos >= MAX_YPOS {
		t.ypos_up = false
		t.ypos = MAX_YPOS
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
