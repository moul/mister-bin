package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const mapMaxSize = 1e4

type Binary struct {
	Asset []byte
	Name  string
}

func NewBinary(name string) Binary {
	return Binary{
		Name: name,
	}
}

func (b *Binary) FileName() string {
	return strings.Replace(b.Name, "/", "_", -1)
}

func (b *Binary) FilePath() string {
	return fmt.Sprintf("/tmp/mb-%s", b.FileName())
}

func (b *Binary) Uninstall() error {
	logrus.Infof("Uninstalling asset: %s", b.Name)

	if _, err := os.Stat(b.FilePath()); os.IsNotExist(err) {
		logrus.Warnf("Asset %q not installed, nothing to do", b.FilePath())
		return nil
	}

	return os.Remove(b.FilePath())
}

func (b *Binary) Install() error {
	logrus.Infof("Installing asset: %q", b.FilePath())

	if _, err := os.Stat(b.FilePath()); !os.IsNotExist(err) {
		logrus.Warnf("Asset already installed: %q", b.FilePath())
		return nil
	}

	asset, err := Asset(b.Name)
	if err != nil {
		return fmt.Errorf("failed to load the asset %q: %v", b.Name, err)
	}

	length := len(asset)
	size := int(unsafe.Sizeof(0)) * length
	if size > mapMaxSize*int(unsafe.Sizeof(0)) {
		return fmt.Errorf("file too big for current map size: %d > %d", size, mapMaxSize*int(unsafe.Sizeof(0)))
	}

	logrus.Infof("Creating map file: %q", b.FilePath())
	mapFile, err := os.Create(b.FilePath())
	if err != nil {
		return fmt.Errorf("failed to create map file: %v", err)
	}

	logrus.Infof("Seeking file")
	if _, err = mapFile.Seek(int64(length-1), 0); err != nil {
		return fmt.Errorf("failed to seek: %v", err)
	}

	logrus.Infof("Writing to file")
	if _, err = mapFile.Write([]byte(" ")); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	logrus.Infof("MMAPing")
	fd := int(mapFile.Fd())
	mmap, err := syscall.Mmap(fd, 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("failed to mmap: %v", err)
	}

	logrus.Infof("Filling array var")
	mapArray := (*[mapMaxSize]byte)(unsafe.Pointer(&mmap[0]))
	for i := 0; i < length; i++ {
		mapArray[i] = asset[i]
	}
	// fmt.Println(*mapArray)

	logrus.Infof("MUNMAPing")
	if err = syscall.Munmap(mmap); err != nil {
		return fmt.Errorf("failed to munmap: %v", err)
	}

	logrus.Infof("Closing")
	if err = mapFile.Close(); err != nil {
		return fmt.Errorf("failed to close: %v", err)
	}

	logrus.Infof("Chmoding binary")
	if err = os.Chmod(b.FilePath(), 0777); err != nil {
		return fmt.Errorf("failed to chmod program: %v", err)
	}

	return nil
}

func (b *Binary) Execute() error {
	logrus.Infof("Executing binary")
	cmd := exec.Command(b.FilePath())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute program: %v", err)
	}
	logrus.Infof("Output: %s", output)
	return nil
}

func GetBinaryByName(name string) (*Binary, error) {
	for _, assetName := range AssetNames() {
		if name == assetName {
			bin := NewBinary(assetName)
			return &bin, nil
		}
	}
	for _, assetName := range AssetNames() {
		if name == filepath.Base(assetName) {
			bin := NewBinary(assetName)
			return &bin, nil
		}
	}
	return nil, fmt.Errorf("No match")
}

func ActionExecute(c *cli.Context) {
	bin, err := GetBinaryByName(c.Command.Name)
	if err != nil {
		logrus.Fatalf("No such binary %q: %v", c.Command.Name, err)
	}
	if err := bin.Execute(); err != nil {
		logrus.Fatalf("Failed to execute binary: %v", err)
	}
}

func ActionInstall(c *cli.Context) {
	for _, name := range AssetNames() {
		bin := NewBinary(name)
		if err := bin.Install(); err != nil {
			logrus.Fatalf("Failed to install binary: %v", err)
		}
	}
}
func ActionUninstall(c *cli.Context) {
	for _, name := range AssetNames() {
		bin := NewBinary(name)
		if err := bin.Uninstall(); err != nil {
			logrus.Fatalf("Failed to uninstall binary: %v", err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Mister Bin"

	app.Commands = []cli.Command{
		{
			Name:   "install",
			Action: ActionInstall,
		},
		{
			Name:   "uninstall",
			Action: ActionUninstall,
		},
	}

	for _, name := range AssetNames() {
		command := cli.Command{
			Name:   filepath.Base(name),
			Action: ActionExecute,
		}
		app.Commands = append(app.Commands, command)
	}

	app.Run(os.Args)
}
