package main

import (
	"go.uber.org/zap"
	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/pub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
	"os"
	"time"
)

func date() string {
	return time.Now().Format(time.ANSIC)
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	sugar := logger.With(zap.String("app", "hussar-hq")).Sugar()

	sugar.Info("starting HQ!")

	var sock mangos.Socket
	if sock, err = pub.NewSocket(); err != nil {
		sugar.Errorw("cannot create socket", "err", err)
		os.Exit(1)
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())

	if err = sock.Listen("tcp://127.0.0.1:40899"); err != nil {
		sugar.Errorw("cannot listen on a given url", "err", err)
		os.Exit(2)
	}

	for {
		d := date()
		sugar.Infow("publishing date", zap.String("date", d))
		if err = sock.Send([]byte(d)); err != nil {
			sugar.Errorw("failed publishing", "err", err)
			os.Exit(3)
		}
		time.Sleep(time.Second)
	}
}
