package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var supportedDBTypes = []string{"postgres", "mysql"}

var supportedControllers = []string{"grpc", "http"}

// runInInteractiveMode prompts the user for input to set the database type, controller type, inclusion of a GitHub Actions workflow, and inclusion of a Dockerfile.
func runInInteractiveMode(data *projectInitializer) {
	log.Println("Running in interactive mode.")

	dbPrompt := PromptContent{
		label:    "Choose a database type",
		errorMsg: "please provide a database type",
		success:  "Database: ",
	}

	data.dbType = dbPrompt.promptSelect(supportedDBTypes)

	controllerPrompt := PromptContent{
		label:    "Choose a controller type",
		errorMsg: "please provide a controller type",
		success:  "Controller: ",
	}

	data.controlType = controllerPrompt.promptSelect(supportedControllers)

	withWorkflowPrompt := PromptContent{
		label:    "Do you want to include a GitHub Actions workflow? (yes/no)",
		errorMsg: "please answer yes or no",
		success:  "Workflow: ",
	}

	if rst := withWorkflowPrompt.promptGetInput(); rst == "yes" || rst == "y" {
		data.withWorkflow = true
	} else {
		data.withWorkflow = false
	}

	withDockerfilePrompt := PromptContent{
		label:    "Do you want to include a Dockerfile? (yes/no)",
		errorMsg: "please answer yes or no",
		success:  "Workflow: ",
	}

	if rst := withDockerfilePrompt.promptGetInput(); rst == "yes" || rst == "y" {
		data.withDockerfile = true
	} else {
		data.withDockerfile = false
	}
}

// runFlagMode validates the inputs provided by the user in flag mode and if they are
// invalid, it prints an error message and exits the program.
func runFlagMode(data *projectInitializer) {
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
