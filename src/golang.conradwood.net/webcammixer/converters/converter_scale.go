package converters

type Scaler struct {
	TargetWidth  uint32
	TargetHeight uint32
}

func (f *Scaler) Convert(ri *RawImage) (*RawImage, error) {
	return ri, nil
}
