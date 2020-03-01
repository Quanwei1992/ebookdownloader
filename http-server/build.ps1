# 获取系统时间
$origin_date=Get-Date
$build_time=$origin_date.ToString('yyyy-MM-dd hh:mm:ss')
$commit_id=git rev-parse --short HEAD

$last_tag_commit_id=git rev-list --tags --max-count=1
$last_tag=git describe --tags $last_tag_commit_id
Write-Host "BuildTime = $build_time"
Write-Host "CommitID = $commit_id"
Write-Host "last_tag = $last_tag"

go build -ldflags "-X main.Commit=$commit_id -X 'main.BuildTime=$build_time' -X main.Version=$last_tag" -o ebookdownloader_http-server.exe
Copy-Item ebookdownloader_http-server.exe ..\