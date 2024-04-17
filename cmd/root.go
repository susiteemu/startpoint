/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "startpoint",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Root: Config file used: %v\n", viper.ConfigFileUsed())
		fmt.Printf("Root: All keys: %v\n", viper.AllKeys())
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

	rootCmd.PersistentFlags().StringP("workspace", "w", "", "Workspace directory")
	viper.BindPFlag("workspace", rootCmd.PersistentFlags().Lookup("workspace"))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.startpoint.yaml)")

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

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("theme.syntax", "native")
	viper.SetDefault("printer.response.formatter", "terminal16m")

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	cobra.CheckErr(err)
	viper.SetDefault("workspace", cwd)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config %v\n", err)
	}

}
