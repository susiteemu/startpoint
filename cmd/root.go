/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "startpoint",
	Short: "Startpoint is a tui app for managing and running HTTP requests",
	Long: `Startpoint is a tui app with which you can manage and run HTTP requests from your terminal. It offers a way for flexible chaining and scripting requests as well as defining them in a simple format.
	`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// disable for now, before custom completion is implemented
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().Bool("help", false, "Displays help")

	rootCmd.PersistentFlags().StringP("workspace", "w", "", "Workspace directory (default is current dir)")
	viper.BindPFlag("workspace", rootCmd.PersistentFlags().Lookup("workspace"))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is a merge of $HOME/.startpoint.yaml and <workspace>/.startpoint.yaml)")

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	logFile := filepath.Join(home, "startpoint.log")
	runLogFile, _ := os.OpenFile(logFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		multi := zerolog.MultiLevelWriter(runLogFile)
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
		log.Info().Msg("Initialized logging")
		return nil
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".startpoint")

	}
	cwd, err := os.Getwd()
	cobra.CheckErr(err)
	viper.SetDefault("workspace", cwd)
	// TODO: where should default values come from?
	viper.SetDefault("printer.response.formatter", "terminal16m")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config %v\n", err)
	}

	// if no specific config file is given, we try to merge config files found in $HOME and workspace
	if cfgFile == "" {
		workspaceViper := viper.New()
		workspaceViper.AddConfigPath(viper.GetString("workspace"))
		workspaceViper.SetConfigType("yaml")
		workspaceViper.SetConfigName(".startpoint")

		// If a config file is found, read it in.
		if err := workspaceViper.ReadInConfig(); err == nil {
			viper.MergeConfigMap(workspaceViper.AllSettings())
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

}
