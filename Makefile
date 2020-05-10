.PHONY: build

build:
	go build -o bin/psioHelperLinux main.go && GOOS=windows GOARCH=386 go build -o bin/psioHelper.exe main.go && GOOS=darwin GOARCH=amd64 go build -o bin/psioHelperMac main.go
