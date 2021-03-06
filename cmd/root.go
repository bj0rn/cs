package cmd

import (
	"fmt"
	"github.com/bj0rn/cs/pkg/switcher"
	"github.com/spf13/cobra"
	"os"
)

var (
	kubeconfig string
	aoconfig   string

	rootCmd = &cobra.Command{
		Use: "cs <env>",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}

			s := switcher.NewSwitcher(kubeconfig, aoconfig)

			cluster := args[0]

			namespace, err := cmd.Flags().GetString("namespace")
			if err != nil {
				fmt.Errorf("Missing parameter: namespace: %s", err)
				return
			}

			err = s.Switch(cluster, namespace)
			if err != nil {
				fmt.Println("Ups: %s", err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file (default is $HOME/.kube/config)")
	rootCmd.PersistentFlags().StringVar(&aoconfig, "aoconfig", "", "aoconfig file (default is $HOME/.ao.json)")
	rootCmd.Flags().StringP("namespace", "n", "aurora", "namespace")
}

func initConfig() {
	home, _ := os.UserHomeDir()
	kubeconfig = home + "/.kube/config"
	aoconfig = home + "/.ao.json"
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
