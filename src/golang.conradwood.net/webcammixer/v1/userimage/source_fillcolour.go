package userimage

import (
	"fmt"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/webcammixer/v1/converters"
	"image"
	"image/draw"
	"time"
)

type fillcoloursource struct {
	c     chan bool
	frame []byte
	stop  bool
	col   *pb.Colour
}

func NewColourSource(h, w uint32, colour *pb.Colour) (*fillcoloursource, error) {
	res := &fillcoloursource{c: make(chan bool), col: colour}
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	col := converters.ImgColourFromProto(colour)
	draw.Draw(img, img.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
	rawimg, err := converters.ConvertToRaw(img)
	if err != nil {
		return nil, err
	}
	res.frame = rawimg.DefaultBytes()
	return res, nil
}
func (src *fillcoloursource) timer_loop() {
	for !src.stop {
		time.Sleep(time.Duration(100) * time.Millisecond)
		src.c <- true
	}
}
func (src *fillcoloursource) Activate() {
	go src.timer_loop()
}
func (src *fillcoloursource) Close() {
	src.stop = true
}
func (src *fillcoloursource) GetFrame() ([]byte, error) {
	return src.frame, nil
}
func (src *fillcoloursource) GetTimingChannel() chan bool {
	return src.c
}
func (src *fillcoloursource) String() string {
	return fmt.Sprintf("fillcolour://%02X%02X%02X", src.col.Red, src.col.Green, src.col.Blue)
}
