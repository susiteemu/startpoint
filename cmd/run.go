/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"fmt"
	"goful-cli/client"
	"goful-cli/printer"

	"github.com/spf13/cobra"
)

type RunConfig struct {
	Plain        bool
	PrintHeaders bool
	PrintBody    bool
}

var runConfig RunConfig

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a http request",
	Long:  `Run a http request`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := client.DoRequest(
			"https://httpbin.org/anything",
			"POST",
			map[string]string{"X-Foo": "bar", "X-Bar": "foo"},
			[]byte("{\"foo\":\"Run run\"}"))

		var resp_str string
		var err error

		if runConfig.Plain {
			resp_str, err = printer.SprintResponse(resp, runConfig.PrintHeaders, runConfig.PrintBody)
		} else {
			resp_str, err = printer.SprintPrettyResponse(resp, runConfig.PrintHeaders, runConfig.PrintBody)
		}
		if err != nil {
			fmt.Errorf("Error %v", err)
		}
		fmt.Print(resp_str)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	const printHeadersP = "h"
	const printBodyP = "b"

	runCmd.PersistentFlags().BoolVarP(&runConfig.Plain, "plain", "p", false, "Print plain response without styling")
	runCmd.PersistentFlags().BoolVarP(&runConfig.PrintHeaders, "headers", printHeadersP, false, "Print response headers")
	runCmd.PersistentFlags().BoolVarP(&runConfig.PrintBody, "body", printBodyP, true, "Print response body")
	runCmd.PersistentFlags().StringSlice("print", []string{}, fmt.Sprintf("Print WHAT\n- '%s'\tPrint response headers\n- '%s'\tPrint response body", printHeadersP, printBodyP))

	runCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd == runCmd {
			printFlags, _ := cmd.Flags().GetStringSlice("print")
			for _, flag := range printFlags {
				if flag == printHeadersP {
					runConfig.PrintHeaders = true
				} else if flag == printBodyP {
					runConfig.PrintBody = true
				}
			}
		}
	}
}
