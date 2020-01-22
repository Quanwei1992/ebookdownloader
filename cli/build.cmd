for /F %%i in ('git rev-parse --short HEAD') do ( set commitid=%%i)
echo commitid=%commitid%

set CURRENT_DATE=%date:~0,4%-%date:~5,2%-%date:~8,2%
set CURRENT_TIME=%time:~0,2%:%time:~3,2%:%time:~6,2%
echo %CURRENT_DATE% %CURRENT_TIME%
set buildtime=%CURRENT_DATE%-%CURRENT_TIME%
go build -ldflags "-X main.Commit=%commitid% -X main.BuildTime=%buildtime%" -o ebookdownloader_cli.exe
copy ebookdownloader_cli.exe ..\
