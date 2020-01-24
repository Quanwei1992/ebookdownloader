for /F %%i in ('git rev-parse --short HEAD') do ( set commitid=%%i)
echo commitid=%commitid%

go build -ldflags "-X main.Commit=%commitid%" -o ebookdownloader_cli.exe
copy ebookdownloader_cli.exe ..\
