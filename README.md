Provide httpdns and showmyip services in golang.
提供 httpdns 和 showmyip 两种服务 in golang

[![Build Status](https://travis-ci.org/Tomasen/httpdns.svg?branch=master)](https://travis-ci.org/Tomasen/httpdns)
[![GoDoc](https://godoc.org/github.com/Tomasen/httpdns?status.svg)](http://godoc.org/github.com/Tomasen/httpdns)

## APIs

| 端口 | 协议 | 请求 | 返回 |
| ------ | ------ | ------ | ------ |
| 1053 | http | /dns?d={$domain} | 字符串，该域名的ip地址 |
| 1053 | http | /myip | 字符串，请求者的IP |
| 1053 | http | /health | 字符串，"OK" |
| 1153 | tcp  | 字符串类型，$domain + "\\n" | 字符串，该域名的ip地址 |
| 1154 | tcp  | 无 | 字符串类型，请求者的IP + '\\n'，并关闭连接 |

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
* Dockerfile
* continue deployment on AWS and aliyun