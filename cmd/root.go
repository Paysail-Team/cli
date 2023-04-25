/*
Copyright © 2023 Syro team <info@syro.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "syro",
	Short: "Easily access and inject your secrets through the Syro CLI",
	Long:  "Syro CLI allows you to access secrets in projects you own or are shared with you and inject those secrets into your CI/CD pipelines.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("env", "", "this has been automatically set for you")
	rootCmd.PersistentFlags().String("serverUrl", "", "this has been automatically set for you")
}
