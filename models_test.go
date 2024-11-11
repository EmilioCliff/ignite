package main

import (
	"testing"
)

func TestIsSupported(t *testing.T) {
	supported := []string{"mysql", "postgres"}
	if !isSupported(supported, "mysql") {
		t.Error("Expected 'mysql' to be supported")
	}
	if isSupported(supported, "sqlite") {
		t.Error("Expected 'sqlite' to be unsupported")
	}
}

func TestRunFlagModeUnsupportedDB(t *testing.T) {
	data := &ProjectInitializer{dbType: "unsupportedDB"}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for unsupported database type")
		}
	}()
	runFlagMode(data)
}
