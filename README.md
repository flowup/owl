# Owl

[![Build Status](https://travis-ci.org/flowup/owl.svg?branch=master)](https://travis-ci.org/flowup/owl) [![Coverage Status](https://coveralls.io/repos/github/flowup/owl/badge.svg?branch=master)](https://coveralls.io/github/flowup/owl?branch=master)

## Overview
:rocket:-speed file-watcher written in Golang, Owl is mostly suitable as an automatic build/run/test tool.

##Installation
```bash
$ go get github.com/flowup/owl/cmd/owl # this will install the binary in $GOBIN
```

## Usage

You can use **owl** to simply run tests when anything within the current folder(recursively) changes. The `-i` flag will ignore a directory named `bin`

```bash
$ owl -r 'go test ./...' -i bin
```

### Flags

- `-i` or `--ignore` to ignore files and folders
- `-r` or `--run` for specific command
- `-t` or `--time` debounce time for filesystem events before command execution in miliseconds (default 500)

### Config file owl.yaml

You can set default settings for the `owl` command within the folder with config file.

> :robot: Note that any environment variables and flags will override this configuration

```
run: "echo \"Hello Owl!\""
time: 5000
verbose: true
ignore:
 - vendor"
 - "bin"
```
