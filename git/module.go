package git

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/systemshift/memex/pkg/sdk/types"
)

// Module implements Git repository management within memex
type Module struct {
	repo     types.Repository
	git      *git.Repository
	repoPath string // Path to the git repository
}

// NewModule creates a new git module
func NewModule(repo types.Repository, path string) *Module {
	// Ensure absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		// Fall back to provided path if abs fails
		absPath = path
	}

	return &Module{
		repo:     repo,
		repoPath: absPath,
	}
}

// SetRepository sets the repository instance
func (m *Module) SetRepository(repo types.Repository) {
	m.repo = repo
}

// ID returns module identifier
func (m *Module) ID() string {
	return "git"
}

// Name returns human-readable name
func (m *Module) Name() string {
	return "Git Management"
}

// Description returns module description
func (m *Module) Description() string {
	return "Manages Git version control within memex repositories"
}

// Commands returns available module commands
func (m *Module) Commands() []types.ModuleCommand {
	return []types.ModuleCommand{
		{
			Name:        "init",
			Description: "Initialize Git repository in memex directory",
			Usage:       "git init",
		},
		{
			Name:        "status",
			Description: "Show working directory status",
			Usage:       "git status",
		},
		{
			Name:        "add",
			Description: "Stage changes for commit",
			Usage:       "git add [path]",
		},
		{
			Name:        "commit",
			Description: "Commit staged changes",
			Usage:       "git commit <message>",
		},
		{
			Name:        "log",
			Description: "Show commit history",
			Usage:       "git log",
		},
	}
}

// ValidateNodeType validates node types (not used for this module)
func (m *Module) ValidateNodeType(nodeType string) bool {
	return false
}

// ValidateLinkType validates link types (not used for this module)
func (m *Module) ValidateLinkType(linkType string) bool {
	return false
}

// ValidateMetadata validates module-specific metadata
func (m *Module) ValidateMetadata(meta map[string]interface{}) error {
	return nil
}

// HandleCommand handles module commands
func (m *Module) HandleCommand(cmd string, args []string) error {
	switch cmd {
	case "init":
		return m.InitRepo()

	case "status":
		return m.ShowStatus()

	case "add":
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		return m.AddFiles(path)

	case "commit":
		if len(args) < 1 {
			return fmt.Errorf("commit message required")
		}
		return m.Commit(args[0])

	case "log":
		return m.ShowLog()

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

// InitRepo initializes a new Git repository
func (m *Module) InitRepo() error {
	// Initialize new git repository
	gitRepo, err := git.PlainInit(m.repoPath, false)
	if err != nil {
		return fmt.Errorf("initializing git repo: %w", err)
	}

	m.git = gitRepo
	return nil
}

// ShowStatus shows current Git status
func (m *Module) ShowStatus() error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	worktree, err := m.git.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("getting status: %w", err)
	}

	fmt.Println(status)
	return nil
}

// AddFiles stages files for commit
func (m *Module) AddFiles(path string) error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	worktree, err := m.git.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	// Resolve path relative to repository root
	absPath := filepath.Join(m.repoPath, path)

	// Add files to staging
	err = worktree.AddWithOptions(&git.AddOptions{
		All:  true,
		Path: absPath,
	})
	if err != nil {
		return fmt.Errorf("adding files: %w", err)
	}

	return nil
}

// Commit commits staged changes
func (m *Module) Commit(message string) error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	worktree, err := m.git.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	// Create commit
	_, err = worktree.Commit(message, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("committing changes: %w", err)
	}

	return nil
}

// ShowLog shows commit history
func (m *Module) ShowLog() error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	// Get commit iterator
	iter, err := m.git.Log(&git.LogOptions{})
	if err != nil {
		return fmt.Errorf("getting log: %w", err)
	}

	// Print commits
	err = iter.ForEach(func(commit *object.Commit) error {
		fmt.Printf("%s %s\n", commit.Hash.String(), commit.Message)
		return nil
	})
	if err != nil {
		return fmt.Errorf("iterating commits: %w", err)
	}

	return nil
}
