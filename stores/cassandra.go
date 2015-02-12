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
	Cluster             *cql.ClusterConfig
	Retentions          []CassandraRetention
	LocalDcName         string
	ReplicationStrategy string
	StrategyOptions     string
}

type CassandraRetention struct {
	Resolution int
	Period     int
}

func retentionToTimeInterval(retention string) int {
	res := strings.Split(retention, "")
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

	password_auth := new(cql.PasswordAuthenticator)
	password_auth.Username = config.StoreCassandra.Username
	password_auth.Password = config.StoreCassandra.Password

	store.Cluster = cql.NewCluster(config.StoreCassandra.Servers...)
	store.Cluster.Authenticator = password_auth
	store.Cluster.Keyspace = config.StoreCassandra.Keyspace
	store.StrategyOptions = config.StoreCassandra.StrategyOptions
	store.ReplicationStrategy = config.StoreCassandra.ReplicationStrategy

	retentions := make([]CassandraRetention, len(config.StoreCassandra.Retentions))
	for i, r := range config.StoreCassandra.Retentions {
		split_r := strings.Split(r, ":")
		resolution := retentionToTimeInterval(split_r[0])
		period := retentionToTimeInterval(split_r[1])
		retention := CassandraRetention{Resolution: resolution, Period: period}
		retentions[i] = retention
	}
	store.Retentions = retentions

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
