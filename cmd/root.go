/*
Copyright © 2023 Teemu Turunen <teturun@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goful-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
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
	runLogFile, _ := os.OpenFile(
		"app.log",
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
	rootCmd.PersistentFlags().Bool("help", false, "Displays help")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("workspace", "w", "", "Workspace directory")
	viper.BindPFlag("workspace", rootCmd.PersistentFlags().Lookup("workspace"))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goful.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".goful")
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
