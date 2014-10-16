package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/Secret-Ironman/boxr/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/Secret-Ironman/boxr/shared/api"
	"github.com/Secret-Ironman/boxr/shared/types"
)

var (
	host = "localhost"
	port = 3000
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
			Name:      "add",
			ShortName: "a",
			Usage:     "Add a new pallet to the wharehouse",

			Action: func(c *cli.Context) {
				name := c.Args().First()
				repo := c.Args()[1]
				// post repo link to api and enqueue
				// it in forklift creation.
				payload := &types.Pallet{}
				payload.Name = name
				payload.Url = repo

				client := api.NewApiClient(c.GlobalString("host"), c.GlobalInt("port"), false)
				resp, err := client.CreatePallet(payload)

				if err != nil {
					log.Fatal(err)
				}

				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				var response api.Response
				e := json.Unmarshal(body, &response)
				if e != nil {
					log.Fatal(err)
				}
				log.Printf("%+v", response)
			},
		},
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "Start a boxr web service",
			Action: func(c *cli.Context) {
				api.Run(":3000")
			},
		},
	}

	app.Run(os.Args)
}
