build:
	mkdir build/
	env GOOS=windows GOARCH=amd64 go build -o build/semver-windows-amd64-$(VERSION).exe main.go
	env GOOS=darwin  GOARCH=amd64 go build -o build/semver-darwin-amd64-$(VERSION) main.go
	env GOOS=linux   GOARCH=amd64 go build -o build/semver-linux-amd64-$(VERSION) main.go

clean:
	rm -rf build/

release-tag:
	git tag v$(go run main.go) -m "Release version $(go run main.go)"