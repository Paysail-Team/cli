/*
Copyright © 2023 Syro team <info@syro.com>
*/
package util

import "fmt"

// Server URL
var LocalServerApiUrl = "http://localhost:1400/v1/functions"
var StagingServerApiUrl = "https://api-staging.syro.com/v1/functions"
var ProductionServerApiUrl = "https://api-production.syro.com/v1/functions"

func GetServerUrl(env string, customServerUrl string) (serverApiUrl string) {
	if len(customServerUrl) > 0 {
		return fmt.Sprintf(`%s/v1/functions`, customServerUrl)
	}
	if env == "local" {
		return LocalServerApiUrl
	} else if env == "staging" {
		return StagingServerApiUrl
	}
	return ProductionServerApiUrl
}

// Auth API endpoints
var CliLogin = "/cli_login"
var CliValidateAccessTokenAndProjectId = "/cli_validate_access_token_and_project_id"
var CliValidateSessionToken = "/cli_validate_session_token"
var CliValidateMemberId = "/cli_validate_member_id"
var CliUserMemberships = "/cli_user_memberships"

// Project API endpoints
var CliValidateProjectId = "/cli_validate_project_id"
var CliProjectItems = "/cli_project_items"
var CliPullProjectItems = "/cli_pull_project_items"
