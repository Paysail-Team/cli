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
			fmt.Println("Something went wrong while loading items from your project's config file.")
		}

		if len(token) > 0 && len(projectId) > 0 {
			companyId, validatedAccessToken, validatedProjectId, err := api.ValidateAccessTokenAndProjectId(token, projectId, env, customServerUrl)
			if err != nil {
				fmt.Println("Something went wrong while validating your access token and project ID.")
			}
			util.SaveCompanyIdAndValidatedInfoToProjectConfigFile(companyId, validatedAccessToken, validatedProjectId)
			return
		}

		if isConfigLoaded == true {
			isSessionTokenValid, err := api.ValidateSessionToken(config.SessionToken, env, customServerUrl)
			if err != nil {
				fmt.Println("Something went wrong while validating your session token. We recommend logging in again.")
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
				fmt.Println("Your session token is invalid. You'll need to log in again.")
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
					fmt.Println("You're all set!\nTo learn more about the app and the CLI commands it offers, you may enter `Syro --help` or `Syro [command] --help`.")
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
	fmt.Println("Please enter your credentials.")
	email, password, err := getLoginCredentials()
	if err != nil {
		fmt.Println("Something went wrong while getting your credentials. Please try again.")
		return "", "", "", err
	}

	companyId, expiresAt, memberId, sessionToken, err := api.Login(email, password, env, customServerUrl)
	if err != nil {
		fmt.Println("Invalid email/password combination. Please input the correct credentials.")
		return "", "", "", err
	}
	fmt.Println("Login successful!")

	err = util.SaveUserAndSessionInfoToProjectConfigFile(companyId, expiresAt, memberId, sessionToken)
	if err != nil {
		fmt.Println("Something went wrong while saving user and session info to your project's config file.")
		return "", "", "", err
	}
	return companyId, memberId, sessionToken, nil
}

func getProjectIdAndUpdateProjectConfigFile(companyId string, memberId string, sessionToken string, env string, customServerUrl string) (projectId string, err error) {
	fmt.Println("Please enter the project ID of a project you own or shared with you.")
	userProjectId, err := util.GetProjectId()
	if err != nil {
		fmt.Println("Something went wrong while getting your project ID. Please try again.")
		return "", err
	}

	isProjectIdValid, err := api.ValidateProjectId(companyId, memberId, userProjectId, sessionToken, env, customServerUrl)
	if !isProjectIdValid {
		fmt.Println("The project ID you entered is not associated with any project you own or shared with you. Please try again.")
		return "", err
	} else {
		fmt.Println("Project ID Validated!")
		err = util.SaveProjectIdToProjectConfigFile(userProjectId)
		if err != nil {
			fmt.Println("Something went wrong while saving user and session info to your project's config file.")
			return "", err
		}
	}
	fmt.Println("You're all set!\nTo learn more about the app and the CLI commands it offers, you may enter `Syro --help` or `Syro [command] --help`.")
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
