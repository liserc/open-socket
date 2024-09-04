package gateway

import (
	"context"
	"github.com/liserc/open-socket/pkg/rpc"
	"github.com/liserc/open-socket/protocol/gateway"
	"github.com/openimsdk/tools/discovery"
	"google.golang.org/grpc"
)

type RpcServer struct {
	gateway.UnimplementedGatewayServer
	rpcPort      int
	socketServer *SocketServer
}

func NewRpcServer(rpcPort int, socketServer *SocketServer) *RpcServer {
	return &RpcServer{
		rpcPort:      rpcPort,
		socketServer: socketServer,
	}
}

func (s *RpcServer) InitServer(ctx context.Context, config *Config, registry discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	s.socketServer.SetDiscoveryRegistry(registry)
	gateway.RegisterGatewayServer(server, s)
	return nil
}

func (s *RpcServer) Start(ctx context.Context, index int, conf *Config) error {
	return rpc.Start(ctx, &conf.Discovery, &conf.Gateway.Prometheus, conf.Gateway.ListenIP,
		conf.Gateway.RPC.RegisterIP, conf.Gateway.RPC.Ports, index,
		conf.Share.RpcRegisterName.Gateway, &conf.Share, conf,
		s.InitServer)
}

func (s *RpcServer) OnlineBatchPushOneMsg(ctx context.Context, req *gateway.OnlineBatchPushOneMsgReq) (*gateway.OnlineBatchPushOneMsgResp, error) {
	// todo implement
	return nil, nil
}
