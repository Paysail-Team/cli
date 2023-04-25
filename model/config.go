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

func (config *Config) UpdateUserAndSessionInfo(companyId string, expiresAt string, memberId string, sessionToken string) {
	config.CompanyId = companyId
	config.ExpiresAt = expiresAt
	config.MemberId = memberId
	config.SessionToken = sessionToken
}

func (config *Config) UpdateProjectId(projectId string) {
	config.ProjectId = projectId
}
