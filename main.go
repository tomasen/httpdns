package main

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"

	"github.com/Tomasen/realip"
)

const (
	MAX_OPENFILE      uint64 = 128000
	MAX_DOMAIN_LENGTH        = 288 // 255 according to rfc3986
)

const (
	QUERY_DNS int8 = 1 << iota
	QUERY_MYIP
)
func LogError(v interface{}) {
	// TODO: log error but with a rate limit and a rate record
}

func QueryDNS(domain string, src_ip string) []byte {
	// TODO: use edns-client-subnet from google or other service provider

	ips, err := net.LookupIP(domain)
	if err != nil {
		LogError(err)
		return nil
	}
	return ips[rand.Intn(len(ips))]
}

func HTTPServerDNS(w http.ResponseWriter, req *http.Request) {
	domain := req.Form.Get("d")
	if len(domain) == 0 {
		http.NotFound(w, req)
		return
	}
	w.Write([]byte(QueryDNS(domain, realip.RealIP(req))))
}

func HTTPServerMYIP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(realip.RealIP(req)))
}

func HTTPServerHealth(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}

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
				}
			}()
			defer c.Close()
			// TODO: detect proxy protocol, eg.
			// PROXY TCP4 192.168.0.1 192.168.0.11 56324 80rn
			// GET / HTTP/1.1rn
			// Host: 192.168.0.11rn
			// rn

			var r []byte
			switch t {
			case QUERY_DNS:
				line, _, err := bufio.NewReaderSize(c, MAX_DOMAIN_LENGTH).ReadLine()
				if err != nil {
					// TODO: log error
					return
				}
				r = QueryDNS(string(line), c.RemoteAddr().String())
			case QUERY_MYIP:
				r = []byte(c.RemoteAddr().String())
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
	if lim.Cur < MAX_OPENFILE || lim.Max < MAX_OPENFILE {
		lim.Cur = MAX_OPENFILE
		lim.Max = MAX_OPENFILE
		syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim)
	}

	ln, err := net.Listen("tcp", ":1153")
	if err != nil {
		log.Fatal(err)
	}
	go TCPServer(ln, QUERY_DNS)

	ln, err = net.Listen("tcp", ":1154")
	if err != nil {
		log.Fatal(err)
	}
	go TCPServer(ln, QUERY_MYIP)

	http.HandleFunc("/myip", HTTPServerMYIP)
	http.HandleFunc("/dns", HTTPServerDNS)
	http.HandleFunc("/health", HTTPServerHealth)
	log.Fatal(http.ListenAndServe(":1053", nil))
}
