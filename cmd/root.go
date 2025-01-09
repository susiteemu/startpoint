/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"os"
	"path/filepath"
	"startpoint/core/configuration"
	"startpoint/core/writer"
	"strings"

	mainview "startpoint/tui"

	"embed"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed config/.startpoint-default.yaml
var config embed.FS

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "startpoint",
	Version: version + " (commit=" + commit + ", build date=" + date + ")",
	Short:   "Startpoint is a tui app for managing and running HTTP requests",
	Long: `Startpoint is a tui app with which you can manage and run HTTP requests from your terminal. It offers a way for flexible chaining and scripting requests as well as defining them in a simple format.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace := viper.GetString("workspace")
		mainview.Start(workspace, mainview.Requests)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
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
		// Default level for this example is info, unless debug flag is present
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if configuration.New().GetBoolWithDefault("debug", false) {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		log.Info().Msgf("Initialized logging with level %s", zerolog.GlobalLevel().String())
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

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if cfgFile == "" {
			// If no specific config file is given, read embedded default config file and write it to $HOME and read it
			file, err := config.ReadFile("config/.startpoint-default.yaml")
			cobra.CheckErr(err)
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			_, err = writer.WriteFile(filepath.Join(home, ".startpoint.yaml"), string(file))
			cobra.CheckErr(err)

			if err := viper.ReadInConfig(); err != nil {
				cobra.CheckErr(err)
			}
		} else {
			cobra.CheckErr(err)
		}
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
