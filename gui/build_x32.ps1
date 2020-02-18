# 获取系统时间
$origin_date=Get-Date
$build_time=$origin_date.ToString('yyyy-MM-dd hh:mm:ss')
$commit_id=git rev-parse --short HEAD

rsrc -manifest ebookdownloader_gui.manifest -ico ebookdownloader.ico -arch 386 -o rsrc_x32.syso

go build -ldflags "-H windowsgui -w -s -X main.Commit=$commit_id -X 'main.BuildTime=$build_time' -linkmode external -extldflags '-static'" -o ebookdownloader_gui.exe
Copy-Item ebookdownloader_gui.exe ..\
Remove-Item rsrc_x32.syso