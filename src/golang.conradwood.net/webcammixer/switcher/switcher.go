package switcher

import (
	"fmt"
	"golang.conradwood.net/webcammixer/interfaces"
	"golang.conradwood.net/webcammixer/userimage"
	"sync"
)

var (
	userframe_active = false
	mixerapp         interfaces.MixerApp
	getuserimagelock sync.Mutex
	def_user_image   *userimage.UserImageProvider
)

type switcher struct {
	source_mixer interfaces.SourceMixer
}

func NewSwitcher() *switcher {
	return &switcher{}
}
func (sw *switcher) SetSourceMixer(s interfaces.SourceMixer) {
	sw.source_mixer = s
}
func (sw *switcher) SetMixerApp(ma interfaces.MixerApp) {
	mixerapp = ma
}
func (sw *switcher) GetCurrentUserImageProvider() *userimage.UserImageProvider {
	ifp := mixerapp.DefaultIdleFrameProvider()
	w := uint32(640)
	h := uint32(480)
	if ifp != nil {
		h, w = ifp.GetDimensions()
	}
	return sw.get_user_image_provider(w, h)
}

func (sw *switcher) get_user_image_provider(x, y uint32) *userimage.UserImageProvider {
	getuserimagelock.Lock()
	defer getuserimagelock.Unlock()
	if def_user_image != nil {
		h, w := def_user_image.GetDimensions()
		if w == x && h == y {
			return def_user_image
		}
		def_user_image.Stop()
		fmt.Printf("Stopping userimage provide with dimensions %d x %d\n", w, h)
	}
	ndef := userimage.NewUserImageProvider(sw.source_mixer, y, x)
	go ndef.Run()
	def_user_image = ndef
	fmt.Printf("Started userimage provide with dimensions %d x %d\n", x, y)
	return def_user_image
}

func (sw *switcher) IsUserFrameActive() bool {
	return userframe_active
}

func (sw *switcher) ActivateUserFrames() error {
	mfp := sw.GetCurrentUserImageProvider()
	loopdev := mixerapp.GetLoopDev()
	if loopdev == nil {
		return fmt.Errorf("(1) not ready yet - try again later")
	}
	loopdev.SetProvider(mfp)
	loopdev.SetTimingSource(mfp)
	userframe_active = true
	return nil
}
func (sw *switcher) DeactivateUserFrames() error {
	getuserimagelock.Lock()
	defer getuserimagelock.Unlock()
	if def_user_image == nil {
		return nil
	}
	def_user_image.Stop()
	def_user_image = nil
	userframe_active = false
	return nil
}
