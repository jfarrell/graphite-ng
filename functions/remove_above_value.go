package functions

import (
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/metrics"
)

func init() {
	Functions["removeAboveValue"] = []string{"RemoveAboveValue", "metric", "float"}
}

func RemoveAboveValue(dep_el chains.ChainEl, threshold float64) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	go func(our_el chains.ChainEl, dep_el chains.ChainEl, threshold float64) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		dep_el.Settings <- from
		dep_el.Settings <- until
		for {
			d := <-dep_el.Link
			if d.Known && (d.Value < threshold){
				our_el.Link <- *metrics.NewDatapoint(d.Ts, d.Value, true)
			} else {
				our_el.Link <- *metrics.NewDatapoint(d.Ts, 0.0, false)
			}
			if d.Ts >= until {
				return
			}
		}
	}(our_el, dep_el, threshold)
	return
}
