# Tesseris - Batteries included framework for Go

## Documentation

See documentation at https://tesseris.iskandervdh.nl/docs

## Releasing a new version

First commit all changes and push them to the repository.

Afterwards the version number can be updated based on what kind of version update it is:

```sh
go run ./next-version [patch|minor|major]
```

Then the version tag can be added to the repository using the following command:

```sh
./add-version-tag.sh
```


