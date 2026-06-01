package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeconfig    string
	namespace     string
	output        string
	allNamespaces bool
)

var rootCmd = &cobra.Command{
	Use:   "sylvactl",
	Short: "A CLI tool to interact with FluxCD resources on a Kubernetes cluster",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file (default: ~/.kube/config)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace to query (default: all namespaces)")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "output format: table, json, yaml")
	rootCmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "query across all namespaces")
}
