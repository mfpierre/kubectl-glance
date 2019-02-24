package main

import (
	"os"

	"github.com/mfpierre/kubectl-glance/pkg/cmd"
	"github.com/spf13/pflag"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-glance", pflag.ExitOnError)
	pflag.CommandLine = flags

	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
