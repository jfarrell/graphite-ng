# Thrift

Our thrift integration is based on three parts. The server and two example clients. All of the example code below starts from the graphite-ng directory and assumes you have go and ruby installed on your system.

## Prerequesites

	go get git.apache.org/thrift.git/lib/go/thrift
	gem install thrift

## Thrift Server

There are two examples. One is a very simple test of the thrift data and connection. The second is embedded into the grahpite-ng application.

	# simple test server
	go run thrift_server.go
	# graphite-ng application with a thrift server
	gom run graphite-ng.go

## Golang Client

	cd thrift
	go run go_client.go -P=json metrics
	go run go_client.go -P=json render

## Ruby Client

	gem install thrift
	cd thrift
	ruby rb_client.rb metrics
	ruby rb_client.rb render

## Generating Client Libraries

The client libraries are generated from the .thrift file for the specified language. These will appear in a gen-\* file (e.g. gen-go) and need to be incorprated into your application. For the examples above, the clients are a simple CLI utility that retrieves the data.

	thrift --gen go -o thrift graphiteng.thrift
	thrift --gen rb -o thrift graphiteng.thrift
