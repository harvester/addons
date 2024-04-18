# addons

Repo contains addon info under `/pkg/templates` and stores them as `go bindata` to be used by
`harvester/harvester` and `harvester/harvester-installer` repos.

The cli generates two ouput types
* raw template: to be copied into the harvester-installer before compilation to ensure that the addons can be enabled/disabled during the install phase
* disabled addon: to be copied into the harvester repo before compilation to ensure same addon info is available in the upgrade path

To update addons please ensure that the templates under `pkg/templates` are update and `make` is executed to update bindata.

The bindata also needs to be committed to ensure correct version is extract during the packaging of harvester and harvester-installer

The repo also contains a `version_info` file which is sourced by `harvester-installer` build-bundle script

Please ensure image and chart info update is also reflected in this file.