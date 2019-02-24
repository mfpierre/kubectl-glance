package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

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
			err := GlobalSettings.GetNamespaces()
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
