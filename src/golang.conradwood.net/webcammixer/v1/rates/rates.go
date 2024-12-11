package rates

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"golang.conradwood.net/go-easyops/utils"
)

var (
	lock             sync.Mutex
	last_printed     time.Time
	rate_calculators = make(map[string]*name_map)
)

type name_map struct {
	rc utils.RateCalculator
	sa *utils.SlidingAverage
}

func Inc(name string) {
	lock.Lock()
	nm, found := rate_calculators[name]
	if !found {
		nm = &name_map{
			rc: utils.NewRateCalculator(name),
			sa: utils.NewSlidingAverage(),
		}
		rate_calculators[name] = nm
	}
	lock.Unlock()
	nm.rc.Add(1)
	nm.sa.Add(1, 1)
	Print()
}
func Print() {
	if time.Since(last_printed) < time.Duration(3)*time.Second {
		return
	}
	last_printed = time.Now()

	var x []string
	lock.Lock()
	rc_copy := make(map[string]*name_map)
	for k, v := range rate_calculators {
		x = append(x, k)
		rc_copy[k] = v
	}
	lock.Unlock()
	sort.Slice(x, func(i, j int) bool {
		return x[i] < x[j]
	})
	s := ""
	deli := ""
	for _, n := range x {
		rc := rc_copy[n]
		f1 := rc.rc.Rate()
		s = s + deli + fmt.Sprintf("Rate \"%s\"=%0.1f", n, f1)
		deli = ", "
		rc.rc.Reset()
	}
	fmt.Println(s)
}
