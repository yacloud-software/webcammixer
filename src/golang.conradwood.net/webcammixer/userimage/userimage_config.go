package userimage

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	pb "golang.conradwood.net/apis/webcammixer"
	// "golang.conradwood.net/webcammixer/labeller"
	"image"
)

const (
	DEFAULT_BINARY = "/usr/local/bin/webcammixer_ext_binary"
)

var (
	current_config *config
)

type ext_converter interface {
	Modify(gg.Context) error
}
type ImageSource interface {
	GetTimingChannel() chan bool
	GetFrame() ([]byte, error)
	Close()
	Activate()
}
type config struct {
	image_source ImageSource

	convdef    string
	ifp        *UserImageProvider
	req        *pb.UserImageRequest
	converters []*converter
}

func (ifp *UserImageProvider) SetConfig(cfg *pb.UserImageRequest) error {
	res := &config{req: cfg, ifp: ifp}
	if cfg.ImageSource == nil {
		return fmt.Errorf("No source specified")
	}
	h, w := ifp.GetDimensions()

	/* configure the converters */
	check_dups := make(map[string]bool) // check for duplicate references
	for _, cv := range cfg.Converters {
		_, exists := check_dups[cv.Reference]
		if exists {
			return fmt.Errorf("Duplicate reference \"%s\"", cv.Reference)
		}
		check_dups[cv.Reference] = true
		fmt.Printf("Config: Adding \"%v\"\n", cv.Type)
		c := &converter{
			input_has_changed: true,
			convdef:           cv,
			cfg:               res,
			typ:               cv.Type,
			tmv:               NewTextMover(h, w),
		}

		c.tmv.red = 60
		c.tmv.green = 145
		c.tmv.blue = 55

		if cv.Type == pb.ConverterType_LABEL {
			//			c.lab = labeller.NewLabellerForCfg(cv)
			s := cv.Text
			fmt.Printf("Adding label with text \"%s\"\n", s)
			c.text = func() string { return s }
		} else if cv.Type == pb.ConverterType_EMOJI {
			c.emoji = NewEmojiConverter(c)
			//
		} else if cv.Type == pb.ConverterType_EXT_BINARY {
			// needs no config
			ebin, err := GetExtBinary(DEFAULT_BINARY)
			if err != nil {
				return err
			}
			c.extbin = ebin
		} else if cv.Type == pb.ConverterType_OVERLAY_IMAGE {
			// an image handler
			oir := cv.OverlayImage
			if oir != nil {
				b := bytes.NewReader(oir.Image)
				img, _, err := image.Decode(b)
				if err != nil {
					return err
				}
				c.image = img
				c.image_x = oir.XPos
				c.image_y = oir.YPos
			}

		} else if cv.Type == pb.ConverterType_WEBCAM {
			/*
				h, w := ifp.GetDimensions()
				vcs, err := res.ifp.sourceMixer.SourceActivateVideoDef(cv.Device.Device, h, w)
				if err != nil {
					return err
				}
				c.vcs = vcs
				tc := make(chan bool)
				vcs.SetTimerTarget(tc)
				res.ifp.timer_source = tc
			*/
			return fmt.Errorf("cannot yet use webcam as converter")
		} else {
			return fmt.Errorf("unsupported type \"%v\"", cv.Type)
		}

		res.converters = append(res.converters, c)
	}

	/* configure the source */
	if cfg.ImageSource.Device != nil {
		cv := cfg.ImageSource.Device
		vcs, err := NewWebcamSource(ifp.sourceMixer, h, w, cv.Device)
		if err != nil {
			return err
		}
		res.image_source = vcs
	} else if cfg.ImageSource.FillColour != nil {
		vcs, err := NewColourSource(h, w, cfg.ImageSource.FillColour)
		if err != nil {
			return err
		}
		res.image_source = vcs
	} else {
		return fmt.Errorf("Unsupported image source")
	}

	current_config = res
	res.ifp.NewSource(res.image_source)
	return nil
}
func (ifp *UserImageProvider) ConverterByReference(ref string) *converter {
	for _, c := range current_config.converters {
		if c.convdef.Reference == ref {
			return c
		}
	}
	return nil
}
func (ifp *UserImageProvider) SetText(s func() string) {
	if current_config == nil {
		return
	}
	for _, c := range current_config.converters {
		if c.typ == pb.ConverterType_LABEL {
			c.text = s
			c.input_has_changed = true
		}
	}
}
func (ifp *UserImageProvider) SetImage(x, y uint32, image image.Image) {
	if current_config == nil {
		return
	}
	for _, c := range current_config.converters {
		if c.typ == pb.ConverterType_OVERLAY_IMAGE {
			c.image = image
			c.image_x = x
			c.image_y = y
		}

		c.input_has_changed = true
	}
}
