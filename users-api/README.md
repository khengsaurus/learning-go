NOTE: receiving status 502 from AWS API url

build command: 
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" main.go
<!-- GOARCH=amd64: x86-64 -->

zip:
zip build/main.zip main