/*
Copyright © 2023 Syro team <info@syro.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"syro/api"
	"syro/model"
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
				memberships, sessionToken, err := loginAndUpdateProjectConfigFile(env, customServerUrl)
				if err != nil {
					return
				}
				companyId, memberId, err := getMembershipSelectionAndUpdateProjectConfigFile(memberships, sessionToken, env, customServerUrl)
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
				memberships, sessionToken, err := loginAndUpdateProjectConfigFile(env, customServerUrl)
				if err != nil {
					return
				}
				companyId, memberId, err := getMembershipSelectionAndUpdateProjectConfigFile(memberships, sessionToken, env, customServerUrl)
				if err != nil {
					return
				}
				_, err = getProjectIdAndUpdateProjectConfigFile(companyId, memberId, sessionToken, env, customServerUrl)
				if err != nil {
					return
				}
			} else {
				if len(config.CompanyId) == 0 {
					memberships, err := getMemberships(config.SessionToken, env, customServerUrl)
					if err != nil {
						return
					}
					companyId, memberId, err := getMembershipSelectionAndUpdateProjectConfigFile(memberships, config.SessionToken, env, customServerUrl)
					if err != nil {
						return
					}
					_, err = getProjectIdAndUpdateProjectConfigFile(companyId, memberId, config.SessionToken, env, customServerUrl)
					if err != nil {
						return
					}
				} else if len(config.ProjectId) == 0 {
					_, err = getProjectIdAndUpdateProjectConfigFile(config.CompanyId, config.MemberId, config.SessionToken, env, customServerUrl)
					if err != nil {
						return
					}
				} else {
					fmt.Printf("You're all set!\nTo learn more about the app and the CLI commands it offers, you may enter `Syro --help` or `Syro [command] --help`.\n")
				}
			}
		} else {
			memberships, sessionToken, err := loginAndUpdateProjectConfigFile(env, customServerUrl)
			if err != nil {
				return
			}
			companyId, memberId, err := getMembershipSelectionAndUpdateProjectConfigFile(memberships, sessionToken, env, customServerUrl)
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

func loginAndUpdateProjectConfigFile(env string, customServerUrl string) (memberships []model.MembershipDetails, sessionToken string, err error) {
	fmt.Printf("Please enter your credentials.\n")
	email, password, err := getLoginCredentials()
	if err != nil {
		fmt.Printf("Failed to get your credentials. Please try again.\n")
		return []model.MembershipDetails{}, "", err
	}

	expiresAt, memberships, sessionToken, err := api.Login(email, password, env, customServerUrl)
	if err != nil {
		fmt.Printf("Login failed!\n")
		return []model.MembershipDetails{}, "", err
	}
	fmt.Print("Login successful!\n")

	err = util.SaveSessionInfoToProjectConfigFile(expiresAt, sessionToken)
	if err != nil {
		fmt.Printf("Failed to save session info to your config file.\n")
		return []model.MembershipDetails{}, "", err
	}
	return memberships, sessionToken, nil
}

func getMemberships(sessionToken string, env string, customServerUrl string) (memberships []model.MembershipDetails, err error) {
	memberships, err = api.FetchUserMemberships(sessionToken, env, customServerUrl)
	if err != nil {
		fmt.Printf("Failed to fetch memberships. Please try again.\n")
		return []model.MembershipDetails{}, err
	}
	return memberships, nil
}

func getMembershipSelectionAndUpdateProjectConfigFile(memberships []model.MembershipDetails, sessionToken string, env string, customServerUrl string) (companyId string, memberId string, err error) {
	selectedMemberId := ""
	if len(memberships) == 0 {
		fmt.Printf("You aren't a member of any company. Your membership from a company may have been revoked or you haven't signed up yet for an account.\n")
		return "", "", errors.New("No memberships from any company.")
	} else if len(memberships) == 1 {
		selectedMemberId = memberships[0].MemberId
	} else {
		fmt.Printf("You are a member of multiple companies. Select a company to continue.\n")
		selectedMemberId, err = util.GetMembershipSelection(memberships)
		if err != nil {
			fmt.Printf("Failed to select a company. Please try again.\n")
			return "", "", err
		}
	}

	isMemberIdValid, companyId, err := api.ValidateMemberId(selectedMemberId, sessionToken, env, customServerUrl)
	if !isMemberIdValid {
		if len(memberships) == 1 {
			fmt.Printf("Failed to validate member ID. Please try again.\n")
		} else {
			fmt.Printf("Failed to validate selected company. The company you selected may have revoked your membership. Please try again.\n")
			return "", "", err
		}
	} else {
		if len(memberships) > 1 {
			fmt.Printf("Company seletion successful!\n")
		}
		err = util.SaveCompanyIdAndMemberIdToProjectConfigFile(companyId, selectedMemberId)
		if err != nil {
			fmt.Printf("Failed to save company ID and member ID to config file.\n")
			return "", "", err
		}
	}
	return companyId, selectedMemberId, nil
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
