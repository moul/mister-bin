all: helloworld-dynamic helloworld-static quine-static quine-dynamic


helloworld-dynamic: helloworld.c
	docker run --rm -v $(PWD):/workdir multiarch/crossbuild cc $< -o $@


helloworld-static: helloworld.c
	docker run --rm -v $(PWD):/workdir multiarch/crossbuild cc $< -o $@ -static


quine-dynamic: quine.c
	docker run --rm -v $(PWD):/workdir multiarch/crossbuild cc $< -o $@


quine-static: quine.c
	docker run --rm -v $(PWD):/workdir multiarch/crossbuild cc $< -o $@ -static
