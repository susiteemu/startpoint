/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"goful/client"
	"goful/client/validator"
	"goful/printer"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type RunConfig struct {
	Plain        bool
	PrintHeaders bool
	PrintBody    bool
}

type RunArgs struct {
	Method string
	Url    string
}

type RunFlags struct {
	Body    string
	Headers []string
}

var runConfig RunConfig
var runFlags RunFlags

var runCmd = &cobra.Command{
	Use:   "run [METHOD] [URL]",
	Short: "Run a http request",
	Long:  `Run a http request`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.RangeArgs(0, 2)(cmd, args); err != nil {
			return err
		}

		if len(args) == 0 {
			return nil
		}

		parsedArgs := ParseArgs(args)

		if !validator.IsValidMethod(parsedArgs.Method) {
			return errors.New(fmt.Sprintf("METHOD must be one of following: %v", strings.Join(validator.ValidMethods, ", ")))
		}
		if !validator.IsValidUrl(parsedArgs.Url) {
			return errors.New(fmt.Sprintf("URL is not valid"))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var resp *http.Response
		runArgs := ParseArgs(args)
		if runArgs != (RunArgs{}) {
			// TODO err
			headers := toHeadersMap(runFlags.Headers)
			resp, _ = client.DoRequest(
				runArgs.Url,
				runArgs.Method,
				headers,
				[]byte(runFlags.Body))
		} else {
			// TODO get from --name
			resp, _ = client.DoRequest(
				"https://httpbin.org/anything",
				"POST",
				map[string]string{"X-Foo": "bar", "X-Bar": "foo"},
				[]byte("{\"foo\":\"Run run\"}"))
		}

		var respStr string
		var err error

		if runConfig.Plain {
			respStr, err = printer.SprintResponse(resp, runConfig.PrintHeaders, runConfig.PrintBody)
		} else {
			respStr, err = printer.SprintPrettyResponse(resp, runConfig.PrintHeaders, runConfig.PrintBody)
		}
		if err != nil {
			fmt.Print(fmt.Errorf("error %v", err))
		}
		fmt.Print(respStr)
	},
}

func ParseArgs(args []string) RunArgs {
	if len(args) == 0 {
		return RunArgs{}
	}
	if len(args) == 1 {
		return RunArgs{validator.DefaultBodilessMethod, args[0]}
	}
	return RunArgs{args[0], args[1]}
}

func toHeadersMap(headers []string) map[string]string {
	var headerMap = map[string]string{}
	for _, h := range headers {
		headerParts := strings.Split(h, ":")
		if len(headerParts) == 2 {
			headerMap[headerParts[0]] = headerParts[1]
		}
	}
	return headerMap
}

func init() {
	rootCmd.AddCommand(runCmd)

	const printHeadersP = "h"
	const printBodyP = "b"

	runConfig.PrintBody = true

	runCmd.PersistentFlags().BoolVarP(&runConfig.Plain, "plain", "p", false, "Print plain response without styling")
	runCmd.PersistentFlags().Bool("no-body", false, "Print no body")
	runCmd.PersistentFlags().StringSlice("print", []string{}, fmt.Sprintf("Print WHAT\n- '%s'\tPrint response headers\n- '%s'\tPrint response body", printHeadersP, printBodyP))
	runCmd.Flags().StringVarP(&runFlags.Body, "body", "b", "", "Request body")
	runCmd.Flags().StringSliceVarP(&runFlags.Headers, "header", "h", []string{}, "Request headers formatted as HeaderName:HeaderValue")

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
			noBody, _ := cmd.PersistentFlags().GetBool("no-body")
			runConfig.PrintBody = !noBody
		}
	}
}
