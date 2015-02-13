package stores

import (
	"fmt"
	cql "github.com/gocql/gocql"
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/config"
	"github.com/graphite-ng/graphite-ng/metrics"
	"log"
	"strconv"
	"strings"
)

type CassandraRetention struct {
	Resolution int
	Period     int
}

type CassandraStore struct {
	Cluster             *cql.ClusterConfig
	Keyspace            string
	Retentions          []CassandraRetention
	LocalDcName         string
	ReplicationStrategy string
	StrategyOptions     map[string]config.CassandraReplication
	Session             *cql.Session
}

type CassandraNode struct {
	Name          string
	ChildOrMetric string
	IsMetric      bool
}

type CassandraMetric struct {
	timestamp int
	value     float64
}

type CassandraMetadata struct {
	Key               string
	AggregationMethod string
	Retentions        []int
	StartTime         int
	TimeStep          int
	// xFilesFactor renamed to make sense
	//Retentions        []CassandraRetention
	NonNilPercentage float64
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
	store.StrategyOptions = config.StoreCassandra.StrategyOptions
	store.ReplicationStrategy = config.StoreCassandra.ReplicationStrategy

	store.LocalDcName = strings.Replace(config.StoreCassandra.LocalDcName, "-", "_", -1)

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
		val := strconv.Itoa(value.Replication)
		mapping[counter] = fmt.Sprintf("'%s' : %s", key, val)
		counter += 1
	}
	return strings.Join(mapping, ",")
}

func createCassandraTables(store *CassandraStore) {
	tables := []string{"global_nodes", "metadata"}
	if store.LocalDcName != "" {
		tables = append(tables, fmt.Sprintf("dc_%s_nodes", store.LocalDcName))
	}

	for _, table_name := range tables {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (key text, column1 text, value text, PRIMARY KEY(key, column1)) WITH COMPACT STORAGE", table_name)
		if err := store.Session.Query(query).Exec(); err != nil {
			panic(err)
		}
	}

	for _, retention := range store.Retentions {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS ts%d (key text, column1 bigint, value float, PRIMARY KEY(key, column1)) WITH COMPACT STORAGE", retention.Resolution)
		if err := store.Session.Query(query).Exec(); err != nil {
			panic(err)
		}
	}
}

func initializeCassandraKeyspace(store *CassandraStore) {
	session, _ := store.Cluster.CreateSession()
	replication_options := strategyOptionsToString(store)
	query := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class' : '%s', %s}", store.Keyspace, store.ReplicationStrategy, replication_options)

	if err := session.Query(query).Exec(); err != nil {
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

func getMetadata(name string, i CassandraStore) CassandraMetadata {
	var key string
	var column_key string
	var value string
	metadata_map := map[string]interface{}{
		"key":        &key,
		"column_key": &column_key,
		"value":      &value,
	}

	metadata := new(CassandraMetadata)
	metadata.Key = name

	query := fmt.Sprintf("SELECT * FROM metadata WHERE key = %s", name)
	iter := i.Session.Query(query).Iter()
	for iter.MapScan(metadata_map) {
		if str, ok := metadata_map["value"].(string); ok {
			switch metadata_map["column_key"] {
			case "aggregationMethod":
				metadata.AggregationMethod = str
			case "retentions":
				val := strings.Split(str, ",")
				vals := make([]int, len(val))
				for i, v := range val {
					if int_v, err := strconv.Atoi(v); err != nil {
						vals[i] = int_v
					} else {
						log.Fatal("Unable to convert a retention value to an integer")

					}
				}
				metadata.Retentions = vals
			case "startTime":
				if num_value, err := strconv.Atoi(str); err != nil {
					log.Fatal("Unable to convert startTime to integer")
				} else {
					metadata.StartTime = num_value
				}
			case "timeStep":
				if num_value, err := strconv.Atoi(str); err != nil {
					log.Fatal("Unable to convert timeStep to integer")
				} else {
					metadata.TimeStep = num_value
				}
			case "xFilesFactor":
				if num_value, err := strconv.ParseFloat(str, 64); err != nil {
					log.Fatal("Unable to convert xFilesFactor to float")
				} else {
					metadata.NonNilPercentage = num_value
				}
			}
		} else {
			log.Fatal("Got a non string type back from metadata query")
		}
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return *metadata
}

func getData(name string, i CassandraStore, metadata CassandraMetadata, from int32, until int32) int {

	var key string
	var timestamp int
	var value float64
	value_map := map[string]interface{}{
		"key":       &key,
		"timestamp": &timestamp,
		"value":     &value,
	}

	// For testing only. I'm too lazy to implement real querying right now
	// TODO Actually implement searching in the other ts* tables for the appropriate historical data
	query := fmt.Sprintf("SELECT * FROM ts60 WHERE key = %s AND column1 >= %d AND column1 <= %d ORDER BY column1", name, from, until)
	iter := i.Session.Query(query).Iter()
	for iter.MapScan(value_map) {
	}
	return 1
}

func (i CassandraStore) Get(name string) (our_el *chains.ChainEl, err error) {
	our_el = chains.NewChainEl()
	go func(our_el *chains.ChainEl) {
		from := <-our_el.Settings
		until := <-our_el.Settings

		metadata := getMetadata(name, i)
		data := getData(name, i, metadata, from, until)

	}(our_el)

	return our_el, nil
}

func (i CassandraStore) Has(name string) (found bool, err error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM dc_%s_nodes WHERE key='%s' LIMIT 1", i.LocalDcName, name)
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
