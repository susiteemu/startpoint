package cmd

import (
	mainview "github.com/susiteemu/startpoint/tui"

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
		mainview.Start(workspace, mainview.Profiles)
	},
}

func init() {
	rootCmd.AddCommand(manageProfilesCmd)
}
