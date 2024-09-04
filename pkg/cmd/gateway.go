package cmd

import (
	"context"
	"github.com/liserc/open-socket/internal/gateway"
	"github.com/liserc/open-socket/pkg/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type GatewayCmd struct {
	*RootCmd
	ctx           context.Context
	configMap     map[string]any
	gatewayConfig *gateway.Config
}

func NewGatewayCmd() *GatewayCmd {
	var gatewayConfig gateway.Config
	ret := &GatewayCmd{gatewayConfig: &gatewayConfig}
	ret.configMap = map[string]any{
		ShareFileName:        &gatewayConfig.Share,
		DiscoveryCfgFilename: &gatewayConfig.Discovery,
		GatewayCfgFileName:   &gatewayConfig.Gateway,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (m *GatewayCmd) Exec() error {
	return m.Execute()
}

func (m *GatewayCmd) runE() error {
	return gateway.Start(m.ctx, m.Index(), m.gatewayConfig)
}
