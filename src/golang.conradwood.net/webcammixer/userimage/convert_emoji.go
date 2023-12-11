package userimage

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"golang.conradwood.net/apis/images"
	//	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	//	"golang.conradwood.net/go-easyops/utils"
	//	"golang.conradwood.net/webcammixer/converters"
	//	"golang.conradwood.net/webcammixer/labeller"
	"image"
	// "image/draw"
)

type emoji_converter struct {
	emoji_utf8  string
	image       image.Image
	gotimage    bool
	has_changed bool
}

func NewEmojiConverter(c *converter) *emoji_converter {
	res := &emoji_converter{
		has_changed: true,
		emoji_utf8:  c.convdef.Emoji,
	}
	return res
}

func (c *emoji_converter) Modify(gctx *gg.Context) (bool, error) {
	if !c.gotimage {
		img, err := get_emoji(c.emoji_utf8)
		if err != nil {
			return false, err
		}
		c.image = img
		c.gotimage = true
	}
	gctx.DrawImage(c.image, 100, 100)
	c.has_changed = false
	return true, nil
}
func (c *emoji_converter) HasChanged() bool {
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

	img, _, err := image.Decode(bytes.NewReader(em.Emoji.PNG))
	if err != nil {
		return nil, err
	}
	fmt.Printf("emoji loaded\n")
	return img, nil
}
