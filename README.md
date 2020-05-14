# july

[![Scc Count Badge](https://sloc.xyz/github/vasyahuyasa/july/)](https://github.com/vasyahuyasa/july/)

July is simple OPDS home catalog without any external dependencies. Currntly supported only listing without covers and meta information.

## Install

```shell
$ go get github.com/vasyahuyasa/july/cmd/july
```

## Docker

https://hub.docker.com/r/vasyahuyasa/july

## Basic usage

```shell
Usage of july:
  -d string
        Root storage directory (default "./")
  -drv string
        Storage driver (can be local, gdrive, yadisk) (default "local")
  -googlecred string
        Path to file with secret for google drive driver (default "credentials.json")
  -i string
        Service network interface (default "0.0.0.0")
  -p int
        Service port (default 80)
```
