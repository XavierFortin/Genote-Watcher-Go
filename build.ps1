# Build for Linux
$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
go build -C src -ldflags "-w -s" -o "../bin/genote-watcher"

# Build for Windows
$Env:GOOS = "windows"; $Env:GOARCH = "amd64"
go build -C src -ldflags "-w -s" -o "..\bin\genote-watcher.exe" 
