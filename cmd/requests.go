/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goful/core/loader"
	requestUI "goful/tui/request"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manageCmd = &cobra.Command{
	Use:   "requests",
	Short: "Starts up tui application to manage and run requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Starting to handle requests cmd with workspace root %s", viper.GetString("workspace"))
		workspace := viper.GetString("workspace")
		loadedRequests, err := loader.ReadRequests(workspace)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read requests")
			fmt.Printf("Failed to read requests %v", err)
			return
		}
		loadedProfiles, err := loader.ReadProfiles(workspace)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read profiles")
			fmt.Printf("Failed to read profiles %v", err)
			return
		}
		log.Info().Msgf("Loaded %d requests and %d profiles", len(loadedRequests), len(loadedProfiles))
		log.Info().Msg("Starting up ui...")
		requestUI.Start(loadedRequests, loadedProfiles)
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
}
