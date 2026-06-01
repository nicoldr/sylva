package cmd

import (
	"context"
	"fmt"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var kustomizationsCmd = &cobra.Command{
	Use:   "kustomizations",
	Short: "List Kustomization objects from the cluster",
	RunE:  runGetKustomizations,
}

func init() {
	getCmd.AddCommand(kustomizationsCmd)
}

func runGetKustomizations(cmd *cobra.Command, args []string) error {
	c, err := buildClient(kubeconfig)
	if err != nil {
		return err
	}

	list := &kustomizev1.KustomizationList{}
	listOpts := []client.ListOption{}
	if !allNamespaces && namespace != "" {
		listOpts = append(listOpts, client.InNamespace(namespace))
	}

	if err := c.List(context.Background(), list, listOpts...); err != nil {
		return fmt.Errorf("failed to list Kustomizations: %w", err)
	}

	if len(list.Items) == 0 {
		fmt.Println("No Kustomizations found.")
		return nil
	}

	if output != "table" {
		return printOutput(list)
	}

	w := newTableWriter()
	fmt.Fprintln(w, "NAMESPACE\tNAME\tREADY\tSTATUS\tAGE")
	for _, ks := range list.Items {
		ready, message := readyCondition(ks.Status.Conditions)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			ks.Namespace, ks.Name, ready, message, sinceCreation(ks.CreationTimestamp.Time))
	}
	w.Flush()

	return nil
}
