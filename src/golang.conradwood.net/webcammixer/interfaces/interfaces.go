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
}
type Switcher interface {
	DeactivateUserFrames() error
	ActivateUserFrames() error
	SetMixerApp(ma MixerApp)
}
