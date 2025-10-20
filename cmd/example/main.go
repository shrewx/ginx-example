package main

import (
	"ginx-example/apis"
	"ginx-example/global"
	"github.com/shrewx/ginx"
	"github.com/spf13/cobra"
)

//go:generate toolx gen openapi

func main() {
	ginx.Launch(func(cmd *cobra.Command, args []string) {

		ginx.Parse(global.Config)

		global.Load()

		ginx.RunServer(&global.Config.Server, apis.V1Router)
	})
}
