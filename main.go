package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/systemshift/memex-git/git"
)

// ModuleInfo represents basic module information
type ModuleInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Get repository path from environment or use current directory
	repoPath := os.Getenv("MEMEX_REPO_PATH")
	if repoPath == "" {
		repoPath = "." // Default to current directory
	}

	// Create git module instance
	gitModule := git.NewModule(nil, repoPath)

	switch command {
	case "--info":
		info := ModuleInfo{
			ID:          "git",
			Name:        "Git Management",
			Description: "Manages Git version control within memex repositories",
			Version:     "1.0.0",
		}
		data, _ := json.Marshal(info)
		fmt.Println(string(data))

	case "init":
		if err := gitModule.InitRepo(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Initialized empty Git repository")

	case "status":
		if err := gitModule.ShowStatus(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "add":
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		if err := gitModule.AddFiles(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added files: %s\n", path)

	case "commit":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Error: commit message required\n")
			fmt.Fprintf(os.Stderr, "Usage: memex-git commit <message>\n")
			os.Exit(1)
		}
		message := args[0]
		if err := gitModule.Commit(message); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Committed changes: %s\n", message)

	case "log":
		if err := gitModule.ShowLog(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "help", "--help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("memex-git - Git version control integration for memex")
	fmt.Println("")
	fmt.Println("Usage: memex-git <command> [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  init                      Initialize a new Git repository")
	fmt.Println("  status                    Show working tree status")
	fmt.Println("  add [path]                Add files to index (default: current directory)")
	fmt.Println("  commit <message>          Commit staged changes with message")
	fmt.Println("  log                       Show commit history")
	fmt.Println("  help, --help              Show this help message")
	fmt.Println("  --info                    Show module metadata (JSON)")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  MEMEX_REPO_PATH           Path to the repository (default: current directory)")
}