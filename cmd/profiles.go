/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"startpoint/core/loader"
	profileUI "startpoint/tui/profilemgmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manageProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		workspace := viper.GetString("workspace")
		loadedProfiles, _ := loader.ReadProfiles(workspace)
		profileUI.Start(loadedProfiles)
	},
}

func init() {
	rootCmd.AddCommand(manageProfilesCmd)
}
