package gateway

import (
	"context"
	"github.com/liserc/open-socket/pkg/config"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

type Config struct {
	Gateway   config.Gateway
	Discovery config.Discovery
}

// Start run ws server.
func Start(ctx context.Context, index int, conf *Config) error {
	log.CInfo(ctx, "gateway initializing", "rpcPorts", conf.Gateway.RPC.Ports, "prometheusPorts",
		conf.Gateway.Prometheus.Ports, "wsPort", conf.Gateway.WS.Ports)
	wsPort, err := datautil.GetElemByIndex(conf.Gateway.WS.Ports, index)
	if err != nil {
		return err
	}
	rpcPort, err := datautil.GetElemByIndex(conf.Gateway.RPC.Ports, index)
	if err != nil {
		return err
	}
	log.CInfo(ctx, "", wsPort, rpcPort)
	return nil
}
