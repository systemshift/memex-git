# Memex Git Module

A memex module that provides Git version control capabilities within memex repositories. This module allows you to maintain version control of your memex content independently of any external git repositories.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/systemshift/memex-git
```

2. Build the module:
```bash
cd memex-git
go build ./...
```

3. Install the module in memex:
```bash
memex module install /path/to/memex-git
```

## Usage

The git module provides several commands to manage version control within your memex repository:

### Initialize a Git Repository

Initialize a new git repository in your memex directory:

```bash
memex git init
```

### View Repository Status

Check the status of your working directory:

```bash
memex git status
```

### Stage Changes

Add files to the staging area:

```bash
memex git add .           # Add all changes
memex git add <path>      # Add specific file or directory
```

### Commit Changes

Commit staged changes:

```bash
memex git commit "Your commit message"
```

### View History

Show commit history:

```bash
memex git log
```

## Development

### Project Structure

- `git/module.go`: Core module implementation
- `plugin/main.go`: Binary plugin interface
- `go.mod`: Module dependencies

### Building from Source

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```
3. Build:
```bash
go build ./...
```

### Running Tests

```bash
go test ./...
```

## License

BSD 3-Clause License - see LICENSE file for details
