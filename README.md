# PipeScope
<p align="center">
  <img src="https://github.com/gregfurman/pipescope/assets/31275503/bb69ef68-0ab7-4098-87e4-480d6a880033" />
</p>

PipeScope is a command-line tool that monitors the status of remote pipelines.

## How it works
PipeScope uses [external clients](https://github.com/xanzy/go-gitlab) in conjunction with a [pure Golang git implementation](https://github.com/go-git/go-git) to dynamically stream the pipeline status of your repo's `HEAD` commit.

## Installation
```shell
go install github.com/gregfurman/pipescope@latest
```

## Usage
Be sure to specify your access token (currently only GitLab supported) via the `--access-token` flag. 

PipeScope _should_ be run wherever your repo's `.git` is located. Else, you will have to specify the location via the `--git-directory` flag.

### Example
```shell
pipescope --access-token=<GitLab access token>
```

Output:
```
2024/04/15 17:40:46 INFO found pipeline url=https://gitlab.com/gregfurman/sample-project/-/pipelines/1253625945 sha=928e2dfffdaaf6fa32f3ee4bf608690b09c6e2c1 project_id=26797650 pipeline_id=5234431782
2024/04/15 17:52:09 INFO Polling Pipeline status=created project_id=26797650 pipeline_id=5234431782
2024/04/15 17:52:24 INFO Polling Pipeline status=running project_id=26797650 pipeline_id=5234431782
2024/04/15 17:52:59 INFO Polling Pipeline status=success project_id=26797650 pipeline_id=5234431782
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
      Polling frequency to pipeline. (default 5
```

## Limitations/Roadmap
- Currently, PipeScope can only monitor GitLab pipelines. Future versions will extend this to include GitHub workflows.
- There is not a lot of flexability to select pipelines or projects via command-line input -- this should be changed to allow custom pipeline IDs to be specified.
