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
	Asset   []byte
	Name    string
	BaseDir string
}

func NewBinary(name string) Binary {
	wd, err := os.Getwd()
	if err != nil {
		wd = "/"
	}
	return Binary{
		Name:    name,
		BaseDir: wd,
	}
}

func (b *Binary) FileName() string {
	return strings.Replace(b.Name, "/", "_", -1)
}

func (b *Binary) FilePath() string {
	return filepath.Join(b.BaseDir, b.FileName())
}

func (b *Binary) TempPath() string {
	// FIXME: ensure using memory instead of disk if p ossible
	return filepath.Join(os.TempDir(), b.FileName())
}

func (b *Binary) Uninstall(filepath string) error {
	logrus.Infof("Uninstalling binary %q (%s)", b.Name, filepath)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		logrus.Debugf("Asset %q not installed, nothing to do", filepath)
		return nil
	}

	return os.Remove(filepath)
}

func (b *Binary) Install(filepath string) error {
	logrus.Infof("Installing binary %q (%s)", b.Name, filepath)

	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		logrus.Debugf("Asset already installed: %q", filepath)
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

	logrus.Debugf("Creating map file: %q", filepath)
	mapFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create map file: %v", err)
	}

	logrus.Debugf("Seeking file")
	if _, err = mapFile.Seek(int64(length-1), 0); err != nil {
		return fmt.Errorf("failed to seek: %v", err)
	}

	logrus.Debugf("Writing to file")
	if _, err = mapFile.Write([]byte(" ")); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	logrus.Debugf("MMAPing")
	fd := int(mapFile.Fd())
	mmap, err := syscall.Mmap(fd, 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("failed to mmap: %v", err)
	}

	logrus.Debugf("Filling array var")
	mapArray := (*[mapMaxSize]byte)(unsafe.Pointer(&mmap[0]))
	for i := 0; i < length; i++ {
		mapArray[i] = asset[i]
	}
	// fmt.Println(*mapArray)

	logrus.Debugf("MUNMAPing")
	if err = syscall.Munmap(mmap); err != nil {
		return fmt.Errorf("failed to munmap: %v", err)
	}

	logrus.Debugf("Closing")
	if err = mapFile.Close(); err != nil {
		return fmt.Errorf("failed to close: %v", err)
	}

	logrus.Debugf("Chmoding binary")
	if err = os.Chmod(filepath, 0777); err != nil {
		return fmt.Errorf("failed to chmod program: %v", err)
	}

	return nil
}

func (b *Binary) Execute(args []string) error {
	filepath := b.TempPath()
	logrus.Infof("Executing binary %q (%s)", b.Name, filepath)
	if len(args) > 0 {
		logrus.Infof("args=%v", args)
	}

	// temporary
	if err := b.Uninstall(filepath); err != nil {
		logrus.Debugf("Failed to uninstall temporary binary %q: %v", filepath, err)
	}

	if err := b.Install(filepath); err != nil {
		return err
	}
	// FIXME: delete early
	defer b.Uninstall(filepath)

	cmd := exec.Command(filepath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()

	return cmd.Wait()
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
	binaryUseContext(&bin, c)

	if err := bin.Execute(c.Args()); err != nil {
		logrus.Fatalf("Failed to execute binary: %v", err)
	}
}

func binaryUseContext(bin *Binary, context *cli.Context) {
	if basedir := context.Parent().String("basedir"); basedir != "" {
		bin.BaseDir = basedir
	}
}

func ActionInstall(c *cli.Context) {
	for _, name := range AssetNames() {
		bin := NewBinary(name)
		binaryUseContext(&bin, c)

		if err := bin.Install(bin.FilePath()); err != nil {
			logrus.Fatalf("Failed to install binary: %v", err)
		}
	}
}
func ActionUninstall(c *cli.Context) {
	for _, name := range AssetNames() {
		bin := NewBinary(name)
		binaryUseContext(&bin, c)

		if err := bin.Uninstall(bin.FilePath()); err != nil {
			logrus.Fatalf("Failed to uninstall binary: %v", err)
		}
	}
}

func hookBefore(c *cli.Context) error {
	if c.Bool("debug") {
		os.Setenv("MB_DEBUG", "1")
	}
	initLogging()
	return nil
}

func initLogging() {
	if os.Getenv("MB_DEBUG") == "1" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func main() {
	// Checking if file is a symlink to mister-bin or mister-bin itself
	fi, err := os.Lstat(os.Args[0])
	if err != nil {
		logrus.Fatalf("Failed to lstat %q: %v", os.Args[0], err)
	}

	// we are a symlink, direct execution
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		initLogging()
		bin, err := GetBinaryByName(filepath.Base(os.Args[0]))
		if err != nil {
			logrus.Fatalf("No such binary %q: %v", filepath.Base(os.Args[0]), err)
		}

		if err := bin.Execute(os.Args[1:]); err != nil {
			logrus.Fatalf("Failed to execute binary: %v", err)
		}
		return
	}

	// we are not a symlink, standard menu
	app := cli.NewApp()
	app.Name = "Mister Bin"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, D",
			EnvVar: "MB_DEBUG",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "basedir",
		},
	}

	app.Before = hookBefore

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
