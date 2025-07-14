@echo off
echo Starting TopV Adaptor Go...

echo.
echo Building project...
go mod tidy
go build -o topv-adaptor.exe main.go push.go

echo.
echo Starting application...
topv-adaptor.exe

pause 