package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

const (
	needle string = "https://randomfakewebsite.com"
)

func checkUrl(uri string) (string, bool) {
	uneedle, _ := url.Parse(needle)
	u, err := url.Parse(uri)
	if err != nil {
		log.Println()
		return "", false
	}

	if u.RawQuery == "" {
		return "", false
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	parameters := u.Query()
	for k := range parameters {
		v := parameters.Get(k)

		parameters.Set(k, needle)
		u.RawQuery = parameters.Encode()

		res, err := client.Get(u.String())
		if err != nil {
			return "", false
		}

		if res.StatusCode >= 300 && res.StatusCode < 400 {
			location, err := res.Location()
			if err != nil {
				continue
			}

			if location.Host == uneedle.Host {
				return u.String(), true
			}
		}

		parameters.Set(k, v)
	}

	return "", false
}

func main() {
	pfile := flag.String("f", "", "File with the target URLs")
	pthreads := flag.Int("t", 10, "Threads to execute")
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

	wg := new(sync.WaitGroup)
	cthreads := make(chan bool, *pthreads)

	for sc.Scan() {
		uri := sc.Text()
		wg.Add(1)

		go func(uri string) {
			cthreads <- true

			openredirect, vulnerable := checkUrl(uri)
			if vulnerable {
				fmt.Println(openredirect)
			}

			<-cthreads
			wg.Done()
		}(uri)
	}

	wg.Wait()
}
