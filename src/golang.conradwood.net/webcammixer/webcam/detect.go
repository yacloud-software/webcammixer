package webcam

import (
	"fmt"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
)

var (
	webcam_lock       sync.Mutex
	known_webcam_info = make(map[string]*WebCamInfo)
)

type WebCamInfo struct {
	DeviceName        string
	Capabilities      v4l2.Capability
	OpenErr           error // if an error occured when opening
	device            *device.Device
	cancel            func()         // context cancel (on stop streaming)
	actualPixelFormat v4l2.PixFormat // once opened, this is the current and actual format
}

func Detect() ([]*WebCamInfo, error) {
	files, err := ioutil.ReadDir("/dev")
	if err != nil {
		return nil, err
	}
	var res []*WebCamInfo
	var wg sync.WaitGroup
	var rl sync.Mutex
	for _, df := range files {
		dfn := "/dev/" + df.Name()
		if !strings.HasPrefix(dfn, "/dev/video") {
			continue
		}
		wg.Add(1)
		go func(fn string) {
			defer wg.Done()
			//		fmt.Printf("f:%s\n", fn)
			wci := GetWebCamInfo(fn)
			rl.Lock()
			res = append(res, wci)
			rl.Unlock()
			dev, err := device.Open(wci.DeviceName,
				device.WithIOType(v4l2.IOTypeMMAP),
			)
			if err != nil {
				wci.OpenErr = err
				fmt.Printf("Failed to open %s: %s\n", wci.DeviceName, wci.OpenErr)
				return
			}
			wci.Capabilities = dev.Capability()
			dev.Close()
		}(dfn)
		//		fmt.Printf("%s Card: %s\n", wci.DeviceName, wci.Capabilities.Card)
	}
	wg.Wait()
	sort.Slice(res, func(i, j int) bool {
		return res[i].DeviceName < res[j].DeviceName
	})

	return res, nil
}

// this is just horrible...
func (w *WebCamInfo) IsCaptureDevice() bool {
	for _, cd := range w.Capabilities.GetDeviceCapDescriptions() {
		s := strings.ToLower(fmt.Sprintf("%v", cd))
		if strings.Contains(s, "video capture") {
			return true
		}
	}
	return false
}
func GetWebCamInfo(devicename string) *WebCamInfo {
	webcam_lock.Lock()
	defer webcam_lock.Unlock()
	wci, found := known_webcam_info[devicename]
	if !found {
		wci = &WebCamInfo{DeviceName: devicename}
		known_webcam_info[devicename] = wci
	}
	return wci

}

func GetCaptureDevices() ([]*WebCamInfo, error) {
	wl, err := Detect()
	if err != nil {
		fmt.Printf("Failed: %s\n", err)
		return nil, err
	}
	var res []*WebCamInfo
	for _, w := range wl {
		if !w.IsCaptureDevice() {
			continue
		}
		if w.OpenErr != nil {
			continue
		}
		xcap := w.Capabilities
		if strings.Contains(xcap.Driver, "loopback") {
			continue
		}
		res = append(res, w)
	}
	return res, nil

}
