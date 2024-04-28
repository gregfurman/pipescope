package gateway

import (
	"net/http"
	"testing"

	"github.com/xanzy/go-gitlab"
)

func Test_GitLab_GetProjectID(t *testing.T) {

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
			mockedAPI, _ := gitlab.NewClient("", gitlab.WithHTTPClient(&http.Client{Transport: &mockRoundTripper{makeJSONResponse(tt.mockedResponseBody)}}))

			client := GitLabClient{
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

func Test_GitLab_GetPipeline(t *testing.T) {

	tests := []struct {
		name               string
		mockedResponseBody string

		path       string
		pipelineID int

		want    Pipeline
		wantErr bool
	}{
		{
			name: "Successfully returns in_progress pipeline",
			mockedResponseBody: `{
				"id": 46,
				"iid": 11,
				"project_id": 1,
				"name": "Build pipeline",
				"status": "pending",
				"ref": "main",
				"sha": "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				"before_sha": "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				"tag": false,
				"yaml_errors": null,
				"user": {
				  "name": "Administrator",
				  "username": "root",
				  "id": 1,
				  "state": "active",
				  "avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				  "web_url": "http://localhost:3000/root"
				},
				"created_at": "2016-08-11T11:28:34.085Z",
				"updated_at": "2016-08-11T11:32:35.169Z",
				"started_at": null,
				"finished_at": "2016-08-11T11:32:35.145Z",
				"committed_at": null,
				"duration": 123,
				"queued_duration": 1,
				"coverage": "30.0",
				"web_url": "https://example.com/gregfurman/pipescope/pipelines/46"
			  }`,
			pipelineID: 46,
			path:       "https://github.com/gregfurman/pipescope",
			want: Pipeline{
				ID:        46,
				ProjectID: "1",
				CommitSha: "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				Status:    "pending",
				URL:       "https://example.com/gregfurman/pipescope/pipelines/46",
			},
		},
		{
			name: "Successfully returns completed and failed workflow",
			mockedResponseBody: `{
				"id": 46,
				"iid": 11,
				"project_id": 1,
				"name": "Build pipeline",
				"status": "failed",
				"ref": "main",
				"sha": "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				"before_sha": "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				"tag": false,
				"yaml_errors": null,
				"user": {
				  "name": "Administrator",
				  "username": "root",
				  "id": 1,
				  "state": "active",
				  "avatar_url": "http://www.gravatar.com/avatar/e64c7d89f26bd1972efa854d13d7dd61?s=80&d=identicon",
				  "web_url": "http://localhost:3000/root"
				},
				"created_at": "2016-08-11T11:28:34.085Z",
				"updated_at": "2016-08-11T11:32:35.169Z",
				"started_at": null,
				"finished_at": "2016-08-11T11:32:35.145Z",
				"committed_at": null,
				"duration": 123,
				"queued_duration": 1,
				"coverage": "30.0",
				"web_url": "https://example.com/gregfurman/pipescope/pipelines/46"
			  }`,
			pipelineID: 46,
			path:       "gregfurman/pipescope",
			want: Pipeline{
				ID:        46,
				ProjectID: "1",
				CommitSha: "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				Status:    "failed",
				URL:       "https://example.com/gregfurman/pipescope/pipelines/46",
			},
		},
		{
			name:               "Fails due to no workflows found",
			mockedResponseBody: `{}`,
			pipelineID:         46,
			path:               "https://example.com/gregfurman/pipescope",
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			mockedAPI, _ := gitlab.NewClient("", gitlab.WithHTTPClient(&http.Client{Transport: &mockRoundTripper{makeJSONResponse(tt.mockedResponseBody)}}))

			client := GitLabClient{
				api: mockedAPI,
			}
			got, err := client.GetPipeline(tt.path, tt.pipelineID)
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

func Test_GitLab_GetPipelineBySha(t *testing.T) {

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
			mockedResponseBody: `[
				{
					"id": 47,
					"iid": 12,
					"project_id": 1,
					"status": "pending",
					"source": "push",
					"ref": "new-pipeline",
					"sha": "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
					"name": "Build pipeline",
					"web_url": "https://example.com/gregfurman/pipescope/pipelines/47",
					"created_at": "2016-08-11T11:28:34.085Z",
					"updated_at": "2016-08-11T11:32:35.169Z"
				}
			]`,
			sha:  "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
			path: "https://github.com/gregfurman/pipescope",
			want: Pipeline{
				ID:        47,
				ProjectID: "1",
				CommitSha: "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				Status:    "pending",
				URL:       "https://example.com/gregfurman/pipescope/pipelines/47",
			},
		},
		{
			name: "Successfully returns completed and failed workflow",
			mockedResponseBody: `[
				{
					"id": 48,
					"iid": 13,
					"project_id": 1,
					"status": "failed",
					"source": "web",
					"ref": "new-pipeline",
					"sha": "eb94b618fb5865b26e80fdd8ae531b7a63ad851a",
					"name": "Build pipeline",
					"web_url": "https://example.com/gregfurman/pipescope/pipelines/48",
					"created_at": "2016-08-12T10:06:04.561Z",
					"updated_at": "2016-08-12T10:09:56.223Z"
				},
				{
				  "id": 47,
				  "iid": 12,
				  "project_id": 1,
				  "status": "pending",
				  "source": "push",
				  "ref": "new-pipeline",
				  "sha": "a91957a858320c0e17f3a0eca7cfacbff50ea29a",
				  "name": "Build pipeline",
				  "web_url": "https://example.com/gregfurman/pipescope/pipelines/47",
				  "created_at": "2016-08-11T11:28:34.085Z",
				  "updated_at": "2016-08-11T11:32:35.169Z"
				}
			  ]`,
			path: "https://example.com/gregfurman/pipescope/pipelines/48",
			sha:  "eb94b618fb5865b26e80fdd8ae531b7a63ad851a",
			want: Pipeline{
				ID:        48,
				ProjectID: "1",
				CommitSha: "eb94b618fb5865b26e80fdd8ae531b7a63ad851a",
				Status:    "failed",
				URL:       "https://example.com/gregfurman/pipescope/pipelines/48",
			},
		},
		{
			name:               "Fails due to no workflows found",
			mockedResponseBody: `[]`,
			path:               "https://example.com/gregfurman/pipescope",
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			mockedAPI, _ := gitlab.NewClient("", gitlab.WithHTTPClient(&http.Client{Transport: &mockRoundTripper{makeJSONResponse(tt.mockedResponseBody)}}))

			client := GitLabClient{
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
