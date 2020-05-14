.PHONY: build

build:
	go build -o build/psioHelperLinux main.go && GOOS=windows GOARCH=386 go build -o build/psioHelper.exe main.go && GOOS=darwin GOARCH=amd64 go build -o build/psioHelperMac main.go
