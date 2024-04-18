package gateway

type Client interface {
	GetPipelineBySha(id, sha string) (*Pipeline, error)
	GetPipeline(id string, pid int) (*Pipeline, error)
	IsStatusPending(status string) bool
}

type Pipeline struct {
	ID        int
	ProjectID string
	CommitSha string
	Status    string
	URL       string
}

type ProviderType string

const (
	GitHub ProviderType = "github"
	GitLab ProviderType = "gitlab"
)

// GitHub access token types typically used with the REST API
// see https://github.blog/2021-04-05-behind-githubs-new-authentication-token-formats/
const (
	// ghp for GitHub personal access tokens.
	GitHubPersonalAccessToken = "ghp_"

	// gho for OAuth access tokens.
	GitHubOAuth = "gho_"

	// ghu for GitHub user-to-server tokens.
	GitHubUserToServer = "ghu_"
)

// GitLab access token types typically used with the REST API
// see https://docs.gitlab.com/ee/security/token_overview.html#token-prefixes
const (
	// glpat- for GitLab personal access tokens.
	GitLabPersonalAccessToken = "glpat-"

	// gloas- for GitLab OAuth Application Secret.
	GitLabOAuthApplicationSecret = "gloas-"

	// gldt- for GitLab Deploy Token.
	GitLabDeployToken = "gldt-"

	// glrt- for GitLab Runner Authentication Token.
	GitLabRunnerAuthenticationToken = "glrt-"

	// glcbt- for GitLab CI/CD Job Token.
	GitLabCICDJobToken = "glcbt-"

	// glptt- for GitLab Trigger Token.
	GitLabTriggerToken = "glptt-"

	// glft- for GitLab Feed Token.
	GitLabFeedToken = "glft-"

	// glimt- for GitLab Incoming Mail Token.
	GitLabIncomingMailToken = "glimt-"

	// glagent- for GitLab Agent for Kubernetes Token.
	GitLabAgentForKubernetesToken = "glagent-"

	// glsoat- for GitLab SCIM Token.
	GitLabSCIMToken = "glsoat-"
)
