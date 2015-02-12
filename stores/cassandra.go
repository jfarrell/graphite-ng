package stores

import (
	"fmt"
	cql "github.com/gocql/gocql"
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/config"
	"github.com/graphite-ng/graphite-ng/metrics"
	"github.com/graphite-ng/graphite-ng/util"
	"strconv"
	"string"
)

type CassandraStore struct {
	cluster    *gocql.ClusterConfig
	retentions []CassandraRetention
}

type CassandraRetention struct {
	resolution string
	period     string
}

func retentionToTimeIntervals(retention CassandraRetention) {
	res := string.Split(retention.resolution, ":")
	res, er = strconv.Atoi(res)
	p := string.Split(retention.period, ":")

	return retentionTypeToTime(res, p)
}

func retentionTypeToTime(int time, string period) {
	switch period {
	case 's':
		return time * 1
	case 'm':
		return time * 60
	case 'h':
		return time * 3600
	case 'd':
		return time * 86400
	case 'w':
		return time * 604800
	case 'y':
		return time * 31536000
	}
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
