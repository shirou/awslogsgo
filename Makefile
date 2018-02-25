VERSION=$(shell git describe --tags --exact-match)

init_tool:
	go get github.com/Songmu/goxz/cmd/goxz
	go get github.com/tcnksm/ghr

build_all:
	goxz -pv=$(VERSION) -os=darwin,linux,windows,freebsd -arch=amd64 -d=dist .

release: build_all
	ghr $(VERSION) ./dist
