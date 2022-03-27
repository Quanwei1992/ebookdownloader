#!/bin/bash
echo "build  ebookdl_ui and pack  it with appimagetool"
echo "if you need to pack ebookdl_ui as the  ebookdl_ui.appimage you need to install appimagetool first"
echo "get appimagetool in github https://github.com/AppImage/AppImageKit/releases"
echo "after downloader appimagetool-$(uname -i).AppImage , rename it as appimagetool"
commit_id=$(git rev-parse --short HEAD)
build_time=$(date  "+%Y-%m-%d %H:%M:%S")
last_tag_commit_id=$(git rev-list --tags --max-count=1)
last_tag=$(git describe --tags $last_tag_commit_id)
cur_dir=$(pwd)
arch=$(uname -i)
go build -ldflags "-w -s -X 'main.Commit=$commit_id' -X 'main.BuildTime=$build_time' -X 'main.Version=$last_tag'" -o ebookdl_ui
cp ebookdl_ui ../ebookdownloader.AppDir/usr/bin/ebookdl_ui
rm -rf  $cur_dir/../ebookdownloader.AppDir/AppRun
ln -s $cur_dir/../ebookdownloader.AppDir/usr/bin/ebookdl_ui  $cur_dir/../ebookdownloader.AppDir/AppRun
ARCH=$arch appimagetool ../ebookdownloader.AppDir
