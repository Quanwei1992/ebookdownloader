#!/bin/bash

CommitID=$(git rev-parse HEAD)
BuildTime=$(date +%Y-%m-%d\ %H:%M)
LDFlags="-H windowsgui -w -s -X main.Commit=${CommitID}  -X 'main.BuildTime=${BuildTime}'"

 go build  -ldflags "${LDFlags}" -o ebookdownloader_gui
 cp ebookdownloader_gui ../