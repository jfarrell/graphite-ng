package functions

import (
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/metrics"
	"unicode"
)

func init() {
	Functions["summarize"] = []string{"ProcessSummarize", "metric"}
}

func ParseTimeOffset(offset string) (num_hours int) {
	var positive bool
	var leadingTerm = 0
	var unit string
	num_hours = 0
	//Check if the offset is signed, positive or not
	if unicode.IsDigit(rune(offset[0])) {
		positive = true
	} else {
		if rune(offset[0]) == '+' {
			positive = true
		} else {
			positive = false
		}
	}

	for i, c := range offset {
		if !unicode.IsDigit(c) {
			unit = offset[i:]
			break
		} else {
			leadingTerm = 10*leadingTerm + int(c-'0')
		}
	}
	//@todo: Fix so that we can work backwards in time
	if positive {
	}
	switch unit {
	case "years":
		num_hours += leadingTerm * 30 * 24
	case "months":
		num_hours += leadingTerm * 365 * 24
	case "week":
		num_hours += leadingTerm * 7 * 24
	case "days":
		num_hours += leadingTerm * 24
	case "hours":
		num_hours += leadingTerm
	}

	return
}

func BucketSum(bucket []float64) (result float64) {
	for _, point := range bucket {
		result += point
	}
	return
}

func BucketAverage(bucket []float64) (result float64) {
	result = BucketSum(bucket) / float64(len(bucket))
	return
}

func BucketMax(bucket []float64) (result float64) {
	for _, value := range bucket {
		if value > result {
			result = value
		}
	}
	return
}

func BucketMin(bucket []float64) (result float64) {
	for _, value := range bucket {
		if value < result {
			result = value
		}
	}
	return
}

func BucketLast(bucket []float64) (result float64) {
	result = bucket[len(bucket)-1]
	return
}

func ProcessSummarize(dep_el chains.ChainEl, interval string, operation string) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	go func(our_el chains.ChainEl, dep_el chains.ChainEl, operation string) {
		//Think bin width. Convert to seconds
		bucket_width := int32(ParseTimeOffset(interval) * 3600)
		buckets := make(map[int32][]float64)
		from := <-our_el.Settings
		until := <-our_el.Settings
		dep_el.Settings <- from
		dep_el.Settings <- until

		for {
			d := <-dep_el.Link
			if d.Known {
				metric_bin := d.Ts - (d.Ts % bucket_width)
				if buckets[metric_bin] == nil {
					buckets[metric_bin] = make([]float64, 0)
					buckets[metric_bin] = append(buckets[metric_bin], float64(d.Value))
				} else {
					buckets[metric_bin] = append(buckets[metric_bin], float64(d.Value))
				}
			}
			if d.Ts >= until {
				break
			}
		}

		bucket_counter := int32(0)
		for _, bucket := range buckets {
			bucket_timestamp := from + bucket_width*bucket_counter
			bucket_counter++

			//Reflect on this
			//@todo: Find the "go" way of speeding this up
			var result float64
			switch operation {
			case "avg":
				result = BucketAverage(bucket)
			case "sum":
				result = BucketSum(bucket)
			case "min":
				result = BucketMin(bucket)
			case "max":
				result = BucketMax(bucket)
			case "last":
				result = BucketLast(bucket)
			}
			our_el.Link <- *metrics.NewDatapoint(bucket_timestamp, result, true)
			if bucket_timestamp >= until {
				return
			}
		}
		return
	}(our_el, dep_el, operation)
	return
}
