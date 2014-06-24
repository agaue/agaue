package main

import (
	"fmt"
	"github.com/agaue/agaue/lib"
	"github.com/codegangsta/cli"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "agaue"
	app.Usage = "make a test"
	app.Version = "0.0.1"
	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{"lang, l", "english", "language for the greeting"},
	// }
	// app.Action = func(c *cli.Context) {
	// 	name := "someone"
	// 	if len(c.Args()) > 0 {
	// 		name = c.Args()[0]
	// 	}
	// 	if c.String("lang") == "spanish" {
	// 		println("Hola", name)
	// 	} else {
	// 		println("Hello", name)
	// 	}
	// }
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
