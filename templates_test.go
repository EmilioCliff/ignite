package main

import (
	"testing"
)

func TestTemplatesExist(t *testing.T) {
	expectedKeys := []string{".gitignore", "Dockerfile", "main.go", "sqlc.yaml", "ci.yml", "Makefile", "README.md"}
	for _, key := range expectedKeys {
		if _, exists := templates[key]; !exists {
			t.Errorf("Expected template %s is missing", key)
		}
	}
}
