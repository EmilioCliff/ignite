package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var supportedDBTypes = []string{"postgres", "mysql"}

var supportedControllers = []string{"grpc", "http"}

func runInInteractiveMode(data *ProjectInitializer) {
	log.Println("Running in interactive mode.")

	dbPrompt := PromptContent{
		label:    "Choose a database type",
		errorMsg: "please provide a database type",
		success:  "Database: ",
	}

	data.dbType = dbPrompt.PromptSelect(supportedDBTypes)

	controllerPrompt := PromptContent{
		label:    "Choose a controller type",
		errorMsg: "please provide a controller type",
		success:  "Controller: ",
	}

	data.controlType = controllerPrompt.PromptSelect(supportedControllers)

	withWorkflowPrompt := PromptContent{
		label:    "Do you want to include a GitHub Actions workflow? (yes/no)",
		errorMsg: "please answer yes or no",
		success:  "Workflow: ",
	}

	if rst := withWorkflowPrompt.PromptGetInput(); rst == "yes" || rst == "y" {
		data.withWorkflow = true
	} else {
		data.withWorkflow = false
	}

	withDockerfilePrompt := PromptContent{
		label:    "Do you want to include a Dockerfile? (yes/no)",
		errorMsg: "please answer yes or no",
		success:  "Workflow: ",
	}

	if rst := withDockerfilePrompt.PromptGetInput(); rst == "yes" || rst == "y" {
		data.withDockerfile = true
	} else {
		data.withDockerfile = false
	}
}

func runFlagMode(data *ProjectInitializer) {
	if data.dbType != "" && !isSupported(supportedDBTypes, data.dbType) {
		fmt.Printf("Error: Unsupported database type '%s'. Supported types are: (%v)\n", data.dbType, strings.Join(supportedDBTypes, ", "))
		os.Exit(1)
	}

	if data.controlType != "" && !isSupported(supportedControllers, data.controlType) {
		fmt.Printf("Error: Unsupported controller type '%s'. Supported types are: (%v)\n", data.controlType, strings.Join(supportedControllers, ", "))
		os.Exit(1)
	}
}

func isSupported(supportedList []string, item string) bool {
	for _, supported := range supportedList {
		if supported == item {
			return true
		}
	}

	return false
}
