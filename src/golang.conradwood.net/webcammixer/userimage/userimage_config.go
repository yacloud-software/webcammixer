package userimage

import (
	"github.com/fogleman/gg"
	pb "golang.conradwood.net/apis/webcammixer"
	// "golang.conradwood.net/webcammixer/labeller"
)

var (
	current_config *config
)

type ext_converter interface {
	Modify(gg.Context) error
}

type config struct {
	convdef    string
	ifp        *UserImageProvider
	req        *pb.UserImageRequest
	converters []*converter
}

func (ifp *UserImageProvider) SetConfig(cfg *pb.UserImageRequest) error {
	res := &config{req: cfg, ifp: ifp}
	for _, cv := range cfg.Converters {
		c := &converter{
			convdef: cv,
			cfg:     res,
			typ:     cv.Type,
			tmv: &text_mover{
				red:   60,
				green: 145,
				blue:  55,
			},
		}
		if cv.Type == pb.ConverterType_LABEL {
			//			c.lab = labeller.NewLabellerForCfg(cv)
			c.text = func() string { return cv.Text }
		}
		if cv.Type == pb.ConverterType_WEBCAM {
			h, w := ifp.GetDimensions()
			vcs, err := res.ifp.sourceMixer.SourceActivateVideoDef(cv.Device.Device, h, w)
			if err != nil {
				return err
			}
			c.vcs = vcs
			tc := make(chan bool)
			vcs.SetTimerTarget(tc)
			res.ifp.timer_source = tc
		}

		res.converters = append(res.converters, c)
	}
	current_config = res
	return nil
}
func (ifp *UserImageProvider) SetText(s func() string) {
	if current_config == nil {
		return
	}
	for _, c := range current_config.converters {
		if c.typ == pb.ConverterType_LABEL {
			c.text = s
		}
	}
}
