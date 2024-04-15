package gitlab

import (
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type Client struct {
	api *gitlab.Client
}

func New(token string) (*Client, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create new GitLab client: %w", err)
	}

	return &Client{
		api: client,
	}, nil
}

func (c *Client) GetProjectID(path string) (int, error) {
	project, _, err := c.api.Projects.GetProject(strings.TrimPrefix(strings.TrimSuffix(path, ".git"), "git@gitlab.com:"), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get project_id: %w", err)
	}

	return project.ID, nil
}

func (c *Client) GetPipelineBySha(id int, sha string) (*Pipeline, error) {
	pipelines, _, err := c.api.Pipelines.ListProjectPipelines(id, &gitlab.ListProjectPipelinesOptions{
		SHA:     gitlab.Ptr(sha),
		OrderBy: gitlab.Ptr("id"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pipelines for project_id=%d and SHA=%s: %w", id, sha, err)
	}

	if len(pipelines) == 0 {
		return nil, fmt.Errorf("no pipelines found for project_id=%d and SHA=%s", id, sha)
	}

	return &Pipeline{
		Status:    pipelines[0].Status,
		ID:        pipelines[0].ID,
		ProjectID: pipelines[0].ProjectID,
		URL:       pipelines[0].WebURL,
		CommitSha: pipelines[0].SHA,
	}, nil
}

func (c *Client) GetPipeline(id, pid int) (*Pipeline, error) {
	pipeline, _, err := c.api.Pipelines.GetPipeline(id, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pipeline %d for project_id=%d: %w", pid, id, err)
	}

	return &Pipeline{
		Status:    pipeline.Status,
		ID:        pipeline.ID,
		ProjectID: pipeline.ProjectID,
		URL:       pipeline.WebURL,
		CommitSha: pipeline.SHA,
	}, nil
}

func (c *Client) GetPipelineJobs(id, pid int) (string, error) {
	pipelines, _, err := c.api.Jobs.ListPipelineJobs(id, pid, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve pipeline jobs for project_id=%d and pipeline=%d: %w", id, pid, err)
	}

	if len(pipelines) == 0 {
		return "", fmt.Errorf("no pipeline jobs found for project_id=%d and pipeline=%d", id, pid)
	}

	return pipelines[0].Status, nil
}
