all:	install

install:
	go install

test:
	go test -v ./...

clean:
	go clean ./...

nuke:
	go clean -i ./...
