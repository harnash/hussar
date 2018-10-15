package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/harnash/hussar/cavalry/transport"
	"github.com/nats-io/nats"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

var nc *nats.Conn

func discoverServices(logger *zap.SugaredLogger) {
	for {
		msg, err := nc.Request("Discovery.hussary", nil, 1000*time.Millisecond)
		if err != nil {
			logger.With("err", err).Error("something went wrong. Waiting 2 seconds before retrying:")
			continue
		}
		fileServerAddressTransport := Transport.DiscoverableServiceTransport{}
		err = proto.Unmarshal(msg.Data, &fileServerAddressTransport)
		if err != nil {
			logger.With("err", err).Error("something went wrong. Waiting 2 seconds before retrying:")
			continue
		}

		logger.With("address", fileServerAddressTransport.Address).Info("detected new service")
	}
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	sugar := logger.With(zap.String("app", "hussar-hq")).Sugar()

	if len(os.Args) != 2 {
		sugar.Error("wrong number of arguments. Need NATS server address.")
		os.Exit(1)
	}

	nc, err = nats.Connect(os.Args[1])
	if err != nil {
		sugar.With("err", err, "address", os.Args[1]).Error("error connecting to NATS server")
		os.Exit(2)
	}

	sugar.With("nats", os.Args[1]).Info("starting HQ!")

	go discoverServices(sugar)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
