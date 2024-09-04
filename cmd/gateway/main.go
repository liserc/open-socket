package main

import (
	"github.com/liserc/open-socket/pkg/cmd"
	"github.com/openimsdk/tools/system/program"
)

func main() {
	if err := cmd.NewGatewayCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
