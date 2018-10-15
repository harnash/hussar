package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/harnash/hussar/cavalry/transport"
	"github.com/nats-io/nats"
	"go.uber.org/zap"
	"os"
	"sync"
)

func RunServiceDiscoverable(address string) {
	nc, err := nats.Connect(address)
	if err != nil {
		fmt.Println("Can't connect to NATS. Service is not discoverable.")
	}
	nc.Subscribe("Discovery.hussary", func(m *nats.Msg) {
		serviceAddressTransport := transport.DiscoverableServiceTransport{Address: "http://localhost:3000"}
		data, err := proto.Marshal(&serviceAddressTransport)
		if err == nil {
			nc.Publish(m.Reply, data)
		}
	})
}

func main() {
	var err error
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	sugar := logger.With(zap.String("app", "hussar-cavalry")).Sugar()

	sugar.With("nats", os.Args[1]).Info("preparing the cavalry!")

	RunServiceDiscoverable(os.Args[1])

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}