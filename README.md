提供 httpdns 和 showmyip 两种服务

[![GoDoc](https://godoc.org/github.com/Tomasen/httpdns?status.svg)](http://godoc.org/github.com/Tomasen/httpdns)

#APIs

| 端口 | 协议 | 请求 | 返回 |
| ------ | ------ | ------ | ------ |
| 1053 | http | /dns?d={$domain} | 字符串，该域名的ip地址 |
| 1053 | http | /myip | 字符串，请求者的IP |
| 1053 | http | /health | 字符串，"OK" |
| 1153 | tcp  | 字符串类型，$domain + "\\n" | 字符串，该域名的ip地址 |
| 1154 | tcp  | 无 | 字符串类型，请求者的IP + '\\n'，并关闭连接 |


#Contribution

Commited code must pass:

`golint`
`go vet`
`gofmt`
`go test`

# TODO

* test cases
* Dockerfile
* continue deployment on AWS and aliyun