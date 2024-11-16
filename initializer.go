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

type projectInitializer struct {
	path           string
	projectName    string
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

// NewProjectInitializer returns a new ProjectInitializer instance.
//
// The parameters are as follows:
//
//   - path: the path where the project will be initialized.
//   - dbType: the database type to be used (e.g. postgres, mysql).
//   - controlType: the type of controller to be generated (e.g. gRPC, HTTP).
//   - withWorkflow: if true, a GitHub Actions workflow will be generated.
//   - withDockerfile: if true, a Dockerfile will be generated.
//   - verbose: if true, the command will log more information to the console.
func NewProjectInitializer(path, dbType, controlType string, withWorkflow, withDockerfile, verbose bool) *projectInitializer {
	return &projectInitializer{
		path:           path,
		dbType:         dbType,
		controlType:    controlType,
		withWorkflow:   withWorkflow,
		withDockerfile: withDockerfile,
		verbose:        verbose,
	}
}

// runSetup initializes a new project given the configuration.
//
// It does the following steps:
//
//   - sets the current working directory.
//   - creates the project structure.
//   - initializes the project modules.
//
// If any step fails, it will return an error.
func (p *projectInitializer) runSetup() error {
	var err error

	p.path, err = MustGetPwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	if err := p.updateProjectStructure(); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	if err := p.initializeModules(); err != nil {
		return fmt.Errorf("failed to initialize project modules: %w", err)
	}

	fmt.Println("Project initialized successfully!")

	return nil
}

// updateProjectStructure creates the project structure.
//
// It takes the default project structure and then modifies it according to the
// configuration. If the database type is set, it adds the corresponding database
// directory to the internal directory. If the controller type is set to grpc, it
// adds the gapi directory with the generated and proto subdirectories. If the
// withWorkflow flag is set, it adds the .github directory with the workflows
// subdirectory. If the withDockerfile flag is set, it adds the Dockerfile.
//
// Finally, it creates the directories according to the modified project structure
// using the createDirectories function.
func (p *projectInitializer) updateProjectStructure() error {
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

	return createDirectories(projectStructure, p.path, templateData{
		DBType:     p.dbType,
		SqlPackage: p.dbType == "postgres",
	})
}

// initializeModules initializes the project's Go module and Git repository.
//
// It runs the following commands in the project directory:
//
//   - go mod init <project_name>
//   - git init
//
// If any of the commands fail, it returns an error.
func (p *projectInitializer) initializeModules() error {
	if err := changeWorkingDir(p.path); err != nil {
		return err
	}

	log.Println("Initializing go module...")

	err := runCommand("go", "mod", "init", p.projectName)
	if err != nil {
		return fmt.Errorf("failed to run go mod init: %w", err)
	}

	log.Println("Initializing git repository...")

	err = runCommand("git", "init")
	if err != nil {
		return fmt.Errorf("failed to run git init: %w", err)
	}

	return nil
}

// createDirectories creates a directory structure based on the given map recursively.
//
// The function logs information about the directories and files it creates, and
// returns an error if any step of the process fails.
func createDirectories(structure map[string]interface{}, basePath string, data templateData) error {
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

			if err := createDirectories(v, fullPath, data); err != nil {
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
					templatePath := fmt.Sprintf("templates/%s.txt", name[:len(name)-len(".yaml")])

					tmpl, err := template.ParseFS(templatesFS, templatePath)
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

// getDefaultProjectStructure returns the default project structure.
//
// It is a map where the keys are the top-level directories and the values are
// either a nested map or nil. If the value is nil, it means that the directory
// is empty and should be created as such. If the value is a map, it will be
// recursively traversed and the same rules will apply.
//
// The default project structure is as follows:
//
//   - .envs: contains .local and configs directories.
//   - .local: contains config.env file.
//   - configs: contains sqlc.yaml file.
//   - cmd: contains server and cli directories.
//   - server: contains main.go file.
//   - cli: contains main.go file.
//   - internal: contains handlers, repository, mock, and services directories.
//   - pkg: contains errors.go file.
//   - README.md
//   - .gitignore
//   - Makefile
func (p *projectInitializer) getDefaultProjectStructure() map[string]interface{} {
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
		"pkg": map[string]interface{}{
			"errors.go": "",
		},
		"README.md":  "",
		".gitignore": "",
		"Makefile":   "",
	}
}

// RunCommand runs a command with the given name and arguments, and returns an error
// if the command fails. It redirects the command's stdout and stderr to the
// corresponding writer in the OS.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func changeWorkingDir(path string) error {
	return os.Chdir(path)
}

func MustGetPwd() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return path, nil
}
