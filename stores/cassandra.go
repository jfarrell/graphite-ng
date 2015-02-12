package stores

import (
	cql "github.com/gocql/gocql"
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/config"
	"github.com/graphite-ng/graphite-ng/metrics"
	"strconv"
	"strings"
)

type CassandraStore struct {
	cluster    *cql.ClusterConfig
	retentions []CassandraRetention
}

type CassandraRetention struct {
	resolution string
	period     string
}

// Dumb function to get a simple retention rate
// This really needs to come from the config
func retentions() {
	retention = new(CassandraRetention)
	retention.resolution = "1m"
	retention.period = "1w"
	return []CassandraRetention{retention}
}

func retentionToTimeIntervals(retention string) int {
	res := strings.Split(retention, ":")
	time, _ := strconv.Atoi(res[0])
	retention_time := retentionTypeToTime(time, res[1])

	return retention_time
}

func retentionTypeToTime(time int, period string) int {
	switch period {
	case "s":
		return time * 1
	case "m":
		return time * 60
	case "h":
		return time * 3600
	case "d":
		return time * 86400
	case "w":
		return time * 604800
	case "y":
		return time * 31536000
	}
	return 0
}

func NewCassandraStore(config config.Main) Store {
	store := new(CassandraStore)

	return store

}

func init() {
	InitFn["influxdb"] = NewCassandraStore
}

func (i CassandraStore) Add(metric metrics.Metric) (err error) {
	return

}

func (i CassandraStore) Get(name string) (our_el *chains.ChainEl, err error) {

	our_el = chains.NewChainEl()
	go func(our_el *chains.ChainEl) {
		from := <-our_el.Settings
		until := <-our_el.Settings

	}(our_el)
	return our_el, nil
}

func (i CassandraStore) Has(name string) (found bool, err error) {

	return
}
func (i CassandraStore) List() (list []string, err error) {

	return
}
