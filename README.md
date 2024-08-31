# Tesseris - Batteries included framework for Go

## Documentation

See documentation at https://tesseris.dev/docs

## Releasing a new version

### Build

Build a new version.

```sh
# Update the version number based on what kind of version you want to release
go run ./next-version [patch|minor|major]

# TODO: make build command
cd cmd/tesseris
go build
```

### Add tag

```sh
./add-version-tag.sh
```


