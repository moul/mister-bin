# mister-bin
:sparkles: Binaries + Steroids = Abracadabra !

An attempt to address [moul/random-ideas#5](https://github.com/moul/random-ideas/issues/5) and [moul/random-ideas#6](https://github.com/moul/random-ideas/issues/6).

## Test

```console
➜  mister-bin git:(master) ✗ make test
go get github.com/jteeuwen/go-bindata/...
go-bindata ./test/linux-x86_64-helloworld-static
go build -o mister-bin .
./mister-bin -h || true
NAME:
   Mister Bin - A new cli application

USAGE:
   Mister Bin [global options] command [command options] [arguments...]

COMMANDS:
   install
   uninstall
   linux-x86_64-helloworld-static
   help, h                             Shows a list of commands or help for one command

./mister-bin ./test/linux-x86_64-helloworld-static
No help topic for './test/linux-x86_64-helloworld-static'
```

## Docker

```console
$ make docker
docker build --no-cache -t mister-bin docker
Sending build context to Docker daemon     4 MB
Step 1 : FROM scratch
--->
Step 2 : ADD ./mister-bin /bin/sh
---> e51ce5af547e
Removing intermediate container 08ebf5bbd60a
Step 3 : RUN /bin/sh install --basedir=/bin --symlinks
---> Running in 9943e91ccd0a
---> cb0f90aee30d
Removing intermediate container 9943e91ccd0a
Successfully built cb0f90aee30d
docker run -it --rm mister-bin /bin/linux-x86_64-helloworld-static
Hello World !
docker export `docker create mister-bin /dont-exists` | tar -tvf -
-rwxr-xr-x  0 0      0           0 Dec 28 17:56 .dockerenv
-rwxr-xr-x  0 0      0           0 Dec 28 17:56 .dockerinit
drwxr-xr-x  0 0      0           0 Dec 28 17:56 bin/
lrwxrwxrwx  0 0      0           0 Dec 28 17:56 bin/linux-x86_64-helloworld-static -> /bin/sh
-rwxr-xr-x  0 0      0     3993656 Dec 28 17:56 bin/sh
drwxr-xr-x  0 0      0           0 Dec 28 17:56 dev/
-rwxr-xr-x  0 0      0           0 Dec 28 17:56 dev/console
drwxr-xr-x  0 0      0           0 Dec 28 17:56 dev/pts/
drwxr-xr-x  0 0      0           0 Dec 28 17:56 dev/shm/
drwxr-xr-x  0 0      0           0 Dec 28 17:56 etc/
-rwxr-xr-x  0 0      0           0 Dec 28 17:56 etc/hostname
-rwxr-xr-x  0 0      0           0 Dec 28 17:56 etc/hosts
lrwxrwxrwx  0 0      0           0 Dec 28 17:56 etc/mtab -> /proc/mounts
-rwxr-xr-x  0 0      0           0 Dec 28 17:56 etc/resolv.conf
drwxr-xr-x  0 0      0           0 Dec 28 17:56 proc/
drwxr-xr-x  0 0      0           0 Dec 28 17:56 sys/
```

## License

MIT
