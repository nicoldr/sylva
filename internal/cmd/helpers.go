package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// Return a human-readable duration since the given time
func sinceCreation(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

// Extract the Ready condition status and message from a condition list
func readyCondition(conditions []metav1.Condition) (ready, message string) {
	for _, cond := range conditions {
		if cond.Type == "Ready" {
			return string(cond.Status), cond.Message
		}
	}
	return "Unknown", ""
}

// Return a tabwriter configured for consistent CLI table output
func newTableWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
}

// printOutput prints any object as json or yaml based on the output flag
func printOutput(obj interface{}) error {
	switch output {
	case "json":
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(obj)
		if err != nil {
			return fmt.Errorf("failed to marshal to YAML: %w", err)
		}
		fmt.Print(string(data))
	default:
		return fmt.Errorf("unknown output format: %s (supported: table, json, yaml)", output)
	}
	return nil
}
