package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
	
	"github.com/agaue/agaue/lib"
	"github.com/urfave/cli"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "agaue"
	app.Usage = "Agaue -- static blog engine"
	app.Version = "0.0.3"

	app.Commands = []cli.Command{
		{
			Name:      "init",
			ShortName: "i",
			Usage:     "init a blog",
			Action: func(c *cli.Context) {
				lib.CreateBlog(c.Args()[0])
			},
		},
		{
			Name:      "generate",
			ShortName: "g",
			Usage:     "generate blog",
			Action: func(c *cli.Context) {
				start := time.Now()
				lib.GenerateSite()
				fmt.Println("Completed, Spend : ", time.Since(start))
			},
		},
		{
			Name:      "server",
			ShortName: "s",
			Usage:     "server blog and auto restart when files change",
			Action: func(c *cli.Context) {
				lib.ServeBlog()
			},
		},
		{
			Name:      "deploy",
			ShortName: "d",
			Usage:     "deploy blog to github",
			Action: func(c *cli.Context) {
				lib.DeploySite()
			},
		},
	}

	app.Run(os.Args)
}
