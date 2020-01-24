#!/bin/bash

CommitID=$(git rev-parse HEAD)
LDFlags="-X main.Commit=${CommitID}"

 go build  -ldflags "${LDFlags}" -o ebookdownloader_cli
 cp ebookdownloader_cli ../