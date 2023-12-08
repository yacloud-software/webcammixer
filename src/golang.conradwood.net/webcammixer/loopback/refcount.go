package loopback

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"strconv"
	"strings"
	"time"
)

const (
	LOOPMOD_REFCOUNT_FILENAME = "/sys/module/v4l2loopback/refcnt"
)

var (
	refcounter_started = false
	last_change        time.Time
	loopback_refcnt    = 0
)

func refcounter() {
	refcounter_started = true
	t := time.Duration(1) * time.Second
	last_change = time.Now()

	for {
		time.Sleep(t)
		t = time.Duration(1) * time.Second
		err := readrefcount()
		if err != nil {
			fmt.Printf("Failed to read refcount: %s\n", err)
			t = time.Duration(1) * time.Second
			continue
		}
	}
}
func readrefcount() error {
	b, err := utils.ReadFile(LOOPMOD_REFCOUNT_FILENAME)
	if err != nil {
		return err
	}
	s := string(b)
	s = strings.Trim(s, "\n")
	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	if num != loopback_refcnt {
		last_change = time.Now()
		fmt.Printf("Loopback: ref-count changed from %d to %d\n", loopback_refcnt, num)
	}
	loopback_refcnt = num
	return nil
}

type LoopbackStatus struct {
	RefCount   int
	LastChange time.Time
}

func IsWatched() bool {
	return loopback_refcnt > 1
}
func Status() *LoopbackStatus {
	return &LoopbackStatus{
		LastChange: last_change,
		RefCount:   loopback_refcnt,
	}
}
