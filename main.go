package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/gregfurman/pipescope/internal/checker"
	"github.com/gregfurman/pipescope/internal/gateway"
	"github.com/gregfurman/pipescope/internal/git"
)

func main() {
	faccessToken := flag.String("access-token", setFromEnv("ACCESS_TOKEN", ""), "API access token where remote pipeline resides (env=ACCESS_TOKEN).")
	fgitDirectoryLoc := flag.String("git-directory", ".", "Location of .git directory.")
	fpollFrequency := flag.Duration("poll-frequency", 5*time.Second, "Polling frequency to pipeline.")

	// Experimental
	fplaySoundOnComplete := flag.Bool("play-sound", false, "Play a noise when pipeline completes (experimental).")

	flag.Parse()

	// Define clients
	gitClient, err := git.New(*fgitDirectoryLoc)
	if err != nil {
		exit(err)
	}

	var gatewayClient gateway.Client

	switch {
	// If an arg is passed in, use it to determine the git provider
	case flag.Arg(0) != "":
		gatewayClient, err = gateway.New(*faccessToken, gateway.ProviderType(flag.Arg(0)))
	// Check the access token's prefix to determine the git provider
	case *faccessToken != "":
		gatewayClient, err = gateway.NewFromToken(*faccessToken)
	// Check the remote git URL to determine the git provider
	default:
		url, _ := gitClient.GetRemoteURL()
		gatewayClient, err = gateway.NewFromRemoteURL(*faccessToken, url)
	}

	if err != nil {
		exit(err)
	}

	// Define service
	svc := checker.New(gatewayClient, gitClient)

	// Run polling
	if err := run(svc, *fpollFrequency); err != nil {
		exit(err)
	}

	if *fplaySoundOnComplete {
		if err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration); err != nil {
			slog.Error("error encountered when playing sound", slog.Any("error", err))
		}
	}
}

func run(svc *checker.Service, pollFrequency time.Duration) error {
	pipeline, err := svc.GetPipeline()
	if err != nil {
		return fmt.Errorf("checker Service GET pipeline failed: %w", err)
	}
	status := pipeline.Status

	slog.Group("pipeline")

	logger := slog.New(slog.Default().Handler()).With(
		slog.Any("url", pipeline.URL),
		slog.Any("sha", pipeline.CommitSha),
		slog.Any("project_id", pipeline.ProjectID),
		slog.Any("pipeline_id", pipeline.ID),
	)

	logger.Info(fmt.Sprintf("Polled Pipeline [status=%s]", pipeline.Status))
	statusCh, _ := svc.PollPipelineStatus(pipeline.ProjectID, pipeline.ID, pollFrequency)
	for s := range statusCh {
		if status == "" {
			continue
		}

		if s != status {
			status = s
			logger.Info(fmt.Sprintf("Polled Pipeline [status=%s]", s))
		}
	}

	return nil
}

func setFromEnv(name, defaultValue string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}

	return defaultValue
}

func exit(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}
