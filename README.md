# libconfig

`libconfig` is a Go library that provides a simple method for populating a struct with data from environment variables.

[![GoDoc](https://godoc.org/github.com/jrudder/libconfig?status.png)](https://godoc.org/github.com/jrudder/libconfig)
[![CircleCI](https://circleci.com/gh/jrudder/libconfig.svg?style=svg)](https://circleci.com/gh/jrudder/libconfig)
[![Coverage](https://codecov.io/gh/jrudder/libconfig/branch/master/graph/badge.svg)](https://codecov.io/gh/jrudder/libconfig)

## Goals

`libconfig` provides the following:

- [x] Parsing to built-in types
- [x] Parsing to slices and pointers
- [x] Support defaults, optionals, base64, and json
- [x] 100% code coverage

## Non-Goals

`libconfig` does not and will not support these scenarios:

* Config reloading: This is the work of deployment, canary releases, and blue/green environments.
* Loading config from a file: Config should come from the environment only so that all environments (local, dev, prod, etc.) use the same code paths.

## Prior Art

* [gonfig](https://github.com/tkanos/gonfig) is a lightweight parser of JSON and env vars.
* [Stært](https://github.com/containous/staert) can be used alone or with [Flæg](https://github.com/containous/flaeg) to provide config parsing from the command line, TOML, and key-value stores using [libkv](https://github.com/docker/libkv).
* [Viper](https://github.com/spf13/viper) is a feature-filled library with support for defaults, various formats (including JSON, TOML, and YAML), watching, and more.
