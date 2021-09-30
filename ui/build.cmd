for /F %%i in ('git rev-parse --short HEAD') do ( set commitid=%%i)
echo commitid=%commitid%

go build -ldflags "-w -s -H windowsgui -X main.Commit=%commitid%" -o ebookdl_ui.exe
copy ebookdl_ui.exe ..\
