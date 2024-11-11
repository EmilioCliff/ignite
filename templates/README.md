# Ignite Project

This project was created with [Ignite](https://github.com/emilio/ignite) — a CLI tool for bootstrapping Go-based applications with flexibility for various configurations.

## Table of Contents

- [Getting Started](#getting-started)
- [Available Commands](#available-commands)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Collaboration](#collaboration)

## Getting Started

To start working with this project, clone the repository, navigate into the project directory, and run the following command to install dependencies:

```sh
go mod tidy

# Make sure you have **Go** and **Git** installed on your system.

# Running the Project

After setting up the project, you can use the following command to start the server:
```

## Available Commands

In the project directory, you can run:

```sh
    make sqlc
    make test
    make race-test
```

## Project Structure

Ignite sets up a flexible folder structure based on hexagonal architecture and repository pattern:

```sh
    .envs                    # Environment configurations
    cmd
    ├── server               # Server main entry point
    └── cli                  # CLI main entry point (if CLI option selected)
    gapi                     # gRPC generated files (if gRPC selected)
    internal
    ├── handlers             # HTTP handler functions
    ├── gapi                 # gRPC service implementations
    ├── repository           # Data access layer
    ├── services             # Business logic layer
    └── mysql/postgres       # Database-related files (queries, migrations, mocks)
    pkg                      # Common utilities and helpers
    .github/workflows        # CI configuration (if --withWorkflow selected)

```

## Configuration

Project configurations are set in environment variables and configuration files:

`.envs/.local/config.env` - for local environment configurations
`.envs/configs/sqlc.yaml` - SQLC configuration for SQL code generation

Adjust these files as needed for different environments.

## Collaboration

We welcome contributions! If you want to add new features, improve the documentation, or fix bugs, please follow these steps:

1. **Fork the repository**: Create a personal copy of the repository to work on.
2. **Create a new branch**: Develop your changes in a separate branch. For example, `feature/new-feature` or `bugfix/fix-issue`.
3. **Commit your changes**: Make sure to write meaningful commit messages describing what your changes do.
4. **Create a pull request**: Once your changes are ready, open a pull request to merge your branch into the main repository.

### Features You Can Help Add:

- **New Commands**: If you'd like to add new subcommands to the CLI tool, feel free to submit an enhancement.
- **Database Integrations**: We currently support SQL-based databases like PostgreSQL and MySQL. Contributions for other databases are welcome!
- **Testing**: Help us write more tests for different use cases and improve test coverage.
- **CI/CD Workflows**: If you have experience with CI/CD tools, improving the `GitHub Actions` workflow for continuous integration is a great way to contribute.

If you have an idea for a new feature or improvement, please open an issue or start a discussion. We'd love to hear your thoughts and collaborate!
