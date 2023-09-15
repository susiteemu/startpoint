package cmd

import (
	"fmt"
	"goful-cli/client"
	"goful-cli/printer"

	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a http request",
	Long:  `Run a http request`,
	Run: func(cmd *cobra.Command, args []string) {
		var url = "https://httpbin.org/anything"
		var headers = map[string]string{"X-Foo": "bar", "X-Bar": "foo"}
		var body = []byte("{\"foo\":\"hello\"}")
		resp, err := client.DoRequest(url, "POST", headers, body)

		if err != nil {
			logrus.Errorf("Failed to perform request %v", err)
		}

		printed, _ := printer.PrintResponse(resp)
		fmt.Print(printed)

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
