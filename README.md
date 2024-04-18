# PipeScope
<p align="center">
  <img src="https://github.com/gregfurman/pipescope/assets/31275503/bb69ef68-0ab7-4098-87e4-480d6a880033" />
</p>

PipeScope is a command-line tool that monitors the status of remote pipelines.

## How it works
PipeScope uses external [GitLab](https://github.com/xanzy/go-gitlab) and [GitHub](https://github.com/google/go-github) clients in conjunction with a [pure Golang git implementation](https://github.com/go-git/go-git) to dynamically stream the pipeline status of your repo's `HEAD` commit.

## Installation
```shell
go install github.com/gregfurman/pipescope@latest
```

## Usage
Be sure to specify your access token via the `--access-token` flag. 

PipeScope _should_ be run wherever your repo's `.git` is located. Else, you will have to specify the location via the `--git-directory` flag.

The Git Provider can be explicitly set as an argument passed to the CLI (can be one of `gitlab` or `github`). Otherwise, PipeScope will attempt to dynamically determine which provider to use from the prefix of your `access-token` OR from the remote URL retrieved via the internal `git` client.

### Example
```shell
pipescope --access-token=<GitLab access token>
```

Output:
```
2024/04/18 17:55:47 INFO Polled Pipeline [status=created]" url=https://gitlab.com/gregfurman/sample-project/-/pipelines/1253625945 sha=928e2dfffdaaf6fa32f3ee4bf608690b09c6e2c1 project_id=26797650 pipeline_id=5234431782
2024/04/18 17:55:58 INFO Polled Pipeline [status=running]" url=https://gitlab.com/gregfurman/sample-project/-/pipelines/1253625945 sha=928e2dfffdaaf6fa32f3ee4bf608690b09c6e2c1 project_id=26797650 pipeline_id=5234431782
2024/04/18 17:56:12 INFO Polled Pipeline [status=success]" url=https://gitlab.com/gregfurman/sample-project/-/pipelines/1253625945 sha=928e2dfffdaaf6fa32f3ee4bf608690b09c6e2c1 project_id=26797650 pipeline_id=5234431782
```

### Command-line Flags
```shell
-access-token string
      API access token where remote pipeline resides (env=ACCESS_TOKEN).
-git-directory string
      Location of .git directory. (default ".")
-play-sound
      Play a noise when pipeline completes (experimental).
-poll-frequency duration
      Polling frequency to pipeline. (default 5s)
```

## Limitations/Roadmap
- <b>Project is currently experimental</b> so please do not use this in production anywhere
- Needs CI/CD, tests, and a Makefile
- <s>Currently, PipeScope can only monitor GitLab pipelines. Future versions will extend this to include GitHub workflows.</s>
- There is not a lot of flexibility to select pipelines or projects via command-line input -- this should be changed to allow custom pipeline IDs to be specified.
- Include more information on pipeline jobs
- Allow better streaming of pipeline/job logs to stdout
