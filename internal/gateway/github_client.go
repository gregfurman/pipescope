package gateway

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v61/github"
)

type GitHubClient struct {
	api *github.Client
}

func NewGitHubClient(token string) (*GitHubClient, error) {
	client := github.NewClient(nil).WithAuthToken(token)

	return &GitHubClient{
		api: client,
	}, nil
}

func (c *GitHubClient) GetProjectID(path string) (int, error) {
	owner, repo, ok := strings.Cut(strings.TrimPrefix(strings.TrimSuffix(path, ".git"), "git@github.com:"), "/")
	if !ok {
		return 0, fmt.Errorf("malformed repository path: 'owner' and 'repo' could not be extracted from %s", path)
	}

	project, _, err := c.api.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return 0, fmt.Errorf("failed to get project_id from GitHub: %w", err)
	}

	return int(project.GetID()), nil
}

func (c *GitHubClient) GetPipelineBySha(id, sha string) (*Pipeline, error) {
	owner, repo, ok := strings.Cut(strings.TrimPrefix(strings.TrimSuffix(id, ".git"), "git@github.com:"), "/")
	if !ok {
		return nil, fmt.Errorf("malformed repository path: 'owner' and 'repo' could not be extracted from %s", id)
	}

	runs, _, err := c.api.Actions.ListRepositoryWorkflowRuns(context.Background(), owner, repo, &github.ListWorkflowRunsOptions{
		HeadSHA:     sha,
		ListOptions: github.ListOptions{Page: 1, PerPage: 1},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pipelines from GitHub: %w", err)
	}

	if len(runs.WorkflowRuns) == 0 {
		return nil, fmt.Errorf("no pipelines found for project %s@%s", id, sha)
	}

	return workflowToPipeline(runs.WorkflowRuns[0]), nil
}

func (c *GitHubClient) GetPipeline(id string, pid int) (*Pipeline, error) {
	owner, repo, ok := strings.Cut(strings.TrimPrefix(strings.TrimSuffix(id, ".git"), "git@github.com:"), "/")
	if !ok {
		return nil, fmt.Errorf("malformed repository path: 'owner' and 'repo' could not be extracted from %s", id)
	}

	workflow, _, err := c.api.Actions.GetWorkflowRunByID(context.Background(), owner, repo, int64(pid))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pipeline from GitHub: %w", err)
	}

	return workflowToPipeline(workflow), nil
}

func workflowToPipeline(wf *github.WorkflowRun) *Pipeline {
	status := wf.GetStatus()
	if status == "completed" && wf.Conclusion != nil {
		status = wf.GetConclusion()
	}

	return &Pipeline{
		Status:    status,
		ID:        int(wf.GetWorkflowID()),
		ProjectID: strconv.FormatInt(wf.GetID(), 10),
		URL:       wf.GetHTMLURL(),
		CommitSha: wf.GetHeadCommit().GetSHA(),
	}
}

func (*GitHubClient) IsStatusPending(status string) bool {
	switch status {
	case "queued", "in_progress":
		return true
	}

	return false
}
