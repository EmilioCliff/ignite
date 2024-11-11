package main

import (
	"os"
	"testing"
)

func TestNewProjectInitializer(t *testing.T) {
	p := NewProjectInitializer("/path/to/project", "postgres", "grpc", true, false, true)
	if p.path != "/path/to/project" || p.dbType != "postgres" || p.controlType != "grpc" {
		t.Error("NewProjectInitializer did not initialize correctly")
	}
}

func TestChangeWorkingDir(t *testing.T) {
	currentDir := "/tmp"
	err := changeWorkingDir(currentDir)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	newDir, _ := os.Getwd()
	if newDir != currentDir {
		t.Errorf("Expected directory to be %s, got %s", currentDir, newDir)
	}
}
