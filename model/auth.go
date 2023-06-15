/*
Copyright © 2023 Syro team <info@syro.com>
*/
package model

type MembershipDetails struct {
	MemberId    string `json:"memberId"`
	CompanyName string `json:"companyName"`
}
type LoginResponseResult struct {
	ExpiresAt    string              `json:"expiresAt"`
	Memberships  []MembershipDetails `json:"memberships"`
	SessionToken string              `json:"sessionToken"`
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

type ValidateMemberIdResponseResult struct {
	IsMemberIdValid bool   `json:"isMemberIdValid"`
	CompanyId       string `json:"companyId"`
}

type ValidateMemberIdResponse struct {
	Result ValidateMemberIdResponseResult `json:"result"`
	Error  string                         `json:"error"`
	Code   int                            `json:"code"`
}

type FetchUserMembershipsResponseResult struct {
	Memberships []MembershipDetails `json:"memberships"`
}

type FetchUserMembershipsResponse struct {
	Result FetchUserMembershipsResponseResult `json:"result"`
	Error  string                             `json:"error"`
	Code   int                                `json:"code"`
}
