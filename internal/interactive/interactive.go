package interactive

import (
	"strconv"

	"github.com/manifoldco/promptui"
)

func AskChoices(message string, choices []string, defaultInput string) (string, error) {
	defaultCursor := 0
	for i, choice := range choices {
		if choice == defaultInput {
			defaultCursor = i
		}
	}
	prompt := promptui.Select{
		Label:     message,
		Items:     choices,
		CursorPos: defaultCursor,
		HideHelp:  true,
		Size:      10,
	}
	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func AskInput(message string, defaultInput string) (string, error) {
	prompt := promptui.Prompt{
		Label:   message,
		Default: defaultInput,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func AskConfirm(message string, defaultInput bool) (bool, error) {
	choices := []bool{true, false}
	defaultCursor := 0
	for i, choice := range choices {
		if choice == defaultInput {
			defaultCursor = i
		}
	}
	prompt := promptui.Select{
		Label:     message,
		Items:     choices,
		CursorPos: defaultCursor,
		HideHelp:  true,
	}
	_, result, err := prompt.Run()

	if err != nil {
		return false, err
	}
	resultBool, _ := strconv.ParseBool(result)

	return resultBool, nil
}
