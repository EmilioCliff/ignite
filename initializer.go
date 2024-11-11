package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"
)

var HOME = "/home/emilio/ignite"

type ProjectInitializer struct {
	path           string
	dbType         string
	controlType    string
	withWorkflow   bool
	withDockerfile bool
	verbose        bool
}
type templateData struct {
	DBType     string
	SqlPackage bool
}

func NewProjectInitializer(path, dbType, controlType string, withWorkflow, withDockerfile, verbose bool) *ProjectInitializer {
	return &ProjectInitializer{
		path:           path,
		dbType:         dbType,
		controlType:    controlType,
		withWorkflow:   withWorkflow,
		withDockerfile: withDockerfile,
		verbose:        verbose,
	}
}

func (p *ProjectInitializer) runSetup() error {
	var err error

	p.path, err = MustGetPwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if err := p.createProjectStructure(); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	if err := p.initializeModules(); err != nil {
		return fmt.Errorf("failed to initialize project modules: %w", err)
	}

	fmt.Println("Project initialized successfully!")

	return nil
}

func (p *ProjectInitializer) createProjectStructure() error {
	projectStructure := p.getDefaultProjectStructure()

	if p.dbType != "" {
		projectInternal, ok := projectStructure["internal"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to access internal directory in project structure")
		}

		projectInternal[p.dbType] = map[string]interface{}{
			"generated":  nil,
			"migrations": nil,
			"queries":    nil,
			"mock":       nil,
		}
	}

	if p.controlType == "grpc" {
		projectStructure["gapi"] = map[string]interface{}{
			"generated": nil,
			"proto":     nil,
		}
	}

	if p.withWorkflow {
		projectStructure[".github"] = map[string]interface{}{
			"workflows": map[string]interface{}{"ci.yml": ""},
		}
	}

	if p.withDockerfile {
		projectStructure["Dockerfile"] = ""
	}

	return CreateDirectories(projectStructure, p.path, templateData{
		DBType:     p.dbType,
		SqlPackage: p.dbType == "postgres",
	})
}

func (p *ProjectInitializer) initializeModules() error {
	if err := changeWorkingDir(p.path); err != nil {
		return err
	}

	log.Println("Initializing go module...")

	err := RunCommand("go", "mod", "init", filepath.Base(p.path))
	if err != nil {
		return fmt.Errorf("failed to run go mod init: %w", err)
	}

	log.Println("Initializing git repository...")

	err = RunCommand("git", "init")
	if err != nil {
		return fmt.Errorf("failed to run git init: %w", err)
	}

	return nil
}

func changeWorkingDir(path string) error {
	return os.Chdir(path)
}

func (p *ProjectInitializer) getDefaultProjectStructure() map[string]interface{} {
	return map[string]interface{}{
		".envs": map[string]interface{}{
			".local":  map[string]interface{}{"config.env": ""},
			"configs": map[string]interface{}{"sqlc.yaml": ""},
		},
		"cmd": map[string]interface{}{
			"server": map[string]interface{}{"main.go": ""},
			"cli":    map[string]interface{}{"main.go": ""},
		},
		"internal": map[string]interface{}{
			"handlers":   nil,
			"repository": nil,
			"mock":       nil,
			"services":   nil,
		},
		"pkg":        nil,
		"README.md":  "",
		".gitignore": "",
		"Makefile":   "",
	}
}

