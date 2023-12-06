package main

import (
	"testing"
	"time"
)

func TestDur(t *testing.T) {
	time_test(t, time.Duration(5)*time.Second)
	time_test(t, time.Duration(70)*time.Second)
	time_test(t, time.Duration(300)*time.Second)
}
func time_test(t *testing.T, tm time.Duration) {
	t.Logf("\"%v\" -> \"%s\"\n", tm, RenderDuration(tm))
}
