/*
Copyright © 2023 Syro team <info@syro.com>
*/
package model

type ValidateProjectIdResponseResult struct {
	IsProjectIdValid bool `json:"isProjectIdValid"`
}

type ValidateProjectIdResponse struct {
	Result ValidateProjectIdResponseResult `json:"result"`
	Error  string                          `json:"error"`
	Code   int                             `json:"code"`
}

type ItemDetails struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type FetchProjectSecretsResponseResult struct {
	Items []ItemDetails `json:"items"`
}

type FetchProjectSecretsResponse struct {
	Result FetchProjectSecretsResponseResult `json:"result"`
	Error  string                            `json:"error"`
	Code   int                               `json:"code"`
}
