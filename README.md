# goffold

Opinionated scaffolding for go.

This project is on hold, because I'm currently determining how my opinion is
about the package structure etc. If you're curious, have a look at [kickstart](https://github.com/fvosberg/kickstart).

# This is a working draft

## Installation

``` bash
go get -u github.com/fvosberg/goffold
```

## CLI

``` bash
goffold new go-service --name=github.com/fvosberg/slimdentity --docker=registry.docker.com/service path
goffold add rest
goffold add kafka --for= --pkg=
```

# TODOs

[ ] Add `new service`
[ ] Add `go mod init {{ .ServiceName }}`
[ ] Add `add gitlab-ci`
[ ] Add `add kubernetes`
[ ] Add templates to binary
[ ] Add `add rest`
[ ] Add `add kafka`
[ ] Support altering go version
[ ] Add license
[ ] Add `add circle-ci`
