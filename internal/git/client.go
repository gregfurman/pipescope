package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

type Client interface {
	GetHead() (string, error)
	GetRemoteURL() (string, error)
}

type ClientImpl struct {
	repo *git.Repository
}

func New(path string) (*ClientImpl, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Git repository: %w", err)
	}

	return &ClientImpl{
		repo: repo,
	}, nil
}

func (c *ClientImpl) GetHead() (string, error) {
	head, err := c.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD of repository: %w", err)
	}

	return head.Hash().String(), nil
}

func (c *ClientImpl) GetRemoteURL() (string, error) {
	remote, err := c.repo.Remote("origin")
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL of repository: %w", err)
	}

	return remote.Config().URLs[0], nil
}
