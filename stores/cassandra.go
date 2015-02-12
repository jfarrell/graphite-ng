package stores

import (
	"fmt"
	cql "github.com/gocql/gocql"
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/config"
	"github.com/graphite-ng/graphite-ng/metrics"
	"github.com/graphite-ng/graphite-ng/util"
)

type CassandraStore struct {
}

func NewCassandraStore(config config.Main) Store {

}

func init() {
	InitFn["influxdb"] = NewCassandraStore
}

func (i CassandraStore) Add(metric metrics.Metric) (err error) {

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
