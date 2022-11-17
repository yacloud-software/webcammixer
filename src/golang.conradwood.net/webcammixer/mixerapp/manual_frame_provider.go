package mixerapp

type ManualFrameProvider struct {
	lastFrame []byte
	notify    chan bool
}

func NewManualFrameProvider() *ManualFrameProvider {
	return &ManualFrameProvider{}
}
func (mfp *ManualFrameProvider) GetID() string {
	return "manualframeprovider"
}
func (mfp *ManualFrameProvider) GetFrame() ([]byte, error) {
	return mfp.lastFrame, nil
}
func (mfp *ManualFrameProvider) NewFrame(frame []byte) {
	mfp.lastFrame = frame
	c := mfp.notify
	if c != nil {
		c <- true
	}
}
func (mfp *ManualFrameProvider) SetTimerTarget(c chan bool) error {
	mfp.notify = c
	return nil
}
