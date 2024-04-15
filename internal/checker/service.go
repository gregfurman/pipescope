package checker

import (
	"fmt"
	"time"

	"github.com/gregfurman/pipescope/internal/git"
	"github.com/gregfurman/pipescope/internal/gitlab"
)

type Service struct {
	gitlabClient *gitlab.Client
	gitClient    *git.Client
}

func New(glc *gitlab.Client, gc *git.Client) *Service {
	return &Service{
		gitlabClient: glc,
		gitClient:    gc,
	}
}

func (s *Service) GetPipelineStatus() (string, error) {
	pipeline, err := s.GetPipeline()
	if err != nil {
		return "", err
	}

	return pipeline.Status, nil
}

func (s *Service) GetPipeline() (*gitlab.Pipeline, error) {
	url, err := s.gitClient.GetRemoteURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote url from git: %w", err)
	}

	sha, err := s.gitClient.GetHead()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD SHA ID from git: %w", err)
	}

	id, err := s.gitlabClient.GetProjectID(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get project ID from GitLab: %w", err)
	}

	pipeline, err := s.gitlabClient.GetPipelineBySha(id, sha)
	if err != nil {
		return nil, fmt.Errorf("failed to get project ID from GitLab: %w", err)
	}

	return pipeline, nil
}

func (s *Service) GetPipelineStatusByID(id, pid int) (string, error) {
	pipeline, err := s.gitlabClient.GetPipeline(id, pid)
	if err != nil {
		return "", fmt.Errorf("failed to get project from GitLab: %w", err)
	}

	return pipeline.Status, nil
}

func (s *Service) PollPipelineStatus(id, pid int, freq time.Duration) (chan string, chan struct{}) {
	ticker := time.NewTicker(freq)

	doneCh := make(chan struct{})
	statusCh := make(chan string)

	go func() {
		for {
			select {
			case <-ticker.C:
				status, err := s.GetPipelineStatusByID(id, pid)
				statusCh <- status

				if err != nil || !gitlab.IsStatusPending(status) {
					close(doneCh)
				}
			case <-doneCh:
				close(statusCh)
				ticker.Stop()

				return
			}
		}
	}()

	return statusCh, doneCh
}
