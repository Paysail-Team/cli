/*
Copyright © 2023 Syro team <info@syro.com>
*/
package model

type Config struct {
	CompanyId            string `json:"companyId"`
	ExpiresAt            string `json:"expiresAt"`
	MemberId             string `json:"memberId"`
	ProjectId            string `json:"projectId"`
	SessionToken         string `json:"sessionToken"`
	ValidatedAccessToken string `json:"validatedAccessToken"`
	ValidatedProjectId   string `json:"validatedProjectId"`
}

func (config *Config) UpdateSessionInfo(expiresAt string, sessionToken string) {
	config.ExpiresAt = expiresAt
	config.SessionToken = sessionToken
}

func (config *Config) UpdateMembershipInfo(companyId string, memberId string) {
	config.CompanyId = companyId
	config.MemberId = memberId
}

func (config *Config) UpdateProjectId(projectId string) {
	config.ProjectId = projectId
}
