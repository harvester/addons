package main

import (
	"fmt"
	"github.com/harvester/addons/pkg/render"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const (
	templateSource = "./pkg/templates"
)

func main() {
	var generateAddons bool
	var path string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "generateAddons",
				Value:       false,
				Usage:       "generate disabled addon yaml manifests",
				Destination: &generateAddons,
			},
			&cli.StringFlag{
				Name:        "path",
				Value:       ".",
				Usage:       "destination for output files",
				Destination: &path,
			},
		},

		Action: func(ctx *cli.Context) error {
			if !generateAddons {
				return fmt.Errorf("generateAddons need to be specified")
			}

			if generateAddons {
				if err := render.Addon(templateSource, path); err != nil {
					return fmt.Errorf("error during rendering addons: %v", err)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
