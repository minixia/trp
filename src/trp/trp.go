package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"regexp"
	"proxy"
	"runtime"
)

const (
	Version = "0.1.0"
)


// backend regular expression (<host>:<port>)
var backendRe *regexp.Regexp = regexp.MustCompile("^[^:]+:[0-9]+$")

// isValidBackend returns true if backend is in "host:port" format
func isValidBackend(backend string) bool {
	return backendRe.MatchString(backend)
}

// parseBackends parses string in format "host:port,host:port" and return list of backends
func parseBackends(str string) ([]string, error) {
	backends := strings.Split(str, ",")
	if len(backends) == 0 {
		return nil, fmt.Errorf("no backends")
	}

	for i, v := range backends {
		backends[i] = strings.TrimSpace(v)
		if !isValidBackend(backends[i]) {
			return nil, fmt.Errorf("'%s' is not valid network address", backends[i])
		}
	}

	return backends, nil
}

// die prints error message and aborts the program
func die(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: trp -l LISTEN_ADDR -b BACKENDS\n")
		fmt.Fprintf(os.Stderr, "command line switches:\n")
		flag.PrintDefaults()
	}
	listenAddr := flag.String("l", ":80", "listen ip:port")
	version := flag.Bool("v", false, "show version and exit")
	backends := flag.String("b", "127.0.0.1:1935,127.0.0.2:1935", "backends for proxy")
	flag.Parse()

	if *version {
		fmt.Printf("Tcp Reverse Proxy v%s build by %s\n", Version, runtime.Version())
		os.Exit(0)
	}

	if flag.NFlag() < 2 {
		flag.Usage()
		os.Exit(1)
	}


	backendList, err := parseBackends(*backends)
	if err != nil {
		die(fmt.Sprintf("%s", err))
	}
	proxy.ProxyServer(backendList, *listenAddr)
}
