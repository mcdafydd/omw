[Unreleased]

- Integrate omw CLI with omw progessive web app, controlled by configuration

[v0.7.0] - 2020-01-20

Remove `omw server`

    - Removed localhost TCP server - all omw functionality will now be local
    - Rename `start` to `end` in timesheet for clarity - format should be compatible with
    the PWA Dexie format
    - Update command help text

[v0.6.6] - 2019-12-24

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
