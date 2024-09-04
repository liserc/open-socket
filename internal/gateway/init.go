package gateway

import (
	"context"
	"github.com/liserc/open-socket/pkg/config"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

type Config struct {
	Share     config.Share
	Discovery config.Discovery
	Gateway   config.Gateway
}

func Start(ctx context.Context, index int, conf *Config) error {
	log.ZInfo(ctx, "gateway initializing", "rpcPorts", conf.Gateway.RPC.Ports, "prometheusPorts", conf.Gateway.Prometheus.Ports, "wsPort", conf.Gateway.WS.Ports)
	wsPort, err := datautil.GetElemByIndex(conf.Gateway.WS.Ports, index)
	if err != nil {
		return err
	}
	rpcPort, err := datautil.GetElemByIndex(conf.Gateway.RPC.Ports, index)
	if err != nil {
		return err
	}
	socketServer := NewSocketServer(
		conf,
		WithPort(wsPort),
		WithMaxConnNum(int64(conf.Gateway.WS.WebsocketMaxConnNum)),
		WithMessageMaxMsgLength(conf.Gateway.WS.WebsocketMaxMsgLen),
	)
	rpcServer := NewRpcServer(rpcPort, socketServer)
	done := make(chan error)
	go func() {
		err = rpcServer.Start(ctx, index, conf)
		done <- err
	}()
	return socketServer.Run(done)
}
