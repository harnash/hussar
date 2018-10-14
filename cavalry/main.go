package main

import (
	"go.uber.org/zap"
	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
	"os"
)

func main() {
	var err error
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	sugar := logger.With(zap.String("app", "hussar-cavalry")).Sugar()

	sugar.Info("preparing the cavalry!")

	var sock mangos.Socket
	var msg []byte

	if sock, err = sub.NewSocket(); err != nil {
		sugar.Errorw("can't get new sub socket", "err", err)
		os.Exit(1)
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial("tcp://127.0.0.1:40899"); err != nil {
		sugar.Errorw("can't dial on sub socket", "err", err)
		os.Exit(2)
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		sugar.Errorw("cannot subscribe", "err", err)
		os.Exit(3)
	}
	for {
		if msg, err = sock.Recv(); err != nil {
			sugar.Errorw("cannot recv", "err", err)
			os.Exit(4)
		}
		sugar.Infow("received message", "msg", msg)
	}
}