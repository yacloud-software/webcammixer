package interfaces

import (
	"golang.conradwood.net/webcammixer/loopback"
	"golang.conradwood.net/webcammixer/providers"
)

type MixerApp interface {
	DefaultIdleFrameProvider() *providers.IdleFrameProvider
	GetLoopDev() *loopback.LoopBackDevice
	Start() error
	SetSwitcher(Switcher)
	GetCurrentProvider() loopback.FrameProvider
}
type Switcher interface {
	DeactivateUserFrames() error
	ActivateUserFrames() error
	SetMixerApp(ma MixerApp)
	SetSourceMixer(SourceMixer)
}
type SourceMixer interface {
	SourceActivateVideoDef(devicename string, height, width uint32) (VideoCamSource, error)
}
type VideoCamSource interface {
	GetFrame() ([]byte, error)
	GetID() string
	SetTimerTarget(c chan bool) error
}