func CreateDirectories(structure map[string]interface{}, basePath string, data templateData) error {
	var mutex sync.Mutex

	for name, item := range structure {
		fullPath := filepath.Join(basePath, name)

		switch v := item.(type) {
		case map[string]interface{}:
			log.Println("Creating directory:", fullPath)

			mutex.Lock()
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				mutex.Unlock()

				return fmt.Errorf("failed to create directory %s: %v", fullPath, err)
			}
			mutex.Unlock()

			if err := CreateDirectories(v, fullPath, data); err != nil {
				return err
			}
		case nil:
			log.Println("Creating empty directory:", fullPath)

			mutex.Lock()
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				mutex.Unlock()

				return fmt.Errorf("failed to create directory %s: %v", fullPath, err)
			}
			mutex.Unlock()
		case string:
			log.Println("Creating file:", fullPath)
			mutex.Lock()

			file, err := os.Create(fullPath)
			if err != nil {
				mutex.Unlock()

				return fmt.Errorf("failed to create file %s: %v", fullPath, err)
			}

			if content, exists := templates[name]; exists {
				if content == "" {
					templatePath := fmt.Sprintf("%s/templates/%s", HOME, name[:len(name)-len(".yaml")]+".txt")

					tmpl, err := template.ParseFiles(templatePath)
					if err != nil {
						file.Close()
						mutex.Unlock()

						return fmt.Errorf("failed to parse template %s: %v", templatePath, err)
					}

					if err := tmpl.Execute(file, data); err != nil {
						file.Close()
						mutex.Unlock()

						return fmt.Errorf("failed to execute template for %s: %v", templatePath, err)
					}
				} else {
					if _, err := file.WriteString(content); err != nil {
						file.Close()
						mutex.Unlock()

						return fmt.Errorf("failed to write to file %s: %v", fullPath, err)
					}
				}
			}

			file.Close()
			mutex.Unlock()
		default:
			return fmt.Errorf("invalid item type for %s", fullPath)
		}
	}

	return nil
}

func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func MustGetPwd() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return path, nil
}

// type ProjectInitializer struct {
// 	db         Database
// 	controller Controller
// }

// func NewProjectInitializer(db Database, controller Controller) *ProjectInitializer {
// 	return &ProjectInitializer{db: db, controller: controller}
// }

// func (p *ProjectInitializer) RunSetup() error {
// 	if err := p.db.Setup(); err != nil {
// 		return fmt.Errorf("database setup failed: %w", err)
// 	}
// 	if err := p.controller.Setup(); err != nil {
// 		return fmt.Errorf("controller setup failed: %w", err)
// 	}
// 	fmt.Println("Project initialized successfully!")
// 	return nil
// }

// func MustGetPwd() {}

// here

// var (
// 	currentPath string // To track the working directory path
// 	mu          sync.Mutex
// )

// var templates = map[string]string{
// 	".gitignore": `
// # Binaries
// bin/
// *.exe
// *.dll
// *.so
// *.dylib

// # Configs
// *.env

// # Logs
// *.log
// `,

// 	"Dockerfile": `
// FROM golang:1.23-alpine3.20 AS builder
// WORKDIR /app
// COPY . .
// RUN go build -o main /app/cmd/server/main.go

// EXPOSE 3030

// CMD ["./main"]
// `,

// 	"main.go": `
// package main

// import "fmt"

// func main() {
// 	fmt.Println("Hello World!")
// }
// `,

// 	"sqlc.yaml": "",

// 	// 	"sqlc.yaml": `
// 	// version: "1"
// 	// packages:
// 	//   - path: "internal/mysql/queries"
// 	//     queries: "./queries/"
// 	//     schema: "./migrations/"
// 	//     engine: "mysql"
// 	// `,

// 	"ci.yml": `
// name: ci-test

// on:
// 	push:
// 	branches: [main]
// 	pull_request:
// 	branches: [main]

// jobs:
// 	test:
// 	name: Test
// 	runs-on: ubuntu-latest
// 	steps:
// 		- uses: actions/checkout@v2

// 		- name: Set up Go
// 		uses: actions/setup-go@v4
// 		with:
// 			go-version: "^1.21"

// 		- name: Run tests
// 		run: make race-test
// `,

// 	"Makefile": `
// test:
// 	go test -v ./...

// race-test:
// 	go test -v -race ./...

// sqlc:
// 	cd .envs/configs && sqlc generate

// run:
// 	cd cmd/server && go run main.go

// .PHONY: test race-test sqlc run
// 	`,

// 	"README.md": `
// 	# Ignite Project

// 	This project was created with [Ignite](https://github.com/emilio/ignite) — a CLI tool for bootstrapping Go-based applications with flexibility for various configurations.

// 	## Table of Contents

// 	- [Getting Started](#getting-started)
// 	- [Available Commands](#available-commands)
// 	- [Project Structure](#project-structure)
// 	- [Configuration](#configuration)
// 	- [Collaboration](#collaboration)

// 	## Getting Started

// 	To start working with this project, clone the repository, navigate into the project directory, and run the following command to install dependencies:

