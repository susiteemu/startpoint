/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	requestchain "startpoint/core/chaining"
	"startpoint/core/client/runner"
	"startpoint/core/loader"
	"startpoint/core/model"
	"startpoint/core/print"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RunConfig struct {
	Plain          bool
	PrintHeaders   bool
	PrintBody      bool
	PrintTraceInfo bool
}

type RunArgs struct {
	Request string
	Profile string
}

var runConfig RunConfig

var runCmd = &cobra.Command{
	Use:   "run [REQUEST NAME] [PROFILE NAME]",
	Short: "Run a http request from workspace",
	Long:  `Run a http request from workspace`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
			return err
		}

		parsedArgs := ParseArgs(args)

		if len(parsedArgs.Request) == 0 {
			return errors.New("Request name is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		runArgs := ParseArgs(args)

		requests, err := loader.ReadRequests(viper.GetString("workspace"))
		if err != nil {
			fmt.Print(fmt.Errorf("error %v", err))
			return
		}
		var request *model.RequestMold
		for _, m := range requests {
			if m.Name == runArgs.Request {
				request = m
				break
			}
		}
		if request == nil {
			fmt.Printf("Could not find a request with name '%s' under workspace '%s'", runArgs.Request, viper.GetString("workspace"))
			return
		}

		profiles, err := loader.ReadProfiles(viper.GetString("workspace"))
		if err != nil {
			fmt.Print(fmt.Errorf("error %v", err))
			return
		}
		profileName := runArgs.Profile
		if len(profileName) == 0 {
			profileName = "default"
		}
		var profile *model.Profile
		envVars := os.Environ()
		for _, p := range profiles {
			if p.Name == profileName {
				profile = &model.Profile{
					Name:      p.Name,
					Variables: loader.GetProfileValues(p, profiles, envVars),
				}
				break
			}
		}

		runRequests := requestchain.ResolveRequestChain(request, requests)
		responses, err := runner.RunRequestChain(runRequests, profile, func(took time.Duration, statusCode int) {
			log.Info().Msgf("Request responded with status %d and took %s", statusCode, took)
		})
		if err != nil {
			fmt.Print(fmt.Errorf("error %v", err))
			return
		}

		for _, response := range responses {
			printOpts := print.PrintOpts{
				PrettyPrint:    !runConfig.Plain,
				PrintHeaders:   runConfig.PrintHeaders,
				PrintBody:      runConfig.PrintBody,
				PrintTraceInfo: runConfig.PrintTraceInfo,
			}
			responseStr, err := print.SprintResponse(response, printOpts)
			if err != nil {
				fmt.Print(fmt.Errorf("error %v", err))
				return
			}
			fmt.Println(responseStr)
		}

	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		log.Debug().Msgf("ValidArgsFunction args=%v, toComplete=%s, w=%s", args, toComplete, viper.GetString("workspace"))
		if len(args) == 2 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		suggestions := []string{}
		if len(args) == 0 {
			requests, _ := loader.ReadRequests(viper.GetString("workspace"))
			if requests == nil {
				requests = []*model.RequestMold{}
			}
			for _, req := range requests {
				if strings.Contains(req.Name, toComplete) {
					suggestions = append(suggestions, req.Name)
				}
			}
		} else if len(args) == 1 {
			profiles, _ := loader.ReadProfiles(viper.GetString("workspace"))
			if profiles == nil {
				profiles = []*model.Profile{}
			}
			for _, p := range profiles {
				if strings.Contains(p.Name, toComplete) {
					suggestions = append(suggestions, p.Name)
				}
			}
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
}

func ParseArgs(args []string) RunArgs {
	if len(args) == 0 {
		return RunArgs{}
	} else if len(args) == 1 {
		return RunArgs{args[0], ""}
	}
	return RunArgs{args[0], args[1]}
}

func init() {
	rootCmd.AddCommand(runCmd)

	const printTrace = "t"
	const printHeadersP = "h"
	const printBodyP = "b"

	runConfig.PrintBody = true

	runCmd.PersistentFlags().BoolVarP(&runConfig.Plain, "plain", "p", false, "Print plain response without styling")
	runCmd.PersistentFlags().Bool("no-body", false, "Print no body")
	runCmd.PersistentFlags().StringSlice("print", []string{}, fmt.Sprintf("Print WHAT\n- '%s'\tPrint response headers\n- '%s'\tPrint response body\n- '%s'\tPrint trace information", printHeadersP, printBodyP, printTrace))
	runCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if cmd == runCmd {
			printFlags, _ := cmd.Flags().GetStringSlice("print")
			for _, flag := range printFlags {
				if flag == printHeadersP {
					runConfig.PrintHeaders = true
				} else if flag == printBodyP {
					runConfig.PrintBody = true
				} else if flag == printTrace {
					runConfig.PrintTraceInfo = true
				}
			}
			noBody, _ := cmd.PersistentFlags().GetBool("no-body")
			runConfig.PrintBody = !noBody
		}
	}

}
