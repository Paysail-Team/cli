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

func ValidateProjectId(companyId string, memberId string, projectId string, sessionToken string, env string, customServerUrl string) (isProjectIdValid bool, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"companyId":"%s", "memberId":"%s", "projectId":"%s", "sessionToken":"%s"}`, companyId, memberId, projectId, sessionToken)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliValidateProjectId))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return false, err
	}

	validateProjectIdResponse := model.ValidateProjectIdResponse{}
	if err := json.Unmarshal(response.Body(), &validateProjectIdResponse); err != nil {
		fmt.Printf("Unable to unmarshal response from the request.")
		fmt.Printf("Error :: %v\n", err)
		return false, err
	}

	if len(validateProjectIdResponse.Error) > 0 {
		fmt.Print("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", validateProjectIdResponse.Error)
		return false, errors.New(validateProjectIdResponse.Error)
	}

	return validateProjectIdResponse.Result.IsProjectIdValid, nil
}

func FetchProjectItems(companyId string, memberId string, projectId string, sessionToken string, env string, customServerUrl string) (secrets []model.ItemDetails, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"companyId":"%s", "memberId":"%s", "projectId":"%s", "sessionToken":"%s"}`, companyId, memberId, projectId, sessionToken)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliProjectItems))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return []model.ItemDetails{}, err
	}

	fetchProjectSecretsResponse := model.FetchProjectSecretsResponse{}
	if err := json.Unmarshal(response.Body(), &fetchProjectSecretsResponse); err != nil {
		fmt.Printf("Unable to unmarshal response from the request.")
		fmt.Printf("Error :: %v\n", err)
		return []model.ItemDetails{}, err
	}

	if len(fetchProjectSecretsResponse.Error) > 0 {
		fmt.Print("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", fetchProjectSecretsResponse.Error)
		return []model.ItemDetails{}, errors.New(fetchProjectSecretsResponse.Error)
	}

	return fetchProjectSecretsResponse.Result.Items, nil
}

func PullProjectItems(accessToken string, companyId string, projectId string, env string, customServerUrl string) (secrets []model.ItemDetails, err error) {
	client := resty.New()
	requestBody := fmt.Sprintf(`{"accessToken":"%s", "companyId":"%s", "projectId":"%s"}`, accessToken, companyId, projectId)
	serverApiUrl := util.GetServerUrl(env, customServerUrl)
	serverApplicationId := util.GetServerApplicationId(env)

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Parse-Application-Id", serverApplicationId).
		SetBody(requestBody).
		Post(fmt.Sprintf("%s%s", serverApiUrl, util.CliPullProjectItems))

	if err != nil {
		fmt.Printf("Failed to connect to the server.\n")
		fmt.Printf("Error :: %v\n", err)
		return []model.ItemDetails{}, err
	}

	pullProjectSecretsResponse := model.FetchProjectSecretsResponse{}
	if err := json.Unmarshal(response.Body(), &pullProjectSecretsResponse); err != nil {
		fmt.Printf("Unable to unmarshal response from the request.")
		fmt.Printf("Error :: %v\n", err)
		return []model.ItemDetails{}, err
	}

	if len(pullProjectSecretsResponse.Error) > 0 {
		fmt.Print("The server responded with an error.\n")
		fmt.Printf("Error :: %v\n", pullProjectSecretsResponse.Error)
		return []model.ItemDetails{}, errors.New(pullProjectSecretsResponse.Error)
	}

	return pullProjectSecretsResponse.Result.Items, nil
}
