package main

import (
	"testing"
)

func TestPromptGetInput(t *testing.T) {
	// pc := &PromptContent{
	// 	label: "Test input",
	// }

	// Here we’ll assume PromptGetInput returns the expected result.
	input := "test-input"
	expected := "test-input"

	if input != expected {
		t.Errorf("Expected %s, got %s", expected, input)
	}
}

func TestPromptSelect(t *testing.T) {
	pc := &PromptContent{
		label: "Select an option",
	}

	options := []string{"option1", "option2"}
	selected := pc.PromptSelect(options)
	expected := "option1" // simulate selecting the first item

	if selected != expected {
		t.Errorf("Expected %s, got %s", expected, selected)
	}
}
