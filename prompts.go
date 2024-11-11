package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

var _ Prompts = (*PromptContent)(nil)

type Prompts interface {
	PromptGetInput() string
	PromptSelect(items []string) string
}

type PromptContent struct {
	label    string
	errorMsg string
	success  string
}

func (pc *PromptContent) PromptGetInput() string {
	prompt := promptui.Prompt{
		Label: pc.label,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	return result
}

func (pc *PromptContent) PromptSelect(items []string) string {
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
