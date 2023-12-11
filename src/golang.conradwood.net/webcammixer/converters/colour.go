package converters

import (
	pb "golang.conradwood.net/apis/webcammixer"
	"image/color"
)

func ImgColourFromProto(col *pb.Colour) color.Color {
	return color.RGBA{uint8(col.Red), uint8(col.Green), uint8(col.Blue), 255}

}
