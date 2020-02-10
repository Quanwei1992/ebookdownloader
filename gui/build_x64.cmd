set commitid=""
for /F %%i in ('git rev-parse --short HEAD') do ( set commitid=%%i)
echo commitid=%commitid%

rsrc -manifest ebookdownloader_gui.manifest -ico ebookdownloader.ico -arch amd64 -o rsrc_x64.syso

go build -ldflags "-H windowsgui -w -s -X main.Commit=%commitid% -linkmode external -extldflags '-static'" -o ebookdownloader_gui.exe
copy ebookdownloader_gui.exe ..\
del  rsrc_x64.syso
