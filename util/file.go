/*
Copyright © 2023 Syro team <info@syro.com>
*/
package util

import (
	"encoding/json"
	"fmt"
	"os"
	"syro/model"
)

func LoadConfigFromProjectConfigFile() (isConfigLoaded bool, config model.Config, err error) {
	configFileDirectory := "./.syro"
	if err := os.MkdirAll(configFileDirectory, 0755); os.IsExist(err) {
	} else {
	}

	configFileFullPath := "./.syro/config.json"
	if _, error := os.Stat(configFileFullPath); os.IsNotExist(error) {
		return false, model.Config{}, nil
	} else {
		existingConfigBytes, err := os.ReadFile(configFileFullPath)
		if err != nil {
			fmt.Printf("Unable to open the file %s.\n", configFileFullPath)
			fmt.Printf("Error :: %v\n", err)
			return false, model.Config{}, err
		}

		existingConfig := model.Config{}
		err = json.Unmarshal(existingConfigBytes, &existingConfig)
		if err != nil {
			fmt.Printf("Unable to unmarshal contents from config file.\n")
			fmt.Printf("Error :: %v\n", err)
			return false, model.Config{}, err
		}
		return true, existingConfig, nil
	}
}

func SaveUserAndSessionInfoToProjectConfigFile(companyId string, expiresAt string, memberId string, sessionToken string) (err error) {
	_, config, err := LoadConfigFromProjectConfigFile()
	if err != nil {
		fmt.Printf("Unable to load items from config file.\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	config.UpdateUserAndSessionInfo(companyId, expiresAt, memberId, sessionToken)

	configBytes, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Unable to marshal contents for config file.\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	configFileFullPath := "./.syro/config.json"
	newFile, err := os.Create(configFileFullPath)
	if err != nil {
		fmt.Printf("Unable to create config file\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	_, err = newFile.Write(configBytes)
	if err != nil {
		fmt.Printf("Unable to write to config file\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	return nil
}

func SaveProjectIdToProjectConfigFile(projectId string) (err error) {
	_, config, err := LoadConfigFromProjectConfigFile()
	if err != nil {
		fmt.Print("Unable to load items from config file.\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	config.UpdateProjectId(projectId)

	configBytes, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Unable to marshal contents for config file.\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	configFileFullPath := "./.syro/config.json"
	newFile, err := os.Create(configFileFullPath)
	if err != nil {
		fmt.Printf("Unable to create config file\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	_, err = newFile.Write(configBytes)
	if err != nil {
		fmt.Printf("Unable to write contents to config file\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	return nil
}

func SaveCompanyIdAndValidatedInfoToProjectConfigFile(companyId string, validatedAccessToken string, validatedProjectId string) (err error) {
	config := model.Config{
		CompanyId:            companyId,
		ValidatedAccessToken: validatedAccessToken,
		ValidatedProjectId:   validatedProjectId,
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("Unable to marshal contents for config file.\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	configFileFullPath := "./.syro/config.json"
	newFile, err := os.Create(configFileFullPath)
	if err != nil {
		fmt.Printf("Unable to create config file\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	_, err = newFile.Write(configBytes)
	if err != nil {
		fmt.Printf("Unable to write to config file\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	return nil
}

func SaveSecretsToEnvFile(items []model.ItemDetails) (err error) {
	// TO DO: Uncomment once checking of existing env file is required
	/*
		isExistingEnvLoaded, existingEnvItems, err := LoadItemsFromExistingEnvFile()
		if err != nil {
			fmt.Printf("Unable to load items from env file.\n")
			return err
		}
	*/

	envItemsString := ""
	// TO DO: Uncomment once checking of existing env file is required
	/*
		if isExistingEnvLoaded {
			fmt.Printf("An env file already exists for this project.\n")
			isOverriding, err := isOverridingItemsInExistingEnvFile()
			if err != nil {
				fmt.Printf("Something went wrong with the isOverriding prompt.\n")
				return err
			}
			if !isOverriding {
				isMerging, err := isMergingItemsToExistingEnvFile()
				if err != nil {
					fmt.Printf("Something went wrong with the isMerging prompt.\n")
					return err
				}
				if isMerging {

					for index := 0; index < len(existingEnvItems); index++ {
						envItemsString += fmt.Sprintf("%s='%s'\n", existingEnvItems[index].Key, existingEnvItems[index].Value)
					}
				}
			}
		}
	*/

	for index := 0; index < len(items); index++ {
		envItemsString += fmt.Sprintf("%s='%s'\n", items[index].Key, items[index].Value)
	}

	envFileFullPath := ".env"
	newFile, err := os.Create(envFileFullPath)
	if err != nil {
		fmt.Printf("Unable to create env file.\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	_, err = newFile.WriteString(envItemsString)
	if err != nil {
		fmt.Printf("Unable to write to env file.\n")
		fmt.Printf("Error :: %v\n", err)
		return err
	}

	return nil
}

// TO DO: Uncomment once checking of existing env file is required
/*
func LoadItemsFromExistingEnvFile() (isEnvLoaded bool, envItems []api.ItemDetails, err error) {
	isEqualSign := func(char rune) bool {
		return char == '='
	}

	envFileFullPath := ".env"
	if _, error := os.Stat(envFileFullPath); os.IsNotExist(error) {
		return false, []api.ItemDetails{}, nil
	} else {
		existingEnvItemsString, err := ioutil.ReadFile(envFileFullPath)
		if err != nil {
			fmt.Printf("Unable to open the file %s.\n", envFileFullPath)
			return false, []api.ItemDetails{}, err
		}
		existingEnvItemsArray := strings.Fields(string(existingEnvItemsString))
		existingEnvItems := []api.ItemDetails{}
		for index := 0; index < len(existingEnvItemsArray); index++ {
			itemDetails := strings.FieldsFunc(existingEnvItemsArray[index], isEqualSign)
			existingEnvItems = append(existingEnvItems, api.ItemDetails{Key: itemDetails[0], Value: trimQuotes(itemDetails[1])})
		}
		return true, existingEnvItems, nil
	}

}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func isOverridingItemsInExistingEnvFile() (isOverriding bool, err error) {
	prompt := promptui.Select{
		Label: "Would you like to override it with secrets coming from Syro? Select[Yes/No]",
		Items: []string{"No", "Yes"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result == "Yes", nil
}

func isMergingItemsToExistingEnvFile() (isOverriding bool, err error) {
	prompt := promptui.Select{
		Label: "Would you like to merge it with secrets coming from Syro? Select[Yes/No]",
		Items: []string{"No", "Yes"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result == "Yes", nil
}
*/
