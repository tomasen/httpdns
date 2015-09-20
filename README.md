Provide httpdns and show-my-ip services in golang.

[![Build Status](https://travis-ci.org/Tomasen/httpdns.svg?branch=master)](https://travis-ci.org/Tomasen/httpdns)
[![GoDoc](https://godoc.org/github.com/Tomasen/httpdns?status.svg)](http://godoc.org/github.com/Tomasen/httpdns)

## APIs

| port | scheme | request | response |
| ------ | ------ | ------ | ------ |
| 1053 | http | /dns?d={$domain} | type: string, ip address of $domain |
| 1053 | http | /myip | type: string, ip address of client |
| 1053 | http | /health | type: string, "OK" indicate health |
| 1153 | tcp  | type: string, $domain + "\\n" | type: string, ip address of $domain |
| 1154 | tcp  | not required | type: string，ip address of client + '\\n' |

## Deployment

`docker run -p 1053:1053 -p 1153-1154:1153-1154 tomasen/httpdns`

## Developing

Commited code must pass:

* [golint](https://github.com/golang/lint)
* [go vet](https://godoc.org/golang.org/x/tools/cmd/vet)
* [gofmt](https://golang.org/cmd/gofmt)
* [go test](https://golang.org/cmd/go/#hdr-Test_packages):

# TODO

* test cases
