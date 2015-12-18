build: mister-bin

./test/darwin-x86_64-helloworld-dynamic:
	cd test; $(MAKE)

bindata.go: ./test/darwin-x86_64-helloworld-dynamic
	go get github.com/jteeuwen/go-bindata/...
	go-bindata ./test/darwin-x86_64-helloworld-dynamic

mister-bin: mister-bin.go bindata.go
	go build -o mister-bin .

test: mister-bin
	./mister-bin
	ls -la /tmp/mb-test_darwin-x86_64-helloworld-dynamic test/darwin-x86_64-helloworld-dynamic
	/tmp/mb-test_darwin-x86_64-helloworld-dynamic
	./test/darwin-x86_64-helloworld-dynamic

clean:
	rm -f /tmp/mb-test_darwin-x86_64-helloworld-dynamic
