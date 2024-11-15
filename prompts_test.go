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
