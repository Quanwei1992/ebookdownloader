#!/bin/bash

CommitID=$(git rev-parse HEAD)
LDFlags="-H windowsgui -w -s -X main.Commit=${CommitID}' -linkmode external -extldflags '-static'"

 go build  -ldflags "${LDFlags}" -o ebookdownloader_gui
 cp ebookdownloader_gui ../