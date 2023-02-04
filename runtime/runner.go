package runtime

import (
	"bufio"
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
)

const (
	needle string = "https://example.com/"
)

type RunnerConfig struct {
	Threads      int
	ThreadsHost  int
	Smart        bool
	ProxyAddress string
	Insecure     bool
}

type runner struct {
	modify       chan bool
	cthreads     chan bool
	cthreadshost map[string]chan bool
	ithreadshost int
	smart        bool
	needle       *url.URL
	client       *http.Client
}

func NewRunner(cfg RunnerConfig) (*runner, error) {
	runner := new(runner)
	runner.modify = make(chan bool, 1)
	runner.cthreads = make(chan bool, cfg.Threads)
	runner.cthreadshost = make(map[string]chan bool)
	runner.ithreadshost = cfg.ThreadsHost
	runner.needle, _ = url.Parse(needle)
	runner.smart = cfg.Smart

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.Insecure},
	}

	if cfg.ProxyAddress != "" {
		u, err := url.Parse(cfg.ProxyAddress)
		if err != nil {
			return nil, err
		}

		transport.Proxy = http.ProxyURL(u)
	}

	runner.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}

	return runner, nil
}

func (r *runner) checkUrl(u *url.URL) (string, bool) {
	if u.RawQuery == "" {
		return "", false
	}

	parameters := u.Query()
	for k := range parameters {
		if r.smart && !IsOpenRedirectParam(k) {
			continue
		}

		v := parameters.Get(k)

		parameters.Set(k, needle)
		u.RawQuery = parameters.Encode()

		res, err := r.client.Get(u.String())
		if err != nil {
			break
		}

		if res.StatusCode >= 300 && res.StatusCode < 400 {
			location, err := res.Location()
			if err != nil {
				continue
			}

			if location.Host == r.needle.Host {
				return u.String(), true
			}
		}

		parameters.Set(k, v)
	}

	return "", false
}

func (r *runner) Run(sc *bufio.Scanner, output chan string) {
	wg := new(sync.WaitGroup)

	for sc.Scan() {
		u, err := url.Parse(sc.Text())
		if err != nil {
			continue
		}

		_, found := r.cthreadshost[u.Host]
		if !found {
			r.modify <- true
			r.cthreadshost[u.Host] = make(chan bool, r.ithreadshost)
			<-r.modify
		}

		ch, _ := r.cthreadshost[u.Host]

		wg.Add(1)
		go func() {
			r.cthreads <- true
			ch <- true

			link, vulnerable := r.checkUrl(u)
			if vulnerable {
				output <- link
			}

			<-ch
			<-r.cthreads
			wg.Done()
		}()
	}

	wg.Wait()
	close(output)
}
