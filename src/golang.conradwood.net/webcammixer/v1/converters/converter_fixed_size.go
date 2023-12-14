package converters

// fixed size converter crops image or adds blank area (does not scale!)
type FixedSize struct {
	TargetWidth  uint32
	TargetHeight uint32
}

func (f *FixedSize) Convert(ri *RawImage) (*RawImage, error) {
	return ri, nil
}
