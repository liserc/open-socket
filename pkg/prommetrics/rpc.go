package prommetrics

import (
	gp "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/liserc/open-socket/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

const rpcPath = commonPath

var (
	grpcMetrics *gp.ServerMetrics
	rpcCounter  = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_count",
			Help: "Total number of RPC calls",
		},
		[]string{"name", "path", "code"},
	)
)

func RpcInit(cs []prometheus.Collector, prometheusPort int) error {
	reg := prometheus.NewRegistry()
	cs = append(append(
		baseCollector,
		rpcCounter,
	), cs...)
	return Init(reg, prometheusPort, rpcPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}), cs...)
}

func RPCCall(name string, path string, code int) {
	rpcCounter.With(prometheus.Labels{"name": name, "path": path, "code": strconv.Itoa(code)}).Inc()
}

func GetGrpcServerMetrics() *gp.ServerMetrics {
	if grpcMetrics == nil {
		grpcMetrics = gp.NewServerMetrics()
		grpcMetrics.EnableHandlingTimeHistogram()
	}
	return grpcMetrics
}

func GetGrpcCusMetrics(registerName string, share *config.Share) []prometheus.Collector {
	switch registerName {
	case share.RpcRegisterName.Gateway:
		return []prometheus.Collector{OnlineUserGauge}
	case share.RpcRegisterName.Push:
		return []prometheus.Collector{MsgOfflinePushFailedCounter}
	default:
		return nil
	}
}
