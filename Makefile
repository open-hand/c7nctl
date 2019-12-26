.PHONY: clean

nothing:
	echo nothing
tmpdir:
	mkdir c7nctl-${VERSION}
darwin: 
	GO111MODULE="on" GOOS=darwin GOARCH=amd64 go build -o c7nctl-darwin -mod vendor
linux: 
	GO111MODULE="on" GOOS=linux GOARCH=amd64 go build  -o c7nctl-linux -mod vendor
windows:
	GO111MODULE="on" GOOS=windows GOARCH=amd64 go build  -o c7nctl-windows -mod vendor
publish: tmpdir darwin linux windows
	mv c7nctl-darwin c7nctl-${VERSION}/c7nctl
	tar -czf c7nctl-${VERSION}-Darwin-amd64.tar.gz c7nctl-${VERSION}
	mv c7nctl-linux c7nctl-${VERSION}/c7nctl
	tar -czf c7nctl-${VERSION}-Linux-amd64.tar.gz c7nctl-${VERSION}
	mv c7nctl-windows c7nctl-${VERSION}/c7nctl
	tar -czf c7nctl-${VERSION}-Windows-amd64.tar.gz c7nctl-${VERSION}
clean:
	rm -rf c7nctl-${VERSION}