/*
Copyright © 2023 Syro team <info@syro.com>
*/

package cmd

import (
	"fmt"
	"syro/api"
	"syro/util"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull secrets from your Syro project",
	Long:  "Fetch the latest secrets from your chosen project and inject it to your app via an env file.",
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		customServerUrl, _ := cmd.Flags().GetString("serverUrl")

		isConfigLoaded, config, err := util.LoadConfigFromProjectConfigFile()
		if err != nil {
			fmt.Printf("Failed to load items from config file.\n")
		}

		if isConfigLoaded == true {
			items, err := api.PullProjectItems(config.ValidatedAccessToken, config.CompanyId, config.ValidatedProjectId, env, customServerUrl)
			if err != nil {
				fmt.Printf("Failed to pull project secrets.\n")
				return
			}

			err = util.SaveSecretsToEnvFile(items)
			if err != nil {
				fmt.Print("Failed to save project secrets to env file.\n")
				return
			}

		} else {
			fmt.Printf("The Syro CLI is not properly configured yet for this project. Kindly complete the set up first by using the login command.\n")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
