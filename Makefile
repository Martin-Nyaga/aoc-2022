.PHONY: test coverage clean

test:
	go test -v ./...

coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	google-chrome-stable cover.html
	
clean:
	rm -f cover.out cover.html
