package main

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"github.com/codegangsta/cli"

	"github.com/Reisender/acdfuse/acdfs"
)

func main() {
	app := cli.NewApp()
	app.Name = "acdfuse"
	app.Usage = "mount amazon cloud drive as fuse file system"
	app.Action = Run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mountpoint, m",
			Usage: "the directory to mount acdfuse fs to",
		},
	}

	app.Run(os.Args)
}

func Run(c *cli.Context) {
	println("Amazon Cloud Drive mount at", c.String("mountpoint"))
	mountpoint := c.String("mountpoint")
	if mountpoint == "" {
		log.Fatal("no mountpoint! try running \"acdfuse help\"")
	}

	fuseCtx, err := fuse.Mount(
		c.String("mountpoint"),
		fuse.FSName("helloworld"),
		fuse.Subtype("hellofs"),
		fuse.LocalVolume(),
		fuse.VolumeName("Hello world!"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer fuseCtx.Close()

	err = fs.Serve(fuseCtx, acdfs.FS{})
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-fuseCtx.Ready
	if err := fuseCtx.MountError; err != nil {
		log.Fatal(err)
	}
}
