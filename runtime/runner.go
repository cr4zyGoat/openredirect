package runtime

import (
	"bufio"
	"net/http"
	"net/url"
	"sync"
)

const (
	needle string = "https://example.com/"
)

type runner struct {
	modify       chan bool
	cthreads     chan bool
	cthreadshost map[string]chan bool
	ithreadshost int
	needle       *url.URL
	client       *http.Client
}

func NewRunner(threads int, threadshost int) *runner {
	runner := new(runner)
	runner.modify = make(chan bool, 1)
	runner.cthreads = make(chan bool, threads)
	runner.cthreadshost = make(map[string]chan bool)
	runner.ithreadshost = threadshost
	runner.needle, _ = url.Parse(needle)

	runner.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return runner
}

func (r *runner) SetTransport(transport *http.Transport) {
	r.client.Transport = transport
}

func (r *runner) checkUrl(u *url.URL) (string, bool) {
	if u.RawQuery == "" {
		return "", false
	}

	parameters := u.Query()
	for k := range parameters {
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
}
