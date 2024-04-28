package checker

import (
	"fmt"
	"time"

	"github.com/gregfurman/pipescope/internal/gateway"
	"github.com/gregfurman/pipescope/internal/git"
)

type Service struct {
	gatewayClient gateway.Client
	gitClient     git.Client
}

func New(gw gateway.Client, gc git.Client) *Service {
	return &Service{
		gatewayClient: gw,
		gitClient:     gc,
	}
}

func (s *Service) GetPipelineStatus() (string, error) {
	pipeline, err := s.GetPipeline()
	if err != nil {
		return "", err
	}

	return pipeline.Status, nil
}

func (s *Service) GetPipeline() (*gateway.Pipeline, error) {
	url, err := s.gitClient.GetRemoteURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote url from git: %w", err)
	}

	sha, err := s.gitClient.GetHead()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD SHA ID from git: %w", err)
	}

	pipeline, err := s.gatewayClient.GetPipelineBySha(url, sha)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline from Gateway client: %w", err)
	}

	if pipeline.CommitSha == "" {
		pipeline.CommitSha = sha
	}

	return pipeline, nil
}

func (s *Service) GetPipelineStatusByID(id string, pid int) (string, error) {
	pipeline, err := s.gatewayClient.GetPipeline(id, pid)
	if err != nil {
		return "", fmt.Errorf("failed to get project: %w", err)
	}

	return pipeline.Status, nil
}

func (s *Service) PollPipelineStatus(id string, pid int, freq time.Duration) (chan string, chan struct{}) {
	ticker := time.NewTicker(freq)

	doneCh := make(chan struct{})
	statusCh := make(chan string)

	go func() {
		for {
			select {
			case <-ticker.C:
				status, err := s.GetPipelineStatusByID(id, pid)
				statusCh <- status

				if err != nil || !s.gatewayClient.IsStatusPending(status) {
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
