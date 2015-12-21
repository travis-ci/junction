package http

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/travis-ci/junction/junction"
)

func TestListener(t *testing.T) (net.Listener, string) {
	fail := func(format string, args ...interface{}) {
		panic(fmt.Sprintf(format, args...))
	}
	if t != nil {
		fail = t.Fatalf
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fail("err: %s", err)
	}
	addr := "http://" + ln.Addr().String()
	return ln, addr
}

func TestServerWithListener(t *testing.T, ln net.Listener, addr string, core *junction.Core) {
	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: Handler(core),
	}
	go server.Serve(ln)
}

func TestServer(t *testing.T, core *junction.Core) (net.Listener, string) {
	ln, addr := TestListener(t)
	TestServerWithListener(t, ln, addr, core)
	return ln, addr
}
