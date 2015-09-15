package main

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"syscall"

	"github.com/Tomasen/realip"
)

const (
	// max open file should at least be
	_MaxOpenfile uint64 = 128000

	// max domain string length should be 255
	// 255 according to rfc3986. padding is added
	_MaxDomainLength = 288
)

const (
	// types of query
	queryTypeDNS int8 = 1 << iota
	queryTypeMYIP
)

func logError(v ...interface{}) {
	// TODO: log error but with a rate limit and a rate record
}

// QueryDNS lookup IP from specified domain
func QueryDNS(domain string, srcip string) []byte {
	// TODO: use edns-client-subnet from google or other service provider

	ips, err := net.LookupHost(domain)
	if err != nil {
		logError(err)
		return nil
	}
	return []byte(ips[rand.Intn(len(ips))])
}

// HTTPServerDNS is handler of httpdns query
func HTTPServerDNS(w http.ResponseWriter, req *http.Request) {

	domain := req.URL.Query().Get("d")
	if len(domain) == 0 {
		domain := req.Form.Get("d")
		if len(domain) == 0 {
			http.NotFound(w, req)
			return
		}
	}
	w.Write([]byte(QueryDNS(domain, realip.RealIP(req))))
}

// HTTPServerMYIP is handler of show-my-ip query
func HTTPServerMYIP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(realip.RealIP(req)))
}

// HTTPServerHealth is handler of health check query
func HTTPServerHealth(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}

// TCPServer is handler for all tcp queries
func TCPServer(l net.Listener, t int8) {
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			defer func() {
				if r := recover(); r != nil {
					// TODO: log error
					logError("Recovered in", r, ":", string(debug.Stack()))
				}
			}()
			defer c.Close()
			// TODO: detect proxy protocol, eg.
			// PROXY TCP4 192.168.0.1 192.168.0.11 56324 80rn
			// GET / HTTP/1.1rn
			// Host: 192.168.0.11rn
			// rn
			h, _, _ := net.SplitHostPort(c.RemoteAddr().String())
			var r []byte
			switch t {
			case queryTypeDNS:
				line, _, err := bufio.NewReaderSize(c, _MaxDomainLength).ReadLine()
				if err != nil {
					// TODO: log error
					return
				}
				r = QueryDNS(string(line), h)
			case queryTypeMYIP:
				r = []byte(h)
			}
			c.Write(r)
		}(conn)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Setenv("GOTRACEBACK", "crash")

	lim := syscall.Rlimit{}
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim)
	if lim.Cur < _MaxOpenfile || lim.Max < _MaxOpenfile {
		lim.Cur = _MaxOpenfile
		lim.Max = _MaxOpenfile
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
	}

	ln, err := net.Listen("tcp", ":1153")
	if err != nil {
		log.Fatal(err)
	}
	go TCPServer(ln, queryTypeDNS)

	ln, err = net.Listen("tcp", ":1154")
	if err != nil {
		log.Fatal(err)
	}
	go TCPServer(ln, queryTypeMYIP)

	http.HandleFunc("/myip", HTTPServerMYIP)
	http.HandleFunc("/dns", HTTPServerDNS)
	http.HandleFunc("/health", HTTPServerHealth)
	log.Fatal(http.ListenAndServe(":1053", nil))
}
