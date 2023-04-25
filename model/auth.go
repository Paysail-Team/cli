/*
Copyright © 2023 Syro team <info@syro.com>
*/
package model

type LoginResponseResult struct {
	CompanyId    string `json:"companyId"`
	ExpiresAt    string `json:"expiresAt"`
	MemberId     string `json:"memberId"`
	SessionToken string `json:"sessionToken"`
}

type LoginResponse struct {
	Result LoginResponseResult `json:"result"`
	Error  string              `json:"error"`
	Code   int                 `json:"code"`
}

type ValidateAccessTokenAndProjectIdResponseResult struct {
	CompanyId           string `json:"companyId"`
	VerifiedAccessToken string `json:"verifiedAccessToken"`
	VerifiedProjectId   string `json:"verifiedProjectId"`
}

type ValidateAccessTokenAndProjectIdResponse struct {
	Result ValidateAccessTokenAndProjectIdResponseResult `json:"result"`
	Error  string                                        `json:"error"`
	Code   int                                           `json:"code"`
}

type ValidateSessionTokenResponseResult struct {
	IsSessionTokenValid bool `json:"isSessionTokenValid"`
}

type ValidateSessionTokenResponse struct {
	Result ValidateSessionTokenResponseResult `json:"result"`
	Error  string                             `json:"error"`
	Code   int                                `json:"code"`
}
