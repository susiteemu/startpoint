/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"goful/core/loader"
	requestSelectTui "goful/tui/request/select"

	"github.com/spf13/cobra"
)

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select and run a http request",
	Long:  `Launches a tui application where you can query a stored request and run it`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO handle err
		loadedRequests, _ := loader.ReadRequests("tmp")
		requestSelectTui.Start(loadedRequests)
	},
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
