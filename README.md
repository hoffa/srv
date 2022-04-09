# srv

[![Go Report Card](https://goreportcard.com/badge/github.com/hoffa/srv)](https://goreportcard.com/report/github.com/hoffa/srv)
[![.github/workflows/release.yml](https://github.com/hoffa/srv/actions/workflows/release.yml/badge.svg)](https://github.com/hoffa/srv/actions/workflows/release.yml)

Easily serve files over HTTPS.

## Installation

```
go get github.com/hoffa/srv
```

## Example

```
mkdir www
echo "<h1>Hi!</h1>" > www/index.html
mkdir certs
srv -n mydomain.com -d www -c certs
```
