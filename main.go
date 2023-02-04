package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cr4zygoat/openredirect/runtime"
)

func main() {
	pfile := flag.String("f", "", "File with the target URLs")
	pproxy := flag.String("proxy", "", "Proxy with the format 'http://proxyhost:port'")
	pthreads := flag.Int("t", 150, "Max number of URLs to check at the same time")
	pthreadshost := flag.Int("tpd", 3, "Threads to execute per domain")
	psmart := flag.Bool("smart", false, "Do not check all the params")
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

	rconfig := runtime.RunnerConfig{
		Threads:      *pthreads,
		ThreadsHost:  *pthreadshost,
		Smart:        *psmart,
		ProxyAddress: *pproxy,
		Insecure:     *pinsecure,
	}

	runner, err := runtime.NewRunner(rconfig)
	if err != nil {
		log.Fatalln(err)
	}

	output := make(chan string)
	go runner.Run(sc, output)

	for link := range output {
		fmt.Println(link)
	}
}
