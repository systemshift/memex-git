package git

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// StandaloneModule implements Git repository management without memex dependencies
type StandaloneModule struct {
	git      *git.Repository
	repoPath string
}

// NewModule creates a new standalone git module
func NewModule(dummyRepo interface{}, path string) *StandaloneModule {
	// Ensure absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		// Fall back to provided path if abs fails
		absPath = path
	}

	module := &StandaloneModule{
		repoPath: absPath,
	}

	// Try to open existing repository
	if gitRepo, err := git.PlainOpen(absPath); err == nil {
		module.git = gitRepo
	}

	return module
}

// InitRepo initializes a new git repository
func (m *StandaloneModule) InitRepo() error {
	gitRepo, err := git.PlainInit(m.repoPath, false)
	if err != nil {
		return fmt.Errorf("initializing git repository: %w", err)
	}
	m.git = gitRepo
	return nil
}

// ShowStatus shows the working tree status
func (m *StandaloneModule) ShowStatus() error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	workTree, err := m.git.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	status, err := workTree.Status()
	if err != nil {
		return fmt.Errorf("getting status: %w", err)
	}

	if status.IsClean() {
		fmt.Println("On branch main")
		fmt.Println("nothing to commit, working tree clean")
	} else {
		fmt.Println("Changes:")
		for file, status := range status {
			fmt.Printf("  %s %s\n", status.Staging, file)
		}
	}

	return nil
}

// AddFiles adds files to the git index
func (m *StandaloneModule) AddFiles(path string) error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	workTree, err := m.git.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	// Add files
	_, err = workTree.Add(path)
	if err != nil {
		return fmt.Errorf("adding files: %w", err)
	}

	return nil
}

// Commit creates a new commit with the staged changes
func (m *StandaloneModule) Commit(message string) error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	workTree, err := m.git.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	// Create commit
	commit, err := workTree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Memex User",
			Email: "user@memex.local",
		},
	})

	if err != nil {
		return fmt.Errorf("creating commit: %w", err)
	}

	fmt.Printf("Committed: %s\n", commit.String())
	return nil
}

// ShowLog displays the commit history
func (m *StandaloneModule) ShowLog() error {
	if m.git == nil {
		return fmt.Errorf("git repository not initialized")
	}

	// Get commit history
	ref, err := m.git.Head()
	if err != nil {
		return fmt.Errorf("getting HEAD reference: %w", err)
	}

	iter, err := m.git.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return fmt.Errorf("getting log: %w", err)
	}
	defer iter.Close()

	// Print commits
	fmt.Println("Commit History:")
	err = iter.ForEach(func(commit *object.Commit) error {
		fmt.Printf("%s %s\n", commit.Hash.String()[:7], commit.Message)
		fmt.Printf("Author: %s <%s>\n", commit.Author.Name, commit.Author.Email)
		fmt.Printf("Date: %s\n\n", commit.Author.When.Format("Mon Jan 2 15:04:05 2006"))
		return nil
	})

	if err != nil {
		return fmt.Errorf("iterating commits: %w", err)
	}

	return nil
}