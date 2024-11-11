package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"
)

var templatesFS embed.FS

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
					templatePath := fmt.Sprintf("templates/%s.txt", name[:len(name)-len(".yaml")])
					// templatePath, err := MustGetExecutablePath(name)
					// if err != nil {
					// 	file.Close()
					// 	mutex.Unlock()

					// 	return fmt.Errorf("failed to get template path: %v", err)
					// }
					// templatePath := fmt.Sprintf("%s/templates/%s", HOME, name[:len(name)-len(".yaml")]+".txt")

					tmpl, err := template.ParseFS(templatesFS, templatePath)
					if err != nil {
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

func MustGetExecutablePath(name string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not get executable path: %v", err)
	}

	baseDir := filepath.Dir(execPath)

	templatePath := filepath.Join(baseDir, "templates", name[:len(name)-len(".yaml")]+".txt")

	return templatePath, nil
}
