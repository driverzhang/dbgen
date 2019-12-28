package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"

	"github.com/driverzhang/dbgen/tool/db-gen-mongo"
)

func main() {
	app := cli.NewApp()
	app.Name = "dbgen"
	app.Usage = "dbgen工具集"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:    "mongo",
			Aliases: []string{"mo"},
			Usage:   "create mongo db",
			// Flags: []cli.Flag{
			// 	cli.StringFlag{
			// 		Name:        "f",
			// 		Value:       "",
			// 		Usage:       "go file for create mongo db",
			// 		Destination: &db_gen_mongo.P.Path,
			// 	},
			// },
			Action: db_gen_mongo.Mongo2Crud,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "dbgen version",
			Action: func(c *cli.Context) error {
				fmt.Println(getVersion())
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
