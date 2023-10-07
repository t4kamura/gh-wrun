package interactive

import (
	"github.com/AlecAivazis/survey/v2"
)

func AskChoices(message string, choices []string, defaultInput string) (string, error) {
	var answer string

	prompt := &survey.Select{
		Message: message,
		Options: choices,
		Default: defaultInput,
	}
	err := survey.AskOne(prompt, &answer)
	return answer, err
}

func AskInput(message string, defaultInput string) (string, error) {
	var answer string

	prompt := &survey.Input{
		Message: message,
		Default: defaultInput,
	}
	err := survey.AskOne(prompt, &answer)
	return answer, err
}

func AskConfirm(message string, defaultInput bool) (bool, error) {
	var answer bool

	prompt := &survey.Confirm{
		Message: message,
		Default: defaultInput,
	}
	err := survey.AskOne(prompt, &answer)
	return answer, err
}
