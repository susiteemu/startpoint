/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"goful/core/loader"
	requestUI "goful/tui/request"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
	Use:   "requests",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		log.Info().Msg("Starting to handle requests cmd")
		// TODO handle err
		loadedRequests, _ := loader.ReadRequests("tmp")
		loadedProfiles, _ := loader.ReadProfiles("tmp")
		log.Info().Msgf("Loaded %d requests", len(loadedRequests))
		log.Info().Msg("Starting up ui...")
		requestUI.Start(loadedRequests, loadedProfiles)
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
