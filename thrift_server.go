package main

import (
	"crypto/tls"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/graphite-ng/graphite-ng/thrift/gen-go"
)

type GraphiteNGHandler struct{}

func (g *GraphiteNGHandler) Render() (r *graphiteng.RenderData, err error) {
	r = graphiteng.NewRenderData()
	r.Target = "my.test.target"

	for i := 0; i < 10; i++ {
		d := graphiteng.NewDatapoint()
		d.Value = float64(1.234) * float64(i)
		d.Timestamp = int32(1423683330) + int32(i)
		r.Datapoints = append(r.Datapoints, d)
	}

	return r, nil
}

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TServerTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}

	if err != nil {
		return err
	}
	fmt.Printf("%T\n", transport)
	handler := new(GraphiteNGHandler)
	processor := graphiteng.NewGraphiteNGProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	return server.Serve()
}

func main() {
	runServer(thrift.NewTTransportFactory(), thrift.NewTJSONProtocolFactory(), "127.0.0.1:9090", false)
}
