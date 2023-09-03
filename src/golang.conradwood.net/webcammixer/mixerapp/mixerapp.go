package mixerapp

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/defaults"
	"golang.conradwood.net/webcammixer/loopback"
	"golang.conradwood.net/webcammixer/webcam"
	"strings"
	"sync"
	"time"
)

var (
	is_running_lock sync.RWMutex
	loopdev         *loopback.LoopBackDevice
	defidle         *IdleFrameProvider
)

func Start() error {
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
	defidle = NewIdleFrameProvider(loopdev.GetDimensions())
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
	go idle_detect_thread(loopdev)
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

func WaitUntilDone() {
	is_running_lock.RLock()
	is_running_lock.RUnlock()

	fmt.Printf("Done.\n")
}

func GetLoopDev() *loopback.LoopBackDevice {
	return loopdev
}

func idle_detect_thread(l *loopback.LoopBackDevice) {
	var last_restored time.Time
	for {
		time.Sleep(time.Duration(100) * time.Millisecond)
		if time.Since(l.TimeOfLastFrame()) > time.Duration(2)*time.Second {
			if time.Since(last_restored) > time.Duration(5)*time.Second {
				fmt.Printf("restoring idleframeprovider...\n")
				ifp := DefaultIdleFrameProvider()
				l.SetProvider(ifp)
				l.SetTimingSource(ifp)
				last_restored = time.Now()
			}
		}
	}
}

func DefaultIdleFrameProvider() *IdleFrameProvider {
	return defidle
}
