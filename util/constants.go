/*
Copyright © 2023 Syro team <info@syro.com>
*/
package util

var LocalServerApplicationId = "paysail_local_app_id"
var StagingServerApplicationId = "paysail_api_staging_app_id"
var ProductionServerApplicationId = "paysail_api_production_app_id"

func GetServerApplicationId(env string) (serverApiUrl string) {
	if env == "local" {
		return LocalServerApplicationId
	} else if env == "staging" {
		return StagingServerApplicationId
	}
	return ProductionServerApplicationId
}
