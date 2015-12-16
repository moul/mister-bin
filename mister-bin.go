package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Sirupsen/logrus"
)

const mapMaxSize = 1e4

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

	length := len(asset)
	size := int(unsafe.Sizeof(0)) * length
	if size > mapMaxSize*int(unsafe.Sizeof(0)) {
		logrus.Fatalf("File too big for current map size: %d > %d", size, mapMaxSize*int(unsafe.Sizeof(0)))
	}

	filename := strings.Replace(assetName, "/", "_", -1)
	filepath := fmt.Sprintf("/tmp/mb-%s", filename)
	logrus.Infof("Creating map file: %q", filepath)
	mapFile, err := os.Create(filepath)
	if err != nil {
		logrus.Fatalf("Failed to create map file: %v", err)
	}

	logrus.Infof("Seeking file")
	if _, err = mapFile.Seek(int64(length-1), 0); err != nil {
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

	logrus.Infof("Executing binary")
	cmd := exec.Command(filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Fatalf("Failed to execute program: %v", err)
	}
	logrus.Infof("Output: %s", output)
}
