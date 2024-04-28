package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v61/github"
)

type mockRoundTripper struct {
	response *http.Response
}

func (rt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.response, nil
}

func makeJSONResponse(body string) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.Header().Add("Content-Type", "application/json")
	recorder.WriteString(body)
	return recorder.Result()
}

func Test_GitHub_GetProjectID(t *testing.T) {

	tests := []struct {
		name               string
		mockedResponseBody string

		path string

		want    int
		wantErr bool
	}{
		{name: "Success", mockedResponseBody: `{"id":12345}`, want: 12345, path: "owner/repo"},
		{name: "Fails due to malformed path", wantErr: true, path: "incorrect path format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			mockedAPI := github.NewClient(&http.Client{Transport: &mockRoundTripper{makeJSONResponse(tt.mockedResponseBody)}})

			client := GitHubClient{
				api: mockedAPI,
			}

			got, err := client.GetProjectID(tt.path)
			if tt.wantErr {
				if err == nil {
					ts.Error("expected an error, got nil")
				}
				return
			}

			if tt.want != got {
				ts.Errorf("expected %d, got %d", tt.want, got)
			}

		})
	}

}

func Test_GitHub_GetPipelineBySha(t *testing.T) {

	tests := []struct {
		name               string
		mockedResponseBody string

		path string
		sha  string

		want    Pipeline
		wantErr bool
	}{
		{
			name: "Successfully returns in_progress pipeline",
			mockedResponseBody: `{"total_count":1, "workflow_runs": [
				{
					"id": 8858984663,
					"head_branch": "test-workflow",
					"head_sha": "23ebbb3b14c9a026199474d2931bdc55863dfffc",
					"event": "push",
					"status": "in_progress",
					"workflow_id": 95717003,
					"html_url": "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
					"repository": {
					  "full_name": "gregfurman/pipescope"
					}
				  }
			]}`,
			sha:  "23ebbb3b14c9a026199474d2931bdc55863dfffc",
			path: "https://github.com/gregfurman/pipescope",
			want: Pipeline{
				ID:        8858984663,
				ProjectID: "gregfurman/pipescope",
				CommitSha: "23ebbb3b14c9a026199474d2931bdc55863dfffc",
				Status:    "in_progress",
				URL:       "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
			},
		},
		{
			name: "Successfully returns completed and failed workflow",
			mockedResponseBody: `{"total_count":1, "workflow_runs": [
				{
					"id": 8858984663,
					"head_branch": "test-workflow",
					"head_sha": "23ebbb3b14c9a026199474d2931bdc55863dfffc",
					"event": "push",
					"status": "completed",
					"conclusion": "failure",
					"workflow_id": 95717003,
					"html_url": "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
					"repository": {
					  "full_name": "gregfurman/pipescope"
					}
				  }
			]}`,
			sha:  "23ebbb3b14c9a026199474d2931bdc55863dfffc",
			path: "https://github.com/gregfurman/pipescope",
			want: Pipeline{
				ID:        8858984663,
				ProjectID: "gregfurman/pipescope",
				CommitSha: "23ebbb3b14c9a026199474d2931bdc55863dfffc",
				Status:    "failure",
				URL:       "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
			},
		},
		{
			name:               "Fails due to no workflows found",
			mockedResponseBody: `{"total_count":0, "workflow_runs": []}`,
			sha:                "23ebbb3b14c9a026199474d2931bdc55863dfffc",
			path:               "https://github.com/gregfurman/pipescope",
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			mockedAPI := github.NewClient(&http.Client{Transport: &mockRoundTripper{makeJSONResponse(tt.mockedResponseBody)}})

			client := GitHubClient{
				api: mockedAPI,
			}
			got, err := client.GetPipelineBySha(tt.path, tt.sha)
			if err != nil && !tt.wantErr {
				ts.Errorf("unexpected error occurred. expected nil, got %s", err)
				return
			}

			if tt.wantErr {
				if err == nil {
					ts.Error("expected an error, got nil")
				}
				return
			}

			if tt.want != *got {
				ts.Errorf("expected %v, got %v", tt.want, got)
			}

		})
	}

}

func Test_GitHub_GetPipeline(t *testing.T) {

	tests := []struct {
		name               string
		mockedResponseBody string

		path       string
		workflowID int

		want    Pipeline
		wantErr bool
	}{
		{
			name: "Successfully returns in_progress pipeline",
			mockedResponseBody: `{
					"id": 8858984663,
					"head_branch": "test-workflow",
					"head_sha": "23ebbb3b14c9a026199474d2931bdc55863dfffc",
					"event": "push",
					"status": "in_progress",
					"workflow_id": 95717003,
					"html_url": "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
					"repository": {
					  "full_name": "gregfurman/pipescope"
					}
				  }`,
			workflowID: 8858984663,
			path:       "https://github.com/gregfurman/pipescope",
			want: Pipeline{
				ID:        8858984663,
				ProjectID: "gregfurman/pipescope",
				CommitSha: "23ebbb3b14c9a026199474d2931bdc55863dfffc",
				Status:    "in_progress",
				URL:       "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
			},
		},
		{
			name: "Successfully returns completed and failed workflow",
			mockedResponseBody: `{
					"id": 8858984663,
					"head_branch": "test-workflow",
					"head_sha": "23ebbb3b14c9a026199474d2931bdc55863dfffc",
					"event": "push",
					"status": "completed",
					"conclusion": "failure",
					"workflow_id": 95717003,
					"html_url": "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
					"repository": {
					  "full_name": "gregfurman/pipescope"
					}
				  }`,
			workflowID: 8858984663,
			path:       "https://github.com/gregfurman/pipescope",
			want: Pipeline{
				ID:        8858984663,
				ProjectID: "gregfurman/pipescope",
				CommitSha: "23ebbb3b14c9a026199474d2931bdc55863dfffc",
				Status:    "failure",
				URL:       "https://github.com/gregfurman/pipescope/actions/runs/8858984663",
			},
		},
		{
			name:               "Fails due to no workflows found",
			mockedResponseBody: `{}`,
			workflowID:         8858984663,
			path:               "https://github.com/gregfurman/pipescope",
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			mockedAPI := github.NewClient(&http.Client{Transport: &mockRoundTripper{makeJSONResponse(tt.mockedResponseBody)}})

			client := GitHubClient{
				api: mockedAPI,
			}
			got, err := client.GetPipeline(tt.path, tt.workflowID)
			if err != nil && !tt.wantErr {
				ts.Errorf("unexpected error occurred. expected nil, got %s", err)
				return
			}

			if tt.wantErr {
				if err == nil {
					ts.Error("expected an error, got nil")
				}
				return
			}

			if tt.want != *got {
				ts.Errorf("expected %v, got %v", tt.want, got)
			}

		})
	}

}
