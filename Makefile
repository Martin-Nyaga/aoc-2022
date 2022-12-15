.PHONY: test coverage clean

test:
	go test -v ./...

coverage:
	go test -v -coverprofile tmp/cover.out ./...
	go tool cover -html tmp/cover.out -o tmp/cover.html
	google-chrome-stable tmp/cover.html
	
clean:
	rm -f tmp/cover.out tmp/cover.html
