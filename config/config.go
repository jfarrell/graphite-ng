package config

type Main struct {
	ListenAddr     string                 `toml:"listen_addr"`
	Stores         []string               `toml:"stores"`
	StoreCassandra storeCassandraInfo     `toml:"store_cassandra"`
	StoreES        storeElasticsearchInfo `toml:"store_elasticsearch"`
	StoreInflux    storeInfluxdbInfo      `toml:"store_influxdb"`
	StoreText      storeTextInfo          `toml:"store_text"`
}

type storeElasticsearchInfo struct {
	Host       string
	Port       int
	MaxPending int `toml:"max_pending"`
	CarbonPort int `toml:"carbon_port"`
}

type storeInfluxdbInfo struct {
	Host     string
	Username string
	Password string
	Database string
}

type storeTextInfo struct {
	Path string
}

type storeCassandraInfo struct {
	keyspace             string
	servers              []string
	username             string
	password             string
	replication_strategy string
	strategy_options     string
	local_dc_name        string
	retentions           []string
}
