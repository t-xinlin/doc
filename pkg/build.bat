SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o server

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build

export CGO_ENABLED=0
export GOOS=linux
go build -a -installsuffix cgo -o app .