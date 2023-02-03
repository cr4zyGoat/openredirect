package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/cr4zygoat/openredirect/runtime"
)

func main() {
	pfile := flag.String("f", "", "File with the target URLs")
	pproxy := flag.String("proxy", "", "Proxy with the format 'http://proxyhost:port'")
	pthreads := flag.Int("t", 150, "Max number of URLs to check at the same time")
	pthreadshost := flag.Int("tpd", 3, "Threads to execute per host")
	pinsecure := flag.Bool("insecure", false, "Skip TLS verification")
	flag.Parse()

	var sc *bufio.Scanner
	if *pfile == "" {
		sc = bufio.NewScanner(os.Stdin)
	} else {
		pf, err := os.Open(*pfile)
		if err != nil {
			log.Fatalln(err)
		}

		sc = bufio.NewScanner(pf)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *pinsecure},
	}

	if *pproxy != "" {
		u, err := url.Parse(*pproxy)
		if err != nil {
			log.Println(err)
		}

		transport.Proxy = http.ProxyURL(u)
	}

	runner := runtime.NewRunner(*pthreads, *pthreadshost)
	runner.SetTransport(transport)

	output := make(chan string)
	go runner.Run(sc, output)

	for link := range output {
		fmt.Println(link)
	}
}
