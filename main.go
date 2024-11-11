package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var (
		dbType         string
		controlType    string
		withWorkflow   bool
		withDockerfile bool
		path           string
		interactive    bool
		verbose        bool
	)

	var rootCmd = &cobra.Command{
		Use:   "ignite <project_name>",
		Short: "Initialize a new project with the specified name",
		Long: `ignite initializes a new project.

Usage examples:

  ignite my_project
  ignite my_project --interactive 
  ignite my_project -d postgres -c http -p ./path/to/project

Supported Database Types: postgres, mysql, sqlite, mongodb
Supported Controllers: user, auth, product, order`,

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing project name: the first argument must be the project name")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println(dbType, controlType, withWorkflow, withDockerfile, path, interactive, verbose)
			logFile, err := os.OpenFile(".logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening log file: %v\n", err)
				os.Exit(1)
			}
			defer logFile.Close()

			var logOutput io.Writer = logFile
			if verbose {
				logOutput = io.MultiWriter(os.Stdout, logFile)
			}
			log.SetOutput(logOutput)

			p := NewProjectInitializer(
				path,
				strings.ToLower(dbType),
				strings.ToLower(controlType),
				withWorkflow,
				withDockerfile,
				verbose,
			)

			// check if it will run in interactive or manual way
			if interactive || len(args) == 1 && dbType == "" {
				runInInteractiveMode(p)
			} else {
				runFlagMode(p)
			}

			if path == "" {
				path, err = os.Getwd()
				if err != nil {
					log.Panic(err)
				}
			}

			err = changeWorkingDir(path)
			if err != nil {
				log.Panic(err)
			}

			projectName := args[0]
			log.Println("Initializing project", projectName)

			if err := p.runSetup(); err != nil {
				log.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&path, "path", "p", "", "Path to create project (defaults to current directory)")
	rootCmd.Flags().StringVarP(&dbType, "database", "d", "", "Database type (one of: postgres, mysql)")
	rootCmd.Flags().StringVarP(&controlType, "controller", "c", "", "Controller type (one of: grpc, http)")
	rootCmd.Flags().BoolVar(&withWorkflow, "withWorkflow", false, "Include GitHub Actions workflow? (yes/no)")
	rootCmd.Flags().BoolVar(&withDockerfile, "withDockerfile", false, "Include Dockerfile? (yes/no)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "Interactive mode")

	rootCmd.MarkFlagsRequiredTogether("database", "controller")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
