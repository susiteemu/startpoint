package cmd

import (
	"fmt"
	"goful-cli/client"
	"goful-cli/printer"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
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
		if plain, _ := cmd.Flags().GetBool("plain"); plain {
			resp_str, err = printer.SprintResponse(resp)
		} else {
			resp_str, err = printer.SprintPrettyResponse(resp)
		}
		if err != nil {
			fmt.Errorf("Error %v", err)
		}
		fmt.Print(resp_str)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().Bool("plain", false, "Print plain response without styling")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
