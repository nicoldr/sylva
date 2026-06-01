package cmd

import (
	"context"
	"fmt"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var helmReleasesCmd = &cobra.Command{
	Use:   "helmreleases",
	Short: "List HelmRelease objects from the cluster",
	RunE:  runGetHelmReleases,
}

func init() {
	getCmd.AddCommand(helmReleasesCmd)
}

func runGetHelmReleases(cmd *cobra.Command, args []string) error {
	c, err := buildClient(kubeconfig)
	if err != nil {
		return err
	}

	list := &helmv2.HelmReleaseList{}
	listOpts := []client.ListOption{}
	if !allNamespaces && namespace != "" {
		listOpts = append(listOpts, client.InNamespace(namespace))
	}

	if err := c.List(context.Background(), list, listOpts...); err != nil {
		return fmt.Errorf("failed to list HelmReleases: %w", err)
	}

	if len(list.Items) == 0 {
		fmt.Println("No HelmReleases found.")
		return nil
	}

	if output != "table" {
		return printOutput(list)
	}

	w := newTableWriter()
	fmt.Fprintln(w, "NAMESPACE\tNAME\tREADY\tSTATUS\tAGE")
	for _, hr := range list.Items {
		ready, message := readyCondition(hr.Status.Conditions)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			hr.Namespace, hr.Name, ready, message, sinceCreation(hr.CreationTimestamp.Time))
	}
	w.Flush()

	return nil
}
