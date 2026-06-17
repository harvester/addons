# addons

Repo contains addon info under `/pkg/templates` 

The templates need to be generated for harvester-installer packaging as follows:
`go run . -generateTemplates -path $path_to_installer_templates`

For harvester upgrade path, the templates need to be rendered, and easiest way to do the same is to call

`go run . -generateAddons -path $upgrade_path_manifests`

The repo also contains a `version_info` file which is sourced by `harvester-installer` build-bundle script

Please ensure image and chart info update is also reflected in this file.

## Run with `make` and `docker`

All following commands run in similar way like most Harvester repos.

Run `make generate` to generate the addon templates, which is saved under `./output`.

Run `make patch-charts` or scripts/patch-charts to test the patches to the charts (`rancher-monitoring` and `rancher-logging`).
The patched charts are saved under `./output/patched-charts`.

If you want to add more charts and patches, just place the patches under
`pkg/config/templates/patch/$CHART/$VERSION` and add a script `scripts/hack/patch-$CHART`.
The script should contain three functions: `pull_$CHART_chart`, `patch_$CHART_chart` and
`test_$CHART_chart`, see the existing scripts for reference.
