package mixerapp

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/v1/defaults"
	"golang.conradwood.net/webcammixer/v1/interfaces"
	"golang.conradwood.net/webcammixer/v1/loopback"
	"golang.conradwood.net/webcammixer/v1/providers"
	"golang.conradwood.net/webcammixer/v1/webcam"
	"strings"
	"sync"
	"time"
)

var (
	open_close_hook = flag.String("usage_hook", "", "if set to a `filename` will invoke that executable with a single argument (true|false) if videodev changes from being used to not being used and vice versa")
	is_running_lock sync.RWMutex
	loopdev         *loopback.LoopBackDevice
	defidle         *providers.IdleFrameProvider
)

type MixerApp struct {
	switcher interfaces.Switcher
}

func (ma *MixerApp) SetSwitcher(sw interfaces.Switcher) {
	ma.switcher = sw
}

func (ma *MixerApp) Start() error {
	fmt.Printf("Starting webcammixer app...\n")
	is_running_lock.Lock()
	defer is_running_lock.Unlock()
	wl, err := webcam.Detect()
	if err != nil {
		fmt.Printf("Failed: %s\n", err)
		return err
	}
	var wi_loop *webcam.WebCamInfo
	t := utils.Table{}
	t.AddHeaders("Device", "Capture", "Driver", "Card", "Bus")
	for _, w := range wl {
		t.AddString(w.DeviceName)

		s := "no"
		if w.IsCaptureDevice() {
			s = "yes"
		}
		t.AddString(s)

		if w.OpenErr != nil {
			t.AddString(fmt.Sprintf(" Error: %s\n", w.OpenErr))
			t.NewRow()
			continue
		}
		cap := w.Capabilities
		t.AddString(cap.Driver)
		t.AddString(cap.Card)
		t.AddString(cap.BusInfo)

		//		t.AddString(fmt.Sprintf("%v %v", w.Capabilities.Capabilities, w.Capabilities.DeviceCapabilities))
		t.NewRow()
		if strings.Contains(cap.Driver, "loopback") {
			wi_loop = w
		}
	}
	fmt.Println(t.ToPrettyString())
	if wi_loop == nil {
		return fmt.Errorf("No loopback device found.\n")
	}

	fmt.Printf("Using v4l2loopback device: %s\n", wi_loop.DeviceName)
	h, w := defaults.GetDimensions()
	loopdev, err = loopback.Open(wi_loop.DeviceName, h, w)
	if err != nil {
		return fmt.Errorf("failed to open loopback device: %w", err)
	}
	fmt.Printf("Loopback opened: %#v\n", wi_loop.DeviceName)
	defidle = providers.NewIdleFrameProvider(loopdev.GetDimensions())
	loopdev.SetProvider(defidle)
	loopdev.SetTimingSource(defidle)
	//	ifp.SetTimerTarget(loopdev.GetTimerChan())
	var wg sync.WaitGroup
	var terr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = loopdev.StartWriter()
		if err != nil {
			terr = err
		}
	}()

	go func() {
		defer wg.Done()
		err := defidle.Run() // blocks
		if err != nil {
			fmt.Printf("Failed to run idle provider: %s\n", err)
			terr = err
		}
	}()
	go ma.idle_detect_thread(loopdev)
	fmt.Printf("Started webcammixer app\n")
	// create albeit "incorrect" usage. we we add one to waitgroup, but call Done()
	// when any one (or both) threads return.
	// essentially we are waiting for AT LEAST ONE to return
	wg.Wait()
	if terr != nil {
		fmt.Printf("failure: %s\n", terr)
		return err
	}
	return nil
}

func (ma *MixerApp) WaitUntilDone() {
	is_running_lock.RLock()
	is_running_lock.RUnlock()

	fmt.Printf("Done.\n")
}
func (ma *MixerApp) GetCurrentProvider() loopback.FrameProvider {
	lp := ma.GetLoopDev()
	if lp == nil {
		return nil
	}
	return lp.GetFrameProvider()

}
func (ma *MixerApp) GetLoopDev() *loopback.LoopBackDevice {
	return loopdev
}

func (ma *MixerApp) idle_detect_thread(l *loopback.LoopBackDevice) {
	var last_restored time.Time
	last_ref_count := 0
	last_status := 0 // 0==undef, 1=open, 2=idle
	for {
		time.Sleep(time.Duration(100) * time.Millisecond)

		restore := false
		reason := ""
		if time.Since(l.TimeOfLastFrame()) > time.Duration(2)*time.Second {
			reason = "no more frames"
			restore = true
		}
		ls := loopback.Status()
		if ls.RefCount >= 2 {
			if last_status != 1 {
				ma.loopback_open_status_changed(true)
			}
			last_status = 1
			last_ref_count = ls.RefCount
		} else {
			if time.Since(ls.LastChange) > time.Duration(10)*time.Second {
				if last_status != 2 {
					ma.loopback_open_status_changed(false)
				}
				last_status = 2

				if last_ref_count != ls.RefCount {
					last_ref_count = ls.RefCount
					reason = "kernel module not in use"
					restore = true
				}
			}
		}

		if restore && time.Since(last_restored) > time.Duration(5)*time.Second {
			fmt.Printf("restoring idleframeprovider (%s)...\n", reason)
			ma.switcher.DeactivateUserFrames()
			ifp := ma.DefaultIdleFrameProvider()
			l.SetProvider(ifp)
			l.SetTimingSource(ifp)
			last_restored = time.Now()
		}
	}
}

// called if at  least one application uses it or if all apps stopped using it
func (ma *MixerApp) loopback_open_status_changed(is_used bool) {
	oc := *open_close_hook
	if oc == "" {
		return
	}
	go func(b bool) {
		l := linux.New()
		s := "false"
		if b {
			s = "true"
		}
		com := []string{oc, s}
		r, err := l.SafelyExecute(com, nil)
		if err != nil {
			fmt.Printf("executable %s failed: %s\n%s\n", oc, err, string(r))
		}
	}(is_used)
}

func (ma *MixerApp) DefaultIdleFrameProvider() *providers.IdleFrameProvider {
	return defidle
}
