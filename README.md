# srv

[![Go Report Card](https://goreportcard.com/badge/github.com/hoffa/srv)](https://goreportcard.com/report/github.com/hoffa/srv)

Serve files over HTTPS.

## Installation

```
go get github.com/hoffa/srv
```

## Usage

```
Usage of srv:
  -c string
        directory to store TLS certificates (temporary directory if empty)
  -d string
        directory to serve (default ".")
  -n string
        domain name (required)
  -p int
        HTTPS port (default 443)
  -q int
        HTTP redirect port (default 80)
  -t duration
        timeout (default 1m0s)
```

## Example

```
srv -n mydomain.com
```