// 	` + "```" + `sh
// 	go mod tidy

// 	# Make sure you have **Go** and **Git** installed on your system.

// 	# Running the Project

// 	After setting up the project, you can use the following command to start the server:
// 	` + "```" + `sh

// 	## Available Commands

// 	In the project directory, you can run:

// 	` + "```" + `sh
// 		make sqlc
// 		make test
// 		make race-test
// 		` + "```" + `

// 	## Project Structure

// 	Ignite sets up a flexible folder structure based on hexagonal architecture and repository pattern:

// 	` + "```" + `sh
// 		.envs                    # Environment configurations
// 		cmd
// 		├── server               # Server main entry point
// 		└── cli                  # CLI main entry point (if CLI option selected)
// 		gapi                     # gRPC generated files (if gRPC selected)
// 		internal
// 		├── handlers             # HTTP handler functions
// 		├── gapi                 # gRPC service implementations
// 		├── repository           # Data access layer
// 		├── services             # Business logic layer
// 		└── mysql/postgres       # Database-related files (queries, migrations, mocks)
// 		pkg                      # Common utilities and helpers
// 		.github/workflows        # CI configuration (if --withWorkflow selected)

// 		` + "```" + `

// 	## Configuration

// 	Project configurations are set in environment variables and configuration files:

// 	` + "`.envs/.local/config.env`" + ` - for local environment configurations
// 	` + "`.envs/configs/sqlc.yaml`" + ` - SQLC configuration for SQL code generation

// 	Adjust these files as needed for different environments.

// 	## Collaboration

// 	We welcome contributions! If you want to add new features, improve the documentation, or fix bugs, please follow these steps:

// 	1. **Fork the repository**: Create a personal copy of the repository to work on.
// 	2. **Create a new branch**: Develop your changes in a separate branch. For example, ` + "`feature/new-feature` or" + "`bugfix/fix-issue`" + `.
// 	3. **Commit your changes**: Make sure to write meaningful commit messages describing what your changes do.
// 	4. **Create a pull request**: Once your changes are ready, open a pull request to merge your branch into the main repository.

// 	### Features You Can Help Add:

// 	- **New Commands**: If you'd like to add new subcommands to the CLI tool, feel free to submit an enhancement.
// 	- **Database Integrations**: We currently support SQL-based databases like PostgreSQL and MySQL. Contributions for other databases are welcome!
// 	- **Testing**: Help us write more tests for different use cases and improve test coverage.
// 	- **CI/CD Workflows**: If you have experience with CI/CD tools, improving the ` + "`GitHub Actions`" + ` workflow for continuous integration is a great way to contribute.

// 	If you have an idea for a new feature or improvement, please open an issue or start a discussion. We'd love to hear your thoughts and collaborate!

// `,
// }

// // structure defining the folder setup
// // var projectStructure = map[string]interface{}{
// // 	".envs": map[string]interface{}{
// // 		".local":  map[string]interface{}{"config.env": ""},
// // 		"configs": map[string]interface{}{"sqlc.yaml": ""},
// // 	},
// // 	"cmd": map[string]interface{}{
// // 		"server": map[string]interface{}{"main.go": ""},
// // 		"cli":    map[string]interface{}{"main.go": ""},
// // 	},
// // 	"gapi": map[string]interface{}{
// // 		"generated": nil,
// // 		"proto":     nil,
// // 	},
// // 	"internal": map[string]interface{}{
// // 		"handlers":   nil,
// // 		"gapi":       map[string]interface{}{"generated": nil},
// // 		"repository": nil,
// // 		"mock":       nil,
// // 		"services":   nil,
// // 		"mysql": map[string]interface{}{
// // 			"generated":  nil,
// // 			"migrations": nil,
// // 			"queries":    nil,
// // 			"mock":       nil,
// // 		},
// // 	},
// // 	"pkg": nil,
// // 	".github": map[string]interface{}{
// // 		"workflows": map[string]interface{}{"ci.yml": ""},
// // 	},
// // 	"Dockerfile": "",
// // 	"README.md":  "",
// // 	".gitignore": "",
// // }

// var projectStructure = map[string]interface{}{
// 	".envs": map[string]interface{}{
// 		".local":  map[string]interface{}{"config.env": ""},
// 		"configs": map[string]interface{}{"sqlc.yaml": ""},
// 	},
// 	"cmd": map[string]interface{}{
// 		"server": map[string]interface{}{"main.go": ""},
// 		"cli":    map[string]interface{}{"main.go": ""},
// 	},
// 	"internal": map[string]interface{}{
// 		"handlers":   nil,
// 		"repository": nil,
// 		"mock":       nil,
// 		"services":   nil,
// 	},
// 	"pkg":        nil,
// 	"README.md":  "",
// 	".gitignore": "",
// 	"Makefile":   "",
// 	// ".github": map[string]interface{}{
// 	// 	"workflows": map[string]interface{}{"ci.yml": ""},
// 	// },
// 	// "Dockerfile": "",
// }

// type templateData struct {
// 	DBType     string
// 	SqlPackage bool
// }

// func CreateDirectories(structure map[string]interface{}, basePath string, data templateData) error {
// 	for name, item := range structure {
// 		fullPath := filepath.Join(basePath, name)

// 		switch v := item.(type) {
// 		case map[string]interface{}:
// 			// Lock only around directory creation
// 			mu.Lock()
// 			log.Println("Creating directory:", fullPath)

// 			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
// 				mu.Unlock()
// 				return fmt.Errorf("failed to create directory %s: %v", fullPath, err)
// 			}
// 			mu.Unlock()

// 			// Recursive call without lock to avoid deadlock
// 			if err := CreateDirectories(v, fullPath, data); err != nil {
// 				return err
// 			}
// 		case nil:
// 			log.Println("Creating empty directory:", fullPath)
// 			// Create empty directories with lock
// 			mu.Lock()
// 			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
// 				mu.Unlock()
// 				return fmt.Errorf("failed to create directory %s: %v", fullPath, err)
// 			}
// 			mu.Unlock()
// 		case string:
// 			log.Println("Creating file:", fullPath)
// 			mu.Lock()
// 			file, err := os.Create(fullPath)
// 			if err != nil {
// 				mu.Unlock()
// 				return fmt.Errorf("failed to create file %s: %v", fullPath, err)
// 			}

// 			// Write template content if it exists
// 			if content, exists := templates[name]; exists {
// 				if content == "" {
// 					templatePath := fmt.Sprintf("templates/%s", name[:len(name)-len(".yaml")]+".txt")
// 					tmpl, err := template.ParseFiles(templatePath)
// 					if err != nil {
// 						file.Close()
// 						mu.Unlock()
// 						return fmt.Errorf("failed to parse template %s: %v", templatePath, err)
// 					}

// 					if err := tmpl.Execute(file, data); err != nil {
// 						file.Close()
// 						mu.Unlock()
// 						return fmt.Errorf("failed to execute template for %s: %v", templatePath, err)
// 					}
// 				} else {
// 					if _, err := file.WriteString(content); err != nil {
// 						file.Close()
// 						mu.Unlock()
// 						return fmt.Errorf("failed to write to file %s: %v", fullPath, err)
// 					}
// 				}
// 			}

// 			file.Close()
// 			mu.Unlock()
// 		default:
// 			return fmt.Errorf("invalid item type for %s", fullPath)
// 		}
// 	}

// 	return nil
// }

// func checkCommandAvailable(cmd string) bool {
// 	_, err := exec.LookPath(cmd)
// 	return err == nil
// }

// func InitializeModules(projectName, path string) error {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if err := os.Chdir(path); err != nil {
// 		return fmt.Errorf("failed to change to project directory: %v", err)
// 	}

// 	// Initialize go module
// 	if checkCommandAvailable("go") {
// 		log.Println("Initializing go module...")
// 		cmd := exec.Command("go", "mod", "init", projectName)
// 		if err := cmd.Run(); err != nil {
// 			return fmt.Errorf("failed to run go mod init: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Warning: `go` is not available, skipping go mod init.")
// 	}

// 	// Initialize git repository
// 	if checkCommandAvailable("git") {
// 		log.Println("Initializing git repository...")
// 		cmd := exec.Command("git", "init")
// 		if err := cmd.Run(); err != nil {
// 			return fmt.Errorf("failed to run git init: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Warning: `git` is not available, skipping git init.")
// 	}

// 	return nil
// }
