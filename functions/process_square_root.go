package functions

import (
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/metrics"
	"math"
)

func init() {
	Functions["squareRoot"] = []string{"ProcessSquareRoot", "metric"}
}

func ProcessSquareRoot(dep_el chains.ChainEl) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	go func(our_el chains.ChainEl, dep_el chains.ChainEl) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		dep_el.Settings <- from
		dep_el.Settings <- until
		for {
			d := <-dep_el.Link
			if d.Known {
				our_el.Link <- *metrics.NewDatapoint(d.Ts, math.Sqrt(d.Value), true)
			} else {
				our_el.Link <- *metrics.NewDatapoint(d.Ts, 0.0, false)
			}
			if d.Ts >= until {
				return
			}
		}
	}(our_el, dep_el)
	return
}
