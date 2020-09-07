# Notifier

Utility that ...

[![GoDoc](https://godoc.org/github.com/devopsext/notifier?status.svg)](https://godoc.org/github.com/devopsext/notifier)
[![build status](https://img.shields.io/travis/devopsext/notifier/master.svg?style=flat-square)](https://travis-ci.org/devopsext/notifier)

## Features

- Pub mode
- Sub mode

## Build

```sh
git clone https://github.com/devopsext/notifier.git
cd notifier/
go build
```

## Example

### Run in Sub mode

```sh
./notifier sub --output json --file scraper.json \
               --log-format stdout --log-level debug --log-template '{{.msg}}' 
```

### Run in Pub mode

```sh
./notifier pub --output json --file scraper.json \
               --log-format stdout --log-level debug --log-template '{{.msg}}' 
```

## Environment variables

For containerization purpose all command switches have environment variables analogs.

- SCRAPER_LOG_FORMAT
- SCRAPER_LOG_LEVEL
- SCRAPER_LOG_TEMPLATE
