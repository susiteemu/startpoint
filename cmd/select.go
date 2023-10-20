/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	selectui "goful/tui/select"

	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select and run a http request",
	Long:  `Launches a tui application where you can query a stored request and run it`,
	Run: func(cmd *cobra.Command, args []string) {
		selectui.Start()
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
