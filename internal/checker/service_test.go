package checker

import (
	"sync"
	"testing"
	"time"

	"github.com/gregfurman/pipescope/internal/gateway"
)

type providerMock struct {
	mockedGetPipelineBySha func(id, sha string) (*gateway.Pipeline, error)
	mockedGetPipeline      func(id string, pid int) (*gateway.Pipeline, error)
	mockedIsStatusPending  func(status string) bool
	lock                   sync.Mutex
}

func (pm *providerMock) GetPipelineBySha(id, sha string) (*gateway.Pipeline, error) {
	return pm.mockedGetPipelineBySha(id, sha)
}

func (pm *providerMock) GetPipeline(id string, pid int) (*gateway.Pipeline, error) {
	pm.lock.Lock()
	pipeline, err := pm.mockedGetPipeline(id, pid)
	pm.lock.Unlock()

	return pipeline, err
}

func (pm *providerMock) IsStatusPending(status string) bool {
	pm.lock.Lock()
	flag := pm.mockedIsStatusPending(status)
	pm.lock.Unlock()
	return flag
}

type gitMock struct {
	mockedGetHead      func() (string, error)
	mockedGetRemoteURL func() (string, error)
}

func (gm *gitMock) GetHead() (string, error) {
	return gm.mockedGetHead()
}
func (gm *gitMock) GetRemoteURL() (string, error) {
	return gm.mockedGetRemoteURL()
}

func Test_Service(t *testing.T) {
	expectedPipeline := gateway.Pipeline{
		ID:        1,
		ProjectID: "PROJECT_ID",
		Status:    "pending",
		URL:       "www.example.com/repo/owner/pipelines/1",
		CommitSha: "COMMIT_SHA",
	}

	providerClient := &providerMock{
		mockedGetPipelineBySha: func(id, sha string) (*gateway.Pipeline, error) { return &expectedPipeline, nil },
		mockedGetPipeline:      func(id string, pid int) (*gateway.Pipeline, error) { return &expectedPipeline, nil },
		mockedIsStatusPending:  func(status string) bool { return true },
		lock:                   sync.Mutex{},
	}

	gitClient := &gitMock{
		mockedGetHead:      func() (string, error) { return "COMMIT_SHA", nil },
		mockedGetRemoteURL: func() (string, error) { return "www.example.com/repo/owner", nil },
	}

	svc := New(providerClient, gitClient)

	pipeline, err := svc.GetPipeline()
	if err != nil {
		t.Error("did not expect error")
	}

	if *pipeline != expectedPipeline {
		t.Errorf("expected pipeline object %v, got %v", *pipeline, expectedPipeline)
	}

	status, err := svc.GetPipelineStatus()
	if err != nil {
		t.Error("did not expect error")
	}

	if status != pipeline.Status {
		t.Errorf("expected pipeline status %v, got %v", *pipeline, expectedPipeline)
	}

	// Poll every 500ms
	statusCh, _ := svc.PollPipelineStatus("PROJECT_ID", 1, 500*time.Millisecond)

	// Change the status to "success" after 1000ms
	time.AfterFunc(1000*time.Millisecond, func() {
		providerClient.lock.Lock()
		providerClient.mockedIsStatusPending = func(status string) bool { return false }
		providerClient.mockedGetPipeline = func(id string, pid int) (*gateway.Pipeline, error) {
			return &gateway.Pipeline{
				ID:        1,
				ProjectID: "PROJECT_ID",
				Status:    "success",
				URL:       "www.example.com/repo/owner/pipelines/1",
				CommitSha: "COMMIT_SHA",
			}, nil
		}
		providerClient.lock.Unlock()
	})

	var gotStatus string
	for s := range statusCh {
		gotStatus = s
		if svc.gatewayClient.IsStatusPending(gotStatus) && gotStatus != "pending" {
			t.Errorf("expected status to be pending, got %s", gotStatus)
		}
	}

	if gotStatus != "success" {
		t.Errorf("expected status to be success, got %s", gotStatus)
	}

}
