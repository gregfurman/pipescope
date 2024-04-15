package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/gregfurman/pipescope/internal/checker"
	"github.com/gregfurman/pipescope/internal/git"
	"github.com/gregfurman/pipescope/internal/gitlab"
)

func main() {
	faccessToken := flag.String("access-token", SetFromEnv("ACCESS_TOKEN", ""), "API access token where remote pipeline resides (env=ACCESS_TOKEN).")
	fgitDirectoryLoc := flag.String("git-directory", ".", "Location of .git directory.")
	fpollFrequency := flag.Duration("poll-frequency", 5*time.Second, "Polling frequency to pipeline.")

	// Experimental
	fplaySoundOnComplete := flag.Bool("play-sound", false, "Play a noise when pipeline completes (experimental).")

	flag.Parse()

	gitlabClient, err := gitlab.New(*faccessToken)
	if err != nil {
		panic(err)
	}

	gitClient, err := git.New(*fgitDirectoryLoc)
	if err != nil {
		panic(err)
	}

	svc := checker.New(gitlabClient, gitClient)

	pipeline, err := svc.GetPipeline()
	if err != nil {
		panic(err)
	}

	slog.Info("found pipeline",
		slog.Any("url", pipeline.URL),
		slog.Any("sha", pipeline.CommitSha),
		slog.Any("project_id", pipeline.ProjectID),
		slog.Any("pipeline_id", pipeline.ID),
	)

	statusCh, _ := svc.PollPipelineStatus(pipeline.ProjectID, pipeline.ID, *fpollFrequency)

	status := pipeline.Status
	slog.Info("Polling Pipeline",
		slog.Any("status", status),
		slog.Any("project_id", pipeline.ProjectID),
		slog.Any("pipeline_id", pipeline.ID),
	)

	for status = range statusCh {
		slog.Info("\033[1A\033[K"+"Polling Pipeline", slog.Any("status", status), slog.Any("project_id", pipeline.ProjectID), slog.Any("pipeline_id", pipeline.ID))
	}

	if *fplaySoundOnComplete {
		if err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration); err != nil {
			slog.Error("error encountered when playing sound", slog.Any("error", err))
		}
	}
}

func SetFromEnv(name, defaultValue string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}

	return defaultValue
}
