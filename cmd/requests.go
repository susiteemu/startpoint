package cmd

import (
	mainview "github.com/susiteemu/startpoint/tui"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manageCmd = &cobra.Command{
	Use:   "requests",
	Short: "Start up a tui application to manage and run requests",
	Long:  "Start up a tui application to manage and run requests",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Starting to handle requests cmd with workspace root %s", viper.GetString("workspace"))
		log.Debug().Msgf("All configuration values %v", viper.AllSettings())
		workspace := viper.GetString("workspace")
		mainview.Start(workspace, mainview.Requests)
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
}
