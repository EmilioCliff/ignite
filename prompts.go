package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type PromptContent struct {
	label    string
	errorMsg string
	success  string
}

// promptGetInput will prompt the user for input using the Label field of the PromptContent object
// and return the user's response as a string. If the user enters an empty string, an error will
// be returned.
func (pc *PromptContent) promptGetInput() string {
	prompt := promptui.Prompt{
		Label: pc.label,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	return result
}

// promptSelect will prompt the user to select one of the given items from a list. The Label
// field of the PromptContent object will be used as the prompt label. The user will be
// repeatedly prompted until they select a valid item. The selected item will be returned as
// a string.
func (pc *PromptContent) promptSelect(items []string) string {
	index := -1

	var result string

	var err error

	templates := &promptui.SelectTemplates{
		Active:   "{{ . | green }}",
		Inactive: "{{ . | red }}",
		Selected: fmt.Sprintf("✔ {{ . | bold | green }}"),
	}

	for index < 0 {
		prompt := promptui.Select{
			Label:     pc.label,
			Items:     items,
			Templates: templates,
		}

		index, result, err = prompt.Run()
		if err != nil {
			return fmt.Sprintf("Prompt failed %v\n", err)
		}
	}

	return result
}
