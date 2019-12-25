[Unreleased]


[v0.6.4] - 2019-12-24

Safer `omw edit`

    - Backs up current file to the same path with a `.bak` extension
    - Ensures valid TOML
    - Automatically fixes duplicate IDs - makes it easy to manually copy/paste a new entry in edit mode

[v0.6.1] - 2019-12-16

First release

    - Default log file format is now structured TOML
    - Added a standalone conversion utility to migrate from old data format
    - Initial migration of build/release to goreleaser Github action
    - Removed all unnecessary components after separation of web app
