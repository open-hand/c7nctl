.PHONY: clean

nothing:
	echo nothing
tmpdir:
	mkdir c7n-${VERSION}
darwin: 
	GOOS=darwin GOARCH=amd64 go build -o c7n-darwin
linux: 
	GOOS=linux GOARCH=amd64 go build  -o c7n-linux
publish: tmpdir darwin linux
	mv c7n-darwin c7n-${VERSION}/c7n
	tar -czf c7n-${VERSION}-Darwin-amd64.tar.gz c7n-${VERSION}
	mv c7n-linux c7n-${VERSION}/c7n
	tar -czf c7n-${VERSION}-Linux-amd64.tar.gz c7n-${VERSION}
clean:
	rm -rf c7n-${VERSION}