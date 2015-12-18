all: build


build:
	docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE=x86_64 multiarch/crossbuild cc helloworld.c -o linux-x86_64-helloworld-dynamic           || true
	docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE=x86_64 multiarch/crossbuild cc helloworld.c -o linux-x86_64-helloworld-static -static    || true
	docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE=x86_64 multiarch/crossbuild cc quine.c      -o linux-x86_64-quine-dynamic                || true
	docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE=x86_64 multiarch/crossbuild cc quine.c      -o linux-x86_64-quine-static -static         || true
	docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE=darwin multiarch/crossbuild cc helloworld.c -w -o darwin-x86_64-helloworld-dynamic          || true
	#docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE=darwin multiarch/crossbuild cc helloworld.c -o darwin-x86_64-helloworld-static -static   || true
	docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE-darwin multiarch/crossbuild cc quine.c      -w -o darwin-x86_64-quine-dynamic               || true
	#docker run --rm -v $(PWD):/workdir -e CROSS_TRIPLE=darwin multiarch/crossbuild cc quine.c      -o darwin-x86_64-quine-static -static        || true
