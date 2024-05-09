/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"startpoint/core/loader"
	profileUI "startpoint/tui/profilemgmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manageProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Start up a tui application to manage profiles",
	Long:  "Start up a tui application to manage profiles",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Starting to handle profiles cmd with workspace root %s", viper.GetString("workspace"))
		log.Debug().Msgf("All configuration values %v", viper.AllSettings())
		workspace := viper.GetString("workspace")
		loadedProfiles, err := loader.ReadProfiles(workspace)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read profiles")
			fmt.Printf("Failed to read profiles %v", err)
			return
		}
		log.Info().Msgf("Loaded %d profiles", len(loadedProfiles))
		log.Info().Msg("Starting up ui...")
		profileUI.Start(loadedProfiles)
	},
}

func init() {
	rootCmd.AddCommand(manageProfilesCmd)
}
