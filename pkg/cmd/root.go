package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
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
			namespacedResources, err := GlobalSettings.GetNamespacedRessources()
			w := tabwriter.NewWriter(os.Stdout, 3, 4, 1, ' ', tabwriter.AlignRight)
			for resource, value := range namespacedResources {
				fmt.Fprintln(
					w,
					fmt.Sprintf("%s\t %s",
						strings.Title(resource),
						color.GreenString(strconv.Itoa(value)),
					),
				)
			}

			pv, err := GlobalSettings.GetPersistentVolumes()
			fmt.Fprintln(
				w,
				fmt.Sprintf("Persistent volumes\t %s",
					color.GreenString(strconv.Itoa(pv)),
				),
			)
			nodesOK, nodesKO, err := GlobalSettings.GetNodes()
			fmt.Fprintln(
				w,
				fmt.Sprintf("Nodes\t %s (%s Unschedulable)",
					color.GreenString(strconv.Itoa(nodesOK)),
					color.RedString(strconv.Itoa(nodesKO)),
				),
			)
			w.Flush()
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
