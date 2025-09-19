# addons

Repo contains addon info under `/pkg/templates` 

The templates need to be generated for harvester-installer packaging as follows:
`go run . -generateTemplates -path $path_to_installer_templates`

For harvester upgrade path, the templates need to be rendered, and easiest way to do the same is to call

`go run . -generateAddons -path $upgrade_path_manifests`

The repo also contains a `version_info` file which is sourced by `harvester-installer` build-bundle script

Please ensure image and chart info update is also reflected in this file.

## run `make` with dapper

All following commands run in similar way like most Harvester repos.

Run `make generate` to generate the addon templates, which is saved under `./bin`.

Run `make test-chart-patch` to test the patches upon `rancher-monitoring` and `rancher-logging` charts, the patched charts are saved under `./bin/patched-charts`.
