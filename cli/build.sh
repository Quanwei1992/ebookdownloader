#!/bin/bash

CommitID=$(git rev-parse HEAD)
BuildTime=$(date +%Y-%m-%d\ %H:%M)
LDFlags="-X main.Commit=${CommitID}  -X 'main.BuildTime=${BuildTime}'"

 go build  -ldflags "${LDFlags}" -o ebookdownloader_cli
 cp ebookdownloader_cli ../