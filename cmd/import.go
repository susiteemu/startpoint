/*
Copyright © 2024 Teemu Turunen <teturun@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/susiteemu/startpoint/core/import/openapi"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ImportConfig struct {
	Path string
}

var importConfig ImportConfig

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import workspace from OpenAPI Spec v3",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Called import cmd with %v", importConfig)
		workspace := viper.GetString("workspace")
		openapi.ReadSpec(importConfig.Path, workspace)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().StringVarP(&importConfig.Path, "path", "p", "", "OpenAPI Spec v3 location (filepath or url)")
}
