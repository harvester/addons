//go:generate go run ./pkg/bindata

package main

import (
	"fmt"
	"github.com/harvester/addons/pkg/render"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var generateTemplates, generateAddons bool
	var path string
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "generateAddons",
				Value:       false,
				Usage:       "generate disabled addon yaml manifests",
				Destination: &generateAddons,
			},
			&cli.BoolFlag{
				Name:        "generateTemplates",
				Value:       false,
				Usage:       "generate addon template files",
				Destination: &generateTemplates,
			},
			&cli.StringFlag{
				Name:        "path",
				Value:       ".",
				Usage:       "destiation for output files",
				Destination: &path,
			},
		},

		Action: func(ctx *cli.Context) error {
			if !generateTemplates && !generateAddons {
				return fmt.Errorf("either generateTemplates or generateAddons need to be specified")
			}

			if generateAddons {
				if err := render.Addon(path); err != nil {
					return fmt.Errorf("error during rendering addons: %v", err)
				}
			}

			if generateTemplates {
				if err := render.Template(path); err != nil {
					return fmt.Errorf("error during rendering templates: %v", err)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
