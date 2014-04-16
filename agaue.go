package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/agaue/agaue/lib"
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
				os.Mkdir("test", os.ModePerm)
			},
		},
		{
			Name:      "generate",
			ShortName: "g",
			Usage:     "generate html from markdown",
			Action: func(c *cli.Context) {
				lib.
			},
		},
	}

	app.Run(os.Args)
}
