package loopback

import (
	"fmt"
)

func (l *LoopBackDevice) AddProvider(f FrameProvider) {
	l.frameProviderLock.Lock()
	l.frameProviders[f.GetID()] = f
	l.frameProviderLock.Unlock()
}

func (l *LoopBackDevice) RemoveProvider(f FrameProvider) {
	l.frameProviderLock.Lock()
	delete(l.frameProviders, f.GetID())
	l.frameProviderLock.Unlock()
}

// set to a single frame provider
func (l *LoopBackDevice) SetProvider(f FrameProvider) {
	l.frameProviderLock.Lock()
	res := make(map[string]FrameProvider)
	res[f.GetID()] = f
	l.frameProviders = res
	l.frameProviderLock.Unlock()
}

// every time this is called a new frame will be sent to the video device
func (l *LoopBackDevice) NewFrame() error {
	var provs []FrameProvider
	l.frameProviderLock.Lock()
	for _, v := range l.frameProviders {
		provs = append(provs, v)
	}
	l.frameProviderLock.Unlock()
	if len(provs) == 0 {
		fmt.Printf("No frame providers\n")
		return nil
	}
	// TODO: actually multiplex them here
	fp := provs[0]
	frame, err := fp.GetFrame()
	if err != nil {
		return err
	}
	fmt.Printf("loopback: Got new frame (%d bytes). (using %s)\n", len(frame), fp.GetID())
	if uint32(len(frame)) != l.bufsize {
		fmt.Printf("Frame is %d bytes, but bufsize for loopback device is %d bytes!\n", len(frame), l.bufsize)
	}
	err = l.WriteImage(frame)
	if err != nil {
		return err
	}
	return nil

}
