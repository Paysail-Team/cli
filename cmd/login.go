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

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to the Syro CLI",
	Long:  "Input your credentials to authenticate and gain access to the secrets in your projects in Syro.",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := cmd.Flags().GetString("token")
		projectId, _ := cmd.Flags().GetString("projectId")
		env, _ := cmd.Flags().GetString("env")
		customServerUrl, _ := cmd.Flags().GetString("serverUrl")

		isConfigLoaded, config, err := util.LoadConfigFromProjectConfigFile()
		if err != nil {
			fmt.Printf("Failed to load items from config file.\n")
		}

		if len(token) > 0 && len(projectId) > 0 {
			companyId, validatedAccessToken, validatedProjectId, err := api.ValidateAccessTokenAndProjectId(token, projectId, env, customServerUrl)
			if err != nil {
				fmt.Printf("Failed to validate your access token and project ID.\n")
			}
			util.SaveCompanyIdAndValidatedInfoToProjectConfigFile(companyId, validatedAccessToken, validatedProjectId)
			return
		}

		if isConfigLoaded == true {
			isSessionTokenValid, err := api.ValidateSessionToken(config.SessionToken, env, customServerUrl)
			if err != nil {
				fmt.Printf("Failed to validate your session token. We recommend logging in again.\n")
				companyId, memberId, sessionToken, err := loginAndUpdateProjectConfigFile(env, customServerUrl)
				if err != nil {
					return
				}
				_, err = getProjectIdAndUpdateProjectConfigFile(companyId, memberId, sessionToken, env, customServerUrl)
				if err != nil {
					return
				}
			}
			if !isSessionTokenValid {
				fmt.Printf("Your session token is invalid. You'll need to log in again.\n")
				companyId, memberId, sessionToken, err := loginAndUpdateProjectConfigFile(env, customServerUrl)
				if err != nil {
					return
				}
				_, err = getProjectIdAndUpdateProjectConfigFile(companyId, memberId, sessionToken, env, customServerUrl)
				if err != nil {
					return
				}
			} else {
				if len(config.ProjectId) == 0 {
					_, err = getProjectIdAndUpdateProjectConfigFile(config.CompanyId, config.MemberId, config.SessionToken, env, customServerUrl)
					if err != nil {
						return
					}
				} else {
					fmt.Printf("You're all set!\nTo learn more about the app and the CLI commands it offers, you may enter `Syro --help` or `Syro [command] --help`.\n")
				}
			}
		} else {
			companyId, memberId, sessionToken, err := loginAndUpdateProjectConfigFile(env, customServerUrl)
			if err != nil {
				return
			}
			_, err = getProjectIdAndUpdateProjectConfigFile(companyId, memberId, sessionToken, env, customServerUrl)
			if err != nil {
				return
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("token", "t", "", "Specify the access token")
	loginCmd.Flags().StringP("projectId", "p", "", "Specify the project ID")
	loginCmd.MarkFlagsRequiredTogether("token", "projectId")
}

func loginAndUpdateProjectConfigFile(env string, customServerUrl string) (companyId string, memberId string, sessionToken string, err error) {
	fmt.Printf("Please enter your credentials.\n")
	email, password, err := getLoginCredentials()
	if err != nil {
		fmt.Printf("Failed to get your credentials. Please try again.\n")
		return "", "", "", err
	}

	companyId, expiresAt, memberId, sessionToken, err := api.Login(email, password, env, customServerUrl)
	if err != nil {
		fmt.Printf("Login failed!\n")
		return "", "", "", err
	}
	fmt.Print("Login successful!\n")

	err = util.SaveUserAndSessionInfoToProjectConfigFile(companyId, expiresAt, memberId, sessionToken)
	if err != nil {
		fmt.Printf("Failed to save user and session info to your config file.\n")
		return "", "", "", err
	}
	return companyId, memberId, sessionToken, nil
}

func getProjectIdAndUpdateProjectConfigFile(companyId string, memberId string, sessionToken string, env string, customServerUrl string) (projectId string, err error) {
	fmt.Printf("Please enter the project ID of a project you own or shared with you.\n")
	userProjectId, err := util.GetProjectId()
	if err != nil {
		fmt.Printf("Failed to get your project ID. Please try again.\n")
		return "", err
	}

	isProjectIdValid, err := api.ValidateProjectId(companyId, memberId, userProjectId, sessionToken, env, customServerUrl)
	if !isProjectIdValid {
		fmt.Printf("Failed to validate project ID. The project ID you entered may not be associated with any project you own or shared with you. Please try again.\n")
		return "", err
	} else {
		fmt.Printf("Project ID Validated!\n")
		err = util.SaveProjectIdToProjectConfigFile(userProjectId)
		if err != nil {
			fmt.Printf("Failed to save user and session info to config file.\n")
			return "", err
		}
	}
	fmt.Printf("You're all set!\nTo learn more about the app and the CLI commands it offers, you may enter `Syro --help` or `Syro [command] --help`.\n")
	return userProjectId, nil
}

func getLoginCredentials() (email string, password string, err error) {
	userEmail, userEmailErr := util.AskForUserEmail()

	if userEmailErr != nil {
		return "", "", userEmailErr
	}

	userPassword, userPasswordErr := util.AskForUserPassword()

	if userPasswordErr != nil {
		return "", "", userPasswordErr
	}

	return userEmail, userPassword, nil
}
