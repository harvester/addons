package main

import (
	"fmt"
	"os"

	"github.com/harvester/addons/pkg/render"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	templateSource = "./pkg/templates"
)

func main() {
	var generateAddons, generateTemplates bool
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
				Usage:       "generate template manifests",
				Destination: &generateTemplates,
			},
			&cli.StringFlag{
				Name:        "path",
				Value:       ".",
				Usage:       "destination for output files",
				Destination: &path,
			},
		},

		Action: func(ctx *cli.Context) error {
			if !generateAddons && !generateTemplates {
				return fmt.Errorf("generateAddons or generateTemplates need to be specified")
			}

			if generateAddons {
				if err := render.Addon(templateSource, path, "version_info"); err != nil {
					return fmt.Errorf("error during rendering addons: %v", err)
				}
			}

			if generateTemplates {
				if err := render.Template(templateSource, path, "version_info"); err != nil {
					return fmt.Errorf("error during rendering template: %v", err)
				}
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}
