mkdir -p ./build/amd64
mkdir -p ./build/i386

export GOOS=linux
export GOARCH=amd64
go build
mv ./boi ./build/amd64/boi
export GOARCH=386
go build
mv ./boi ./build/i386/boi

export GOOS=windows
export GOARCH=amd64
go build
mv ./boi.exe ./build/amd64/boi.exe
export GOARCH=386
go build
mv ./boi.exe ./build/i386/boi.exe

