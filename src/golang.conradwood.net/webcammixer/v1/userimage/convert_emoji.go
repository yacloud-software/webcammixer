package userimage

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"golang.conradwood.net/apis/images"
	//	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	//	"golang.conradwood.net/go-easyops/utils"
	//	"golang.conradwood.net/webcammixer/v1/converters"
	//	"golang.conradwood.net/webcammixer/v1/labeller"
	"image"
	// "image/draw"
	"flag"
)

var (
	emoji_speed = flag.Int("emoji_speed", 7, "number of pixels to move per frame")
)

type emoji_converter struct {
	c           *converter
	emoji_utf8  string
	image       image.Image
	gotimage    bool
	has_changed bool
	curX        int
	curY        int
	finished    bool
}

func NewEmojiConverter(c *converter) *emoji_converter {
	res := &emoji_converter{
		c:           c,
		has_changed: true,
		emoji_utf8:  c.convdef.Emoji,
	}
	return res
}

func (c *emoji_converter) Modify(gctx *gg.Context) (bool, error) {
	if c.finished {
		return false, nil
	}
	if !c.gotimage {
		img, err := get_emoji(c.emoji_utf8)
		if err != nil {
			return false, err
		}
		c.image = img
		c.gotimage = true
	}
	gctx.DrawImage(c.image, c.curX, 100)
	c.has_changed = true
	return true, nil
}
func (c *emoji_converter) HasChanged() bool {
	if c.finished {
		return false
	}
	c.curX = c.curX + *emoji_speed
	_, w := c.c.cfg.ifp.GetDimensions()
	if c.curX >= int(w) {
		c.finished = true
	}
	return c.has_changed
}

func get_emoji(utf8 string) (image.Image, error) {
	sr := &images.EmojiSearchRequest{Text: utf8}
	ctx := authremote.Context()
	em, err := images.GetImagesClient().SearchEmojis(ctx, sr)
	if err != nil {
		return nil, err
	}
	if len(em.Defs) == 0 {
		return nil, fmt.Errorf("Emoji \"%s\" not found", utf8)
	}
	if em.Emoji == nil {
		return nil, fmt.Errorf("No Emoji png for \"%s\" (%d results)", utf8, len(em.Defs))
	}
	png_data := em.Emoji.PNG

	ge := &images.EmojiRequest{Unicode: em.Emoji.Def.Unicode, Size: 80}
	emm, err := images.GetImagesClient().GetEmoji(ctx, ge)
	if err != nil {
		fmt.Printf("Failed to get emoji, size %d\n", ge.Size)
	} else {
		png_data = emm.PNG
	}
	img, _, err := image.Decode(bytes.NewReader(png_data))
	if err != nil {
		return nil, err
	}
	fmt.Printf("emoji loaded\n")
	return img, nil
}
