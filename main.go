package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Secret-Ironman/boxr/pkg/api"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

var (
	host  = "localhost"
	port  = 3000
	db    = "boxr.db"
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func main() {
	app := cli.NewApp()
	app.Name = "boxr"
	app.Usage = "the boxr client app"
	app.Version = "0.0.1a"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  host,
			Usage:  "Hostname of the boxr api",
			EnvVar: "BOXR_API_HOST",
		},
		cli.IntFlag{
			Name:   "port, p",
			Value:  port,
			Usage:  "Port of the boxr api",
			EnvVar: "BOXR_API_PORT",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "pallet",
			Usage: "Manage Pallets",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					ShortName: "c",
					Usage:     "Create a new pallet in the wharehouse",

					Action: func(c *cli.Context) {

						name := c.Args().First()
						repo := c.Args()[1]
						// Create a new pallet.
						payload := &api.Pallet{
							Name: name,
							Url:  repo,
						}

						client := api.NewApiClient(c.GlobalString("host"), c.GlobalInt("port"), false)
						resp, err := client.CreatePallet(payload)

						if err != nil || !resp.Success {
							log.Fatal(err)
							fmt.Printf("%s Failed to create Pallet...\n", red("✘"))
							return
						}

						message, _ := resp.Message.(map[string]interface{})
						fmt.Printf("%s Succussfully created \"%s\" Pallet in %s\n", green("✔︎"), message["name"], resp.Took)
					},
				}, {
					Name:      "list",
					ShortName: "l",
					Usage:     "List all Pallets in the wharehouse",
					Action: func(c *cli.Context) {
						client := api.NewApiClient(c.GlobalString("host"), c.GlobalInt("port"), false)
						resp, err := client.GetAllPallets()

						if err != nil {
							fmt.Printf("%s Failed to fetch Pallets\n", red("✘"))
							return
						}

						message, _ := resp.Message.([]interface{})
						for _, value := range message {
							pallet := value.(map[string]interface{})
							fmt.Printf(" %s: %s (%s)\n", pallet["name"], pallet["url"], pallet["status"])
						}
					},
				},
			},
		}, {
			Name:      "serve",
			ShortName: "s",
			Usage:     "Start a boxr web service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "db",
					Value:  db,
					Usage:  "Database for the boxr api",
					EnvVar: "BOXR_DB",
				},
			},
			Action: func(c *cli.Context) {
				a, err := api.NewApi(db, 3000)
				if err != nil {
					log.Fatal(err)
				}

				a.Run()
			},
		},
	}

	app.Run(os.Args)
}
