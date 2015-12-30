#BINARY = ./test/darwin-x86_64-helloworld-dynamic
BINARY = ./test/linux-x86_64-helloworld-static
MISTERBIN_DIR = ./misterbin


.PHONY: build
build: mister-bin


$(BINARY):
	cd test; $(MAKE)


$(MISTERBIN_DIR)/bindata.go: $(BINARY)
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -o ./$@ $(BINARY)


mister-bin: $(MISTERBIN_DIR)/mister-bin.go $(MISTERBIN_DIR)/bindata.go
	go build -o ./mister-bin ./$(MISTERBIN_DIR)


.PHONY: test
test: mister-bin
	./mister-bin -h || true
	./mister-bin $(notdir $(BINARY))


.PHONY: docker
docker: docker/mister-bin docker/Dockerfile
	docker build --no-cache -t mister-bin docker
	docker run -it --rm mister-bin /bin/$(notdir $(BINARY))
	docker export `docker create mister-bin /dont-exists` | tar -tvf -


docker/mister-bin: $(MISTERBIN_DIR)/mister-bin.go $(MISTERBIN_DIR)/bindata.go
	cd $(MISTERBIN_DIR); goxc -bc="linux,386" -d=../docker -n=mister-bin -o="{{.Dest}}{{.PS}}{{.ExeName}}{{.Ext}}" -include="" compile


.PHONY: clean
clean:
	./mister-bin --uninstall
