echo killing cloudwebsite
taskkill /IM cloudwebsite.exe /T /F

set GOPATH=%~dp0
set GOBIN=.\src\cloudwebsite\
REM set GOOS=windows

go install -tags debug cloudwebsite 

