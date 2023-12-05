package main

import (
	"fmt"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/images"
	pb "golang.conradwood.net/apis/webcammixer"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"golang.conradwood.net/webcammixer/webcam"
	"sort"
	"time"
)

var (
	capture_devices *pb.CaptureDeviceList
	last_scanned    time.Time
)

func getCaptureDevices() (*pb.CaptureDeviceList, error) {
	if capture_devices == nil {
		return nil, fmt.Errorf("unavailable")
	}
	return capture_devices, nil
}
func cache_webcam_devices() {
	t := time.Duration(3) * time.Second
	for {
		time.Sleep(t)
		t = time.Duration(30) * time.Second
		cd, err := retrieve_capture_devices()
		if err != nil {
			fmt.Printf("failed to scan for capture devices: %s\n", err)
			t = time.Duration(5) * time.Second
			continue
		}
		capture_devices = cd
	}
}
func retrieve_capture_devices() (*pb.CaptureDeviceList, error) {
	ctx := authremote.Context()
	cameras, err := images.GetImagesClient().GetCameras(ctx, &common.Void{})
	if err != nil {
		fmt.Printf("error getting cameras: %s\n", utils.ErrorString(err))
		cameras = &images.CameraList{}
	}
	wlist, err := webcam.GetCaptureDevices()
	if err != nil {
		return nil, err
	}
	res := &pb.CaptureDeviceList{}
	for _, w := range wlist {
		cd := &pb.CaptureDevice{
			Type:   0,
			Device: w.DeviceName,
			Name:   w.Capabilities.Card,
		}
		res.Devices = append(res.Devices, cd)
	}

	for _, c := range cameras.Cameras {
		cd := &pb.CaptureDevice{
			Name:   c.Name,
			Device: c.URL,
			Type:   1,
		}
		res.Devices = append(res.Devices, cd)
	}
	sort.Slice(res.Devices, func(i, j int) bool {
		return res.Devices[i].Device < res.Devices[j].Device
	})
	return res, nil
}
