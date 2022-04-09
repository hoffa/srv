# srv

[![Go Report Card](https://goreportcard.com/badge/github.com/hoffa/srv)](https://goreportcard.com/report/github.com/hoffa/srv)

Easily serve files over HTTP and HTTPS.

## Installation

```
go get github.com/hoffa/srv
```

## Example

```
srv -n mydomain.com
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
