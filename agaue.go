package main

import (
	"github.com/agaue/agaue/lib"
	"github.com/codegangsta/cli"
	"os"
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
			ShortName: "init",
			Usage:     "init a blog",
			Action: func(c *cli.Context) {
				lib.CreateBlog(c.Args()[0])
			},
		},
		{
			Name:      "generate",
			ShortName: "gen",
			Usage:     "generate blog",
			Action: func(c *cli.Context) {
				lib.GenerateSite()
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
