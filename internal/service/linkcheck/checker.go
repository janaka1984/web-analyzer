package linkcheck

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type httpDoer interface {
	Head(ctx context.Context, url string) (*http.Response, error)
	Get(ctx context.Context, url string) (*http.Response, error)
}

type Checker struct {
	workers int
}

func New(workers int) *Checker { return &Checker{workers: workers} }

func (c *Checker) CheckLinks(ctx context.Context, baseURL string, links []string) (internal, external, broken int) {
	type job struct{ url string }
	type res struct {
		url      string
		status   int
		err      error
		internal bool
	}
	jobs := make(chan job)
	results := make(chan res)

	// Determine internal/external by host
	// extract base host
	baseHost := hostOf(baseURL)

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		client := &http.Client{Timeout: 8 * time.Second}
		for j := range jobs {
			r := res{url: j.url, internal: hostOf(j.url) == baseHost}
			// HEAD first
			req, _ := http.NewRequestWithContext(ctx, http.MethodHead, j.url, nil)
			resp, err := client.Do(req)
			if err == nil && resp != nil {
				r.status = resp.StatusCode
				resp.Body.Close()
			} else {
				// fallback GET (some servers reject HEAD)
				req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, j.url, nil)
				resp2, err2 := client.Do(req2)
				if resp2 != nil {
					r.status = resp2.StatusCode
					resp2.Body.Close()
				}
				if err2 != nil {
					r.err = err2
				}
			}
			results <- r
		}
	}

	// start workers
	wg.Add(c.workers)
	for i := 0; i < c.workers; i++ {
		go worker()
	}

	go func() {
		for _, u := range links {
			select {
			case jobs <- job{url: u}:
			case <-ctx.Done():
				break
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.internal {
			internal++
		} else {
			external++
		}
		if r.err != nil || r.status >= 400 || r.status == 0 {
			broken++
		}
	}
	return
}

func hostOf(u string) string {
	// minimal host detector
	// (avoid full parse for performance; ok for test task)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil || req.URL == nil {
		return ""
	}
	return req.URL.Host
}
