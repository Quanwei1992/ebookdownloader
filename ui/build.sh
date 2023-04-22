#!/bin/bash
commit_id=$(git rev-parse --short HEAD)
build_time=$(date  "+%Y-%m-%d %H:%M:%S")
last_tag_commit_id=$(git rev-list --tags --max-count=1)
last_tag=$(git describe --tags $last_tag_commit_id)

go build -ldflags "-w -s -X 'main.Commit=$commit_id' -X 'main.BuildTime=$build_time' -X 'main.Version=$last_tag'" -o ebookdl_ui
