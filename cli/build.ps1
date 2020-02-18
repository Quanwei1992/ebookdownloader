
# 获取系统时间
$origin_date=Get-Date
$build_time=$origin_date.ToString('yyyy-MM-dd hh:mm:ss')
$commit_id=git rev-parse --short HEAD

go build -ldflags "-X main.Commit=$commit_id -X 'main.BuildTime=$build_time'" -o ebookdownloader_cli.exe
Copy-Item ebookdownloader_cli.exe ..\