package gateway

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitLabClient struct {
	api *gitlab.Client
}

func NewGitLabClient(token string) (*GitLabClient, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create new GitLab client: %w", err)
	}

	return &GitLabClient{
		api: client,
	}, nil
}

func (c *GitLabClient) GetProjectID(path string) (int, error) {
	project, _, err := c.api.Projects.GetProject(strings.TrimPrefix(strings.TrimSuffix(path, ".git"), "git@gitlab.com:"), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get project_id: %w", err)
	}

	return project.ID, nil
}

func (c *GitLabClient) GetPipelineBySha(id string, sha string) (*Pipeline, error) {
	pipelines, _, err := c.api.Pipelines.ListProjectPipelines(strings.TrimPrefix(strings.TrimSuffix(id, ".git"), "git@gitlab.com:"), &gitlab.ListProjectPipelinesOptions{
		SHA:     gitlab.Ptr(sha),
		OrderBy: gitlab.Ptr("id"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pipelines: %w", err)
	}

	if len(pipelines) == 0 {
		return nil, fmt.Errorf("no pipelines found")
	}

	return &Pipeline{
		Status:    pipelines[0].Status,
		ID:        pipelines[0].ID,
		ProjectID: strconv.Itoa(pipelines[0].ProjectID),
		URL:       pipelines[0].WebURL,
		CommitSha: pipelines[0].SHA,
	}, nil
}

func (c *GitLabClient) GetPipeline(id string, pid int) (*Pipeline, error) {
	pipeline, _, err := c.api.Pipelines.GetPipeline(id, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pipeline: %w", err)
	}

	return &Pipeline{
		Status:    pipeline.Status,
		ID:        pipeline.ID,
		ProjectID: strconv.Itoa(pipeline.ProjectID),
		URL:       pipeline.WebURL,
		CommitSha: pipeline.SHA,
	}, nil
}

func (*GitLabClient) IsStatusPending(status string) bool {
	switch gitlab.BuildStateValue(status) {
	case gitlab.Created, gitlab.WaitingForResource, gitlab.Preparing, gitlab.Pending, gitlab.Running:
		return true
	}

	return false
}
