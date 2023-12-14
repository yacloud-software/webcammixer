package defaults

import (
	"flag"
	"fmt"
	"github.com/vladimirvivien/go4vl/v4l2"
	"strconv"
	"strings"
)

var (
	dimensions = flag.String("size", "1280x720", "image dimensions")
)

// return height,width
func GetDimensions() (uint32, uint32) {
	s := *dimensions
	sx := strings.Split(s, "x")
	if len(sx) != 2 {
		panic(fmt.Sprintf("dimensions not in format WIDTHxHEIGHT (%s)", s))
	}
	w, err := strconv.Atoi(sx[0])
	if err != nil {
		panic(fmt.Sprintf("invalid width specified: %s", err))
	}
	h, err := strconv.Atoi(sx[1])
	if err != nil {
		panic(fmt.Sprintf("invalid height specified: %s", err))
	}
	return uint32(h), uint32(w)
}

// format, multiplier to get at size
func GetPreferredFormat() (uint32, uint32) {
	//	return v4l2.PixelFmtRGB24, 3
	return v4l2.PixelFmtYUYV, 2
}
