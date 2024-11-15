# 🔥 Ignite CLI Tool

Ignite is a command-line tool designed to help you quickly generate project structures and boilerplate code using predefined templates. With Ignite, users can initialize projects by running a single command, selecting specific templates, and customizing their project layout.🚀

## ✨ Features

- Generates project files and structure based on predefined templates.
- Supports dynamic configuration and project naming.
- Simple installation with Go install.
- Embeds templates directly in the binary for easy portability.

## 📥 Installation

To install Ignite, ensure you have Go installed on your system, and run the following command:

```bash
go install github.com/EmilioCliff/ignite@latest
```

This will install the ignite binary to your $GOPATH/bin directory, making it accessible from any directory in your terminal.

## ⚡ Usage

Once installed, you can use the ignite command from any directory. The basic usage is as follows:

There are two modes to run ignite in:

1.  **Interactive Mode**

In Interactive Mode, ignite guides the user step-by-step through a series of prompts to collect the necessary inputs for generating the project structure. This mode is user-friendly and ideal for first-time users or those who prefer not to memorize or use command-line flags.

**How to Use Interactive Mode:**
Simply run the ignite command without any flags:

```bash
ignite <project_name>
```

or

```bash
ignite <project_name> --interactive
```

**Steps:**

**Prompt 1:** Choose a database type: (postgres or http)
**Prompt 2:** Choose a controller type: (grpc or http)
**Prompt 3:** Do you want to include a GitHub Actions workflow? (yes/no):

#### Required Inputs

Project Name **(required)**: The name of the project to be created.
`--interactive` **(optional)**: Sets the mode to interactive when flag is passed interactive mode is set.
`--path` **(optional)**: Sets the path to create directory (defaults to current dir).
`--verbose` **(optional)**: logs the output to the terminal (defaults to false)

2. **Flag Mode**

In Flag Mode, ignite allows advanced users to specify all necessary inputs directly via command-line flags. This mode is faster for experienced users who are familiar with the tool and their desired configuration.

**How to Use Flag Mode:**
Run the ignite command with the appropriate flags:

```bash
ignite <project_name> -d <database> -c <controller> [other flags]
```

#### Required Flags/Inputs:

Project Name **(required)**: The name of the project to be created.
`--database` **(required)**: Specifies the database type (e.g., postgres, mysql).
`--controller` **(required)**: Specifies the controller type (e.g., http, grpc).
`--path` **(optional)**: Sets the path to create directory (defaults to current dir).
`--interactive` **(optional)**: Sets the mode to interactive when flag is passed interactive mode is set.
`--withDockerfile` **(optional)**: Sets if a dockerfile will also be generated (defaults to false)
`--withWorkflow` **(optional)**: Sets if a github workflow will also be generated (defaults to false)
`--verbose` **(optional)**: logs the output to the terminal (defaults to false)

No interactive prompts are shown.

> If required flags are missing, ignite will return an error with a list of missing inputs.

## 🛠️ Troubleshooting

If need help there is the `-h` or `--help` flag and will be guided

```bash
ignite --help
```

Example Output:

```bash
ignite initializes a new project.

Usage examples:

  ignite my_project
  ignite my_project --interactive
  ignite my_project -d postgres -c http -p ./path/to/project

Supported Database Types: postgres, mysql, sqlite, mongodb
Supported Controllers: user, auth, product, order

Usage:
  ignite <project_name> [flags]

Flags:
  -c, --controller string   Controller type (one of: grpc, http)
  -d, --database string     Database type (one of: postgres, mysql)
  -h, --help                help for ignite
      --interactive         Interactive mode
  -p, --path string         Path to create project (defaults to current directory)
  -v, --verbose             verbose output
      --withDockerfile      Include Dockerfile? (yes/no)
      --withWorkflow        Include GitHub Actions workflow? (yes/no)
```

## 🤝 Contribution

Contributions are welcome! 💡 Please feel free to submit a pull request or report issues.

1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Commit your changes and push your branch.
4. Open a pull request for review.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
