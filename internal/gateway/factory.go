package gateway

import (
	"errors"
	"strings"
)

func New(token string, t ProviderType) (Client, error) { //nolint:ireturn
	switch t {
	case GitHub:
		return NewGitHubClient(token)
	case GitLab:
		return NewGitLabClient(token)
	default:
		return nil, errors.New("gateway client does not exist")
	}
}

func NewFromToken(token string) (Client, error) { //nolint:ireturn
	switch {
	case isGitHubToken(token):
		return NewGitHubClient(token)
	case isGitLabToken(token):
		return NewGitLabClient(token)

	default:
		return nil, errors.New("could not determine gateway client from access token")
	}
}

func NewFromRemoteURL(token, remoteURL string) (Client, error) { //nolint:ireturn
	switch {
	case strings.Contains(remoteURL, "github.com"):
		return NewGitHubClient(token)
	case strings.Contains(remoteURL, "gitlab.com"):
		return NewGitLabClient(token)
	default:
		return nil, errors.New("could not determine gateway client from git repository's remote url")
	}
}

func isGitHubToken(token string) bool {
	for _, prefix := range []string{
		GitHubPersonalAccessToken,
		GitHubOAuth,
		GitHubUserToServer,
	} {
		if strings.HasPrefix(token, prefix) {
			return true
		}
	}

	return false
}

func isGitLabToken(token string) bool {
	for _, prefix := range []string{
		GitLabPersonalAccessToken,
		GitLabOAuthApplicationSecret,
		GitLabDeployToken,
		GitLabRunnerAuthenticationToken,
		GitLabCICDJobToken,
		GitLabTriggerToken,
		GitLabFeedToken,
		GitLabIncomingMailToken,
		GitLabAgentForKubernetesToken,
		GitLabSCIMToken,
	} {
		if strings.HasPrefix(token, prefix) {
			return true
		}
	}

	return false
}
