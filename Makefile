
all: livestatusd

livestatusd: lsd/lsd
	cp $< $@

lsd/lsd: lsd/*.go
	cd lsd; go fmt; go build

lsd/parser.go: lsd/parser.y
	cd lsd; go generate

test: lsd/parser.go
	cd lsd; go test

deps:
	go get github.com/Sirupsen/logrus
