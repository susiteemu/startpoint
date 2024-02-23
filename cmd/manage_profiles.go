/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"goful/core/loader"
	profileManageTui "goful/tui/profile/manage"

	"github.com/spf13/cobra"
)

// manageCmd represents the manage command
var manageProfilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO handle err
		loadedProfiles, _ := loader.ReadProfiles("tmp")
		profileManageTui.Start(loadedProfiles)
	},
}

func init() {
	manageCmd.AddCommand(manageProfilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// manageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// manageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
