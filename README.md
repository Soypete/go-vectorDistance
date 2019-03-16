# go-vectorDistance
 this repo tests methods for measuring vector distances. 

This repo used the library (https://github.com/xrash/smetrics) to calculate vector distances.

## Configuration
 ` go get https://github.com/xrash/smetrics`


  `go run main.go`

## Dependency Example

## Dep

```bash
brew install dep
dep init
dep ensure
go run main.go
```

## Go Modules
```bash
export GO111MODULES="on"
go mod init
go get https://github.com/xrash/smetrics
go build ./...
./go-vectorDistance
```
