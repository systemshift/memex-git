package main

import (
	"os"

	"memex-git/git"
)

func main() {
	// Get repository path from environment or argument
	repoPath := os.Getenv("MEMEX_REPO_PATH")
	if repoPath == "" && len(os.Args) > 1 {
		repoPath = os.Args[1]
	}
	if repoPath == "" {
		repoPath = "." // Default to current directory
	}

	// Create module instance
	module := git.NewModule(nil, repoPath)

	// Handle plugin protocol
	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "id":
		os.Stdout.WriteString(module.ID())
	case "name":
		os.Stdout.WriteString(module.Name())
	case "describe":
		os.Stdout.WriteString(module.Description())
	case "run":
		if len(os.Args) < 3 {
			os.Stderr.WriteString("command required")
			os.Exit(1)
		}
		err := module.HandleCommand(os.Args[2], os.Args[3:])
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	default:
		os.Stderr.WriteString("unknown command")
		os.Exit(1)
	}
}
