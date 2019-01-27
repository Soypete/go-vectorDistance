# go-vectorDistance

Using the data provided in (https://git.nav.com/mpeterson/go-vectorDistance/blob/master/20180622/20180622_results_searches.csv), this repo tests methods for measuring vector distances. Some of the methods have been implemented by the DandB business search function in Pudge (https://git.nav.com/backend/pudge).

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
