package userimage

import (
	"github.com/fogleman/gg"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/webcammixer/labeller"
)

var (
	current_config *config
)

type ext_converter interface {
	Modify(gg.Context) error
}

type config struct {
	ifp        *UserImageProvider
	req        *pb.UserImageRequest
	converters []*converter
}

type converter struct {
	cfg  *config
	typ  pb.ConverterType
	lab  *labeller.Labeller // implements ext_converter
	text func() string
}

func (ifp *UserImageProvider) SetConfig(cfg *pb.UserImageRequest) error {
	res := &config{req: cfg, ifp: ifp}
	for _, cv := range cfg.Converters {
		c := &converter{
			cfg: res,
			typ: cv.Type,
		}
		if cv.Type == pb.ConverterType_LABEL {
			c.lab = labeller.NewLabellerForCfg(cv)
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
