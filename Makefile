#BINARY = ./test/darwin-x86_64-helloworld-dynamic
BINARY = ./test/linux-x86_64-helloworld-static
MISTERBIN_DIR = ./cmd/misterbin
MISTERBUILD_DIR = ./cmd/misterbuild


.PHONY: build
build: misterbin misterbuild


misterbuild: $(MISTERBUILD_DIR)/mister-build.go $(MISTERBUILD_DIR)/bindata.go
	go build -o ./$@ ./$(MISTERBUILD_DIR)


$(MISTERBUILD_DIR)/bindata.go: $(MISTERBIN_DIR)/mister-bin.go
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -prefix=$(dir $<) -nocompress -o ./$@ $<
	ls -la $@


$(BINARY):
	cd test; $(MAKE)


$(MISTERBIN_DIR)/bindata.go: $(BINARY)
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -o ./$@ $(BINARY)
	ls -la $@


misterbin: $(MISTERBIN_DIR)/mister-bin.go $(MISTERBIN_DIR)/bindata.go
	go build -o ./$@ ./$(MISTERBIN_DIR)


.PHONY: test
test: misterbin
	./misterbin -h || true
	./misterbin $(notdir $(BINARY))


.PHONY: docker
docker: docker/misterbin docker/Dockerfile
	docker build --no-cache -t misterbin docker
	docker run -it --rm misterbin /bin/$(notdir $(BINARY))
	docker export `docker create misterbin /dont-exists` | tar -tvf -


docker/misterbin: $(MISTERBIN_DIR)/mister-bin.go $(MISTERBIN_DIR)/bindata.go
	cd $(MISTERBIN_DIR); goxc -bc="linux,386" -d=../docker -n=misterbin -o="{{.Dest}}{{.PS}}{{.ExeName}}{{.Ext}}" -include="" compile


.PHONY: clean
clean:
	./misterbin --uninstall
