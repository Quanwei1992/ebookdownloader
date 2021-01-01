#!/bin/bash

# export exec name
Exec=ebookdownloader
# Set Release folder
ReleaseFolder=ebookdownloader.AppDir
BuildPath=$GOPATH/src/github.com/sndnvaps/ebookdownloader
des=$BuildPath/$ReleaseFolder
desBin=$des/usr/bin
VERSION=$(git rev-parse --short HEAD)
BuildTime=$(date +'%Y-%m-%d_%T')
ARCH=x86_64


# Start Build Process
function Build {
    cd $BuildPath/qtgui
    go build -ldflags " -X main.Version=${VERSION} -X main.Commit=${VERSION} -X main.BuildTime=${BUILDTime}" -o $desBin/ebookdownloader_gui
   cd $BuildPath
}

# Clean the Exec file
function Clean {
    if [ -f $desBin/ebookdownloader_gui ]
    then
     rm -rf  $desBin/*
     fi
}

# Copy the needed file into Release folder
function CopyFiles {
    mkdir -p $desBin/tools
    cp   $BuildPath/tools/kindlegenLinux $desBin/tools/kindlegenLinux
}

function PackAppImage {
    cp /usr/bin/desktop-file-validate $desBin/
    ARCH=x86_64 linuxdeployqt $des/usr/share/applications/$Exec.desktop -appimage \
     -executable=$des/usr/bin/desktop-file-validate
}

Clean
Build
CopyFiles
PackAppImage