package functions

import (
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/metrics"
	"math/rand"
	"time"
)

func init() {
	Functions["randomWalk"] = []string{"ProcessRandomWalk", "string"}
	Functions["randomWalkFunction"] = []string{"ProcessRandomWalk", "string"}
}

func RandomPoint()(float64){
	return rand.Float64()
}

func ProcessRandomWalk(metricName string) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	rand.Seed(time.Now().UTC().UnixNano())
	step_size := int32(60) // seconds
	step_counter := int32(0)
	go func(string) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		var timestamp int32
		for {
				timestamp = from + step_size*step_counter
				our_el.Link <- *metrics.NewDatapoint(timestamp, RandomPoint(), true)
				step_counter++
			if timestamp >= until {
				return
			}
		}
	}(metricName)
	return
}
