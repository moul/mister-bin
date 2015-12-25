#BINARY = ./test/darwin-x86_64-helloworld-dynamic
BINARY = ./test/linux-x86_64-helloworld-static


.PHONY: build
build: mister-bin


$(BINARY):
	cd test; $(MAKE)


bindata.go: $(BINARY)
	go get github.com/jteeuwen/go-bindata/...
	go-bindata $(BINARY)


mister-bin: mister-bin.go bindata.go
	go build -o mister-bin .


.PHONY: test
test: mister-bin
	./mister-bin -h || true
	./mister-bin $(BINARY)


.PHONY: docker
docker: docker/mister-bin docker/Dockerfile
	docker build --no-cache -t mister-bin docker
	docker run -it --rm mister-bin /bin/$(notdir $(BINARY))


docker/mister-bin: mister-bin.go bindata.go
	goxc -bc="linux,386" -d=docker -o="{{.Dest}}{{.PS}}{{.ExeName}}{{.Ext}}" -include="" compile


.PHONY: clean
clean:
	./mister-bin --uninstall
