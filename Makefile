.PHONY: clean

nothing:
	echo nothing
tmpdir:
	mkdir c7nctl-${VERSION}-darwin
	mkdir c7nctl-${VERSION}-linux
	mkdir c7nctl-${VERSION}-windows
darwin: 
	GO111MODULE="on" GOOS=darwin GOARCH=amd64 go build -o c7nctl-darwin -mod vendor
linux: 
	GO111MODULE="on" GOOS=linux GOARCH=amd64 go build  -o c7nctl-linux -mod vendor
windows:
	GO111MODULE="on" GOOS=windows GOARCH=amd64 go build  -o c7nctl-windows.exe -mod vendor
publish: tmpdir darwin linux windows
	mv c7nctl-darwin c7nctl-${VERSION}-darwin/c7nctl
	tar -czf c7nctl-${VERSION}-Darwin-amd64.tar.gz c7nctl-${VERSION}-darwin

	mv c7nctl-linux c7nctl-${VERSION}-linux/c7nctl
	tar -czf c7nctl-${VERSION}-Linux-amd64.tar.gz c7nctl-${VERSION}-linux

	mv c7nctl-windows.exe c7nctl-${VERSION}-windows/c7nctl.exe
	tar -czf c7nctl-${VERSION}-Windows-amd64.tar.gz c7nctl-${VERSION}-windows
clean:
	rm -rf c7nctl-${VERSION}