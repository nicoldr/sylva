package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get FluxCD resources from the cluster",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
