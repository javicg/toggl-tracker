package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print current project version",
		Long:  "Print current project version",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "toggl-sync v0.0.1")
			return
		},
	}
}
