GOFLAGS = -ldflags="-extldflags=-static"
GOENV = CGO_ENABLED=0

bootstrap: bootstrap-linux bootstrap-windows bootstrap-darwin

bootstrap-linux:
	$(GOENV) GOOS=linux GOARCH=amd64 go build -o bin/builder.x86_64 $(GOFLAGS) ./builder
	$(GOENV) GOOS=linux GOARCH=386 go build -o bin/builder.i386 $(GOFLAGS) ./builder
	$(GOENV) GOOS=linux GOARCH=arm64 go build -o bin/builder.aarch64 $(GOFLAGS) ./builder
	$(GOENV) GOOS=linux GOARCH=arm go build -o bin/builder.arm $(GOFLAGS) ./builder

bootstrap-windows:
	$(GOENV) GOOS=windows GOARCH=amd64 go build -o bin/builder.x86_64.exe $(GOFLAGS) ./builder
	$(GOENV) GOOS=windows GOARCH=386 go build -o bin/builder.i386.exe $(GOFLAGS) ./builder
	$(GOENV) GOOS=windows GOARCH=arm64 go build -o bin/builder.arm64.exe $(GOFLAGS) ./builder
	$(GOENV) GOOS=windows GOARCH=arm go build -o bin/builder.arm.exe $(GOFLAGS) ./builder

bootstrap-darwin:
	$(GOENV) GOOS=darwin GOARCH=amd64 go build -o bin/builder+mach_o.x86_64 $(GOFLAGS) ./builder
	$(GOENV) GOOS=darwin GOARCH=arm64 go build -o bin/builder+mach_o.arm64 $(GOFLAGS) ./builder

	chmod +x bin/lipo.$(shell uname -m)
	bin/lipo.$(shell uname -m) -create -output bin/builder.mach_o bin/builder+mach_o.x86_64 bin/builder+mach_o.arm64

clean:
	rm -f bin/builder.x86_64 bin/builder.i386 bin/builder.aarch64 bin/builder.arm
	rm -f bin/builder.x86_64.exe bin/builder.i386.exe bin/builder.arm64.exe bin/builder.arm.exe
	rm -f bin/builder+mach_o.x86_64 bin/builder+mach_o.arm64 bin/builder.mach_o

.PHONY: bootstrap
