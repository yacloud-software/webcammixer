package rates

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"sort"
	"sync"
	"time"
)

var (
	lock             sync.Mutex
	last_printed     time.Time
	rate_calculators = make(map[string]utils.RateCalculator)
)

func Inc(name string) {
	lock.Lock()
	rc, found := rate_calculators[name]
	if !found {
		rc = utils.NewRateCalculator(name)
		rate_calculators[name] = rc
	}
	lock.Unlock()
	rc.Add(1)
	Print()
}
func Print() {
	if time.Since(last_printed) < time.Duration(3)*time.Second {
		return
	}
	last_printed = time.Now()

	var x []string
	lock.Lock()
	rc_copy := make(map[string]utils.RateCalculator)
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
		s = s + deli + fmt.Sprintf("Rate \"%s\"=%0.1f", n, rc.Rate())
		deli = ", "
	}
	fmt.Println(s)
}
