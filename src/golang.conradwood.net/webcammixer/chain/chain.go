package chain

import (
	pb "golang.conradwood.net/apis/webcammixer"
)

/*
a chain has an image source (e.g. e.g. a USB webcam, a video stream or a stream of produced images (e.g. a "black image")
a chain has zero or more converters to do something to the image from the source
a chain may send frames to an output.
*/
type Chain interface {
	// configure the current converters
	Configure(*pb.ChainConfig) error
	// get one of the current converters
	ConverterByReference(ref string) *pb.ChainConverter
	// get current Source
	GetCurrentSource() Source
	// where to send the output to
	SetOutput(Output)
}
type Source interface {
}
type Output interface {
	NewFrame()
}
