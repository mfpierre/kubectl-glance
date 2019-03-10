package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"

	// Required auth libraries
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var (
	// RootCmd is the root command of the plugin
	RootCmd = &cobra.Command{
		Use:   "kubectl glance",
		Short: "Kubectl glance",
		RunE: func(cmd *cobra.Command, args []string) error {
			GlobalSettings.InitClient()
			ns, err := GlobalSettings.GetNamespaces()
			fmt.Printf(
				"%s Namespaces\n",
				color.GreenString(strconv.Itoa(ns)),
			)
			nodesOK, nodesKO, err := GlobalSettings.GetNodes()
			fmt.Printf(
				"%s Nodes (%s KO)\n",
				color.GreenString(strconv.Itoa(nodesOK)),
				color.RedString(strconv.Itoa(nodesKO)),
			)
			return err
		},
	}
	// GlobalSettings holds common options/utils
	GlobalSettings *globalSettings
)

type globalSettings struct {
	configFlags *genericclioptions.ConfigFlags
	client      *kubernetes.Clientset
	namespace   string
	restConfig  *rest.Config
}

func init() {
	flags := genericclioptions.NewConfigFlags()
	flags.AddFlags(RootCmd.PersistentFlags())
	GlobalSettings = &globalSettings{
		configFlags: flags,
	}
}
