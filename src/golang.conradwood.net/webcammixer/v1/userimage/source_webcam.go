package userimage

import (
	"fmt"
	"golang.conradwood.net/webcammixer/v1/interfaces"
)

type webcamsource struct {
	vd      interfaces.VideoCamSource
	dev     string
	vd_chan chan bool
}

func NewWebcamSource(sourceMixer interfaces.SourceMixer, h, w uint32, dev string) (*webcamsource, error) {
	res := &webcamsource{vd_chan: make(chan bool, 50), dev: dev}
	vd, err := sourceMixer.SourceActivateVideoDef(dev, h, w)
	if err != nil {
		return nil, err
	}
	res.vd = vd
	res.vd.SetTimerTarget(res.vd_chan)
	fmt.Printf("Webcam set timer target\n")
	return res, nil
}
func (src *webcamsource) Activate() {
	src.vd.SetTimerTarget(src.vd_chan)
}
func (src *webcamsource) GetFrame() ([]byte, error) {
	return src.vd.GetFrame()
}
func (src *webcamsource) GetTimingChannel() chan bool {
	return src.vd_chan
}
func (src *webcamsource) Close() {
	fmt.Printf("Webcamsource closed\n")
	src.vd.SetTimerTarget(nil)
}
func (src *webcamsource) String() string {
	return fmt.Sprintf("webcam://%s", src.dev)
}
