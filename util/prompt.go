/*
Copyright © 2023 Syro team <info@syro.com>
*/
package util

import (
	"errors"

	"github.com/manifoldco/promptui"
)

func GetProjectId() (projectId string, err error) {
	projectIdValidator := func(input string) error {
		if len(input) < 1 {
			return errors.New("Invalid project ID.")
		}
		return nil
	}

	passwordPrompt := promptui.Prompt{
		Label:    "Project ID",
		Validate: projectIdValidator,
	}

	inputProjectId, err := passwordPrompt.Run()

	if err != nil {
		return "", err
	}

	return inputProjectId, nil
}

func AskForUserEmail() (email string, err error) {
	emailValidator := func(input string) error {
		isValid := IsEmailAddressValid(input)
		if isValid {
			return nil
		}
		return errors.New("Invalid email address.")
	}

	emailPrompt := promptui.Prompt{
		Label:    "Email",
		Validate: emailValidator,
	}

	inputEmail, err := emailPrompt.Run()

	if err != nil {
		return "", err
	}
	return inputEmail, nil
}

func AskForUserPassword() (password string, err error) {
	passwordValidator := func(input string) error {
		if len(input) < 1 {
			return errors.New("Invalid password.")
		}
		return nil
	}

	passwordPrompt := promptui.Prompt{
		Label:    "Password",
		Validate: passwordValidator,
		Mask:     '*',
	}

	inputPassword, err := passwordPrompt.Run()

	if err != nil {
		return "", err
	}

	return inputPassword, nil
}
