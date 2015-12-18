build: mister-bin

mister-bin: mister-bin.go bindata.go
	go build -o mister-bin .

test: mister-bin
	go-bindata ./test/darwin-x86_64-helloworld-dynamic
	./mister-bin
	ls -la /tmp/mb-test_darwin-x86_64-helloworld-dynamic test/darwin-x86_64-helloworld-dynamic
	/tmp/mb-test_darwin-x86_64-helloworld-dynamic
	./test/darwin-x86_64-helloworld-dynamic

clean:
	rm -f /tmp/mb-test_darwin-x86_64-helloworld-dynamic
