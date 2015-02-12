package stores

import (
	"fmt"
	cql "github.com/gocql/gocql"
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/config"
	"github.com/graphite-ng/graphite-ng/metrics"
	"strconv"
	"strings"
)

type CassandraStore struct {
	Cluster             *cql.ClusterConfig
	Keyspace            string
	Retentions          []CassandraRetention
	LocalDcName         string
	ReplicationStrategy string
	StrategyOptions     map[string]int
	Session             *cql.Session
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
	store.Cluster = cql.NewCluster(config.StoreCassandra.Servers...)

	if config.StoreCassandra.Username != "" && config.StoreCassandra.Password != "" {
		password_auth := new(cql.PasswordAuthenticator)
		password_auth.Username = config.StoreCassandra.Username
		password_auth.Password = config.StoreCassandra.Password
		store.Cluster.Authenticator = password_auth
	}

	store.Keyspace = config.StoreCassandra.Keyspace
	store.StrategyOptions = config.StoreCassandra.StrategyOptions.Options
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

	initializeCassandraKeyspace(store)
	store.Cluster.Keyspace = store.Keyspace

	store.Session, _ = store.Cluster.CreateSession()
	createCassandraTables(store)

	return store
}

func strategyOptionsToString(store *CassandraStore) string {
	mapping := make([]string, len(store.StrategyOptions))
	counter := 0
	for key, value := range store.StrategyOptions {
		mapping[counter] = fmt.Sprintf("'%s' : %s", key, value)
		counter += 1
	}
	return strings.Join(mapping, ",")
}

func createCassandraTables(store *CassandraStore) {
	tables := []string{"global_nodes", "metadata"}
	if store.LocalDcName != "" {
		tables = append(tables, store.LocalDcName)
	}

	for _, table_name := range tables {
		query := fmt.Sprintf("CREATE TABLE %s (key text, column1 text, value text, PRIMARY KEY(key, column1));", table_name)
		if err := store.Session.Query(query); err != nil {
			panic(err)
		}
	}
}

func initializeCassandraKeyspace(store *CassandraStore) {
	session, _ := store.Cluster.CreateSession()
	replication_options := strategyOptionsToString(store)
	query := fmt.Sprintf("CREATE KESYPACE IF NOT EXISTS %s WITH REPLICATION {'class' %s : %s};", store.Keyspace, store.ReplicationStrategy, replication_options)

	if err := session.Query(query); err != nil {
		panic(err)
	}
	session.Close()
}

func init() {
	InitFn["cassandra"] = NewCassandraStore
}

func (i CassandraStore) Add(metric metrics.Metric) (err error) {

	return
}

func (i CassandraStore) Get(name string) (our_el *chains.ChainEl, err error) {

	return
}

func (i CassandraStore) Has(name string) (found bool, err error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM dc_%s_nodes WHERE key='%s' LIMIT 1", name, i.LocalDcName)
	count := 0
	if err := i.Session.Query(query).Scan(&count); err != nil {
		panic(err)
	} else if count > 0 {
		found = true
	}
	return
}
func (i CassandraStore) List() (list []string, err error) {

	return
}
