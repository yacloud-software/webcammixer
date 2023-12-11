package userimage

import (
	"golang.conradwood.net/webcammixer/interfaces"
)

type webcamsource struct {
	vd      interfaces.VideoCamSource
	vd_chan chan bool
}

func NewWebcamSource(sourceMixer interfaces.SourceMixer, h, w uint32, dev string) (*webcamsource, error) {
	res := &webcamsource{vd_chan: make(chan bool)}
	vd, err := sourceMixer.SourceActivateVideoDef(dev, h, w)
	if err != nil {
		return nil, err
	}
	res.vd = vd
	res.vd.SetTimerTarget(res.vd_chan)
	return res, nil
}
func (src *webcamsource) GetFrame() ([]byte, error) {
	return src.vd.GetFrame()
}
func (src *webcamsource) GetTimingChannel() chan bool {
	return src.vd_chan
}
func (src *webcamsource) Close() {
	src.vd.SetTimerTarget(nil)
}
