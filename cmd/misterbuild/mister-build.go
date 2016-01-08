package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/jteeuwen/go-bindata"
)

func BuildAction(c *cli.Context) {
	if len(c.Args()) < 1 {
		logrus.Fatalf("You need to specify at least 1 binary")
	}

	tmpName := "misterbuild-tmp"
	workdir := filepath.Join(os.Getenv("GOPATH"), "src", "tmp", tmpName)
	// FIXME: use `ioutil.TempDir()` instead

	// Generating bindata with target binaries
	cfg := bindata.NewConfig()
	cfg.Output = filepath.Join(workdir, "bindata.go")
	// cfg.Prefix = goroot
	cfg.Input = []bindata.InputConfig{}
	for _, assetPath := range c.Args() {
		cfg.Input = append(cfg.Input, bindata.InputConfig{
			Path: assetPath,
		})
	}
	if err := bindata.Translate(cfg); err != nil {
		logrus.Fatalf("Failed to generate source based on binaries: %v", err)
	}

	// Exporting mister-bin.go
	if err := RestoreAsset(workdir, "mister-bin.go"); err != nil {
		logrus.Fatalf("Failed to restore mister-bin.go: %v", err)
	}

	// Build
	cmd := exec.Command("go", "build", filepath.Join("tmp", tmpName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Failed to build binary: %v", err)
	}
	// FIXME: try to use go/build package directly

	logrus.Infof("Success: binary built: %q", tmpName)

	// Cleanup
	// FIXME: todo
}

func main() {
	app := cli.NewApp()
	app.Name = "Mister Build"
	app.Usage = "mister-bin builder"

	app.Commands = []cli.Command{
		{
			Name:   "build",
			Action: BuildAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "basedir",
					Usage: "Base dir to install binaries to",
				},
				cli.BoolFlag{
					Name:  "symlinks",
					Usage: " Create symlinks instead of real binaries",
				},
			},
		},
	}

	app.Run(os.Args)
}
