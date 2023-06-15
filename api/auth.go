/*
Copyright © 2023 Syro team <info@syro.com>
*/
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"syro/model"
	"syro/util"

	"github.com/go-resty/resty/v2"
)

func Login(email string, password string, env string, customServerUrl string) (expiresAt string, memberships []model.MembershipDetails, sessionToken string, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliLogin))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return "", []model.MembershipDetails{}, "", err
	}

	loginResponse := model.LoginResponse{}
	if err := json.Unmarshal(response.Body(), &loginResponse); err != nil {
		fmt.Print("Unable to unmarshal response from authentication.\n")
		fmt.Printf("Error :: %v\n", err)
		return "", []model.MembershipDetails{}, "", err
	}

	if len(loginResponse.Error) > 0 {
		fmt.Printf("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", loginResponse.Error)
		return "", []model.MembershipDetails{}, "", errors.New(loginResponse.Error)
	}

	return loginResponse.Result.ExpiresAt, loginResponse.Result.Memberships, loginResponse.Result.SessionToken, nil
}

func ValidateAccessTokenAndProjectId(accessToken string, projectId string, env string, customServerUrl string) (companyId string, verifiedAccessToken string, verifiedProjectId string, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"accessToken":"%s", "projectId":"%s"}`, accessToken, projectId)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliValidateAccessTokenAndProjectId))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return "", "", "", err
	}

	validateAccessTokenAndProjectIdResponse := model.ValidateAccessTokenAndProjectIdResponse{}
	if err := json.Unmarshal(response.Body(), &validateAccessTokenAndProjectIdResponse); err != nil {
		fmt.Printf("Could not unmarshal response from validate access token and project ID request.\n")
		fmt.Printf("Error :: %v\n", err)
		return "", "", "", err
	}

	if len(validateAccessTokenAndProjectIdResponse.Error) > 0 {
		fmt.Print("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", validateAccessTokenAndProjectIdResponse.Error)
		return "", "", "", errors.New(validateAccessTokenAndProjectIdResponse.Error)
	}

	return validateAccessTokenAndProjectIdResponse.Result.CompanyId, validateAccessTokenAndProjectIdResponse.Result.VerifiedAccessToken, validateAccessTokenAndProjectIdResponse.Result.VerifiedProjectId, nil
}

func ValidateSessionToken(sessionToken string, env string, customServerUrl string) (isSessionTokenValid bool, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"sessionToken":"%s"}`, sessionToken)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliValidateSessionToken))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return false, err
	}

	validateSessionTokenResponse := model.ValidateSessionTokenResponse{}
	if err := json.Unmarshal(response.Body(), &validateSessionTokenResponse); err != nil {
		fmt.Printf("Could not unmarshal response from validate session token request.\n")
		fmt.Printf("Error :: %v\n", err)
		return false, err
	}

	if len(validateSessionTokenResponse.Error) > 0 {
		fmt.Printf("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", validateSessionTokenResponse.Error)
		return false, errors.New(validateSessionTokenResponse.Error)
	}

	return validateSessionTokenResponse.Result.IsSessionTokenValid, nil
}

func ValidateMemberId(memberId string, sessionToken string, env string, customServerUrl string) (isMemberIdValid bool, companyId string, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"memberId":"%s", "sessionToken":"%s"}`, memberId, sessionToken)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliValidateMemberId))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return false, "", err
	}

	validateMemberIdResponse := model.ValidateMemberIdResponse{}
	if err := json.Unmarshal(response.Body(), &validateMemberIdResponse); err != nil {
		fmt.Printf("Unable to unmarshal response from the request.")
		fmt.Printf("Error :: %v\n", err)
		return false, "", err
	}

	if len(validateMemberIdResponse.Error) > 0 {
		fmt.Print("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", validateMemberIdResponse.Error)
		return false, "", errors.New(validateMemberIdResponse.Error)
	}

	return validateMemberIdResponse.Result.IsMemberIdValid, validateMemberIdResponse.Result.CompanyId, nil
}

func FetchUserMemberships(sessionToken string, env string, customServerUrl string) (memberships []model.MembershipDetails, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"sessionToken":"%s"}`, sessionToken)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliUserMemberships))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return []model.MembershipDetails{}, err
	}

	fetchUserMembershipsResponse := model.FetchUserMembershipsResponse{}
	if err := json.Unmarshal(response.Body(), &fetchUserMembershipsResponse); err != nil {
		fmt.Printf("Unable to unmarshal response from the request.")
		fmt.Printf("Error :: %v\n", err)
		return []model.MembershipDetails{}, err
	}

	if len(fetchUserMembershipsResponse.Error) > 0 {
		fmt.Print("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", fetchUserMembershipsResponse.Error)
		return []model.MembershipDetails{}, errors.New(fetchUserMembershipsResponse.Error)
	}

	return fetchUserMembershipsResponse.Result.Memberships, nil
}
