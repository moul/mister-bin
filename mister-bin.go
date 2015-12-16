package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Sirupsen/logrus"
)

func main() {
	for _, asset := range AssetNames() {
		mmapAsset(asset)
	}
}

func mmapAsset(assetName string) {
	logrus.Infof("Exporting asset: %s", assetName)
	asset, err := Asset(assetName)
	if err != nil {
		logrus.Fatalf("Failed to load the asset %q: %v", assetName, err)
	}

	const mapMaxSize = 1e4
	length := len(asset)
	size := int(unsafe.Sizeof(0)) * mapMaxSize

	filename := strings.Replace(assetName, "/", "_", -1)
	logrus.Infof("Creating map file: '/tmp/mb-%s'", filename)
	mapFile, err := os.Create(fmt.Sprintf("/tmp/mb-%s", filename))
	if err != nil {
		logrus.Fatalf("Failed to create map file: %v", err)
	}

	logrus.Infof("Seeking file")
	if _, err = mapFile.Seek(int64(size-1), 0); err != nil {
		logrus.Fatalf("Failed to seek: %v", err)
	}

	logrus.Infof("Writing to file")
	if _, err = mapFile.Write([]byte(" ")); err != nil {
		logrus.Fatalf("Failed to write to file: %v", err)
	}

	logrus.Infof("MMAPing")
	fd := int(mapFile.Fd())
	mmap, err := syscall.Mmap(fd, 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		logrus.Fatalf("Failed to mmap: %v", err)
	}

	logrus.Infof("Filling array var")
	mapArray := (*[mapMaxSize]byte)(unsafe.Pointer(&mmap[0]))
	for i := 0; i < length; i++ {
		mapArray[i] = asset[i]
	}
	// fmt.Println(*mapArray)

	logrus.Infof("MUNMAPing")
	if err = syscall.Munmap(mmap); err != nil {
		logrus.Fatalf("Failed to munmap: %v", err)
	}

	logrus.Infof("Closing")
	if err = mapFile.Close(); err != nil {
		logrus.Fatalf("Failed to close: %v", err)
	}
}
