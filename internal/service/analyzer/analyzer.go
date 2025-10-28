package analyzer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/janaka/web-analyzer/internal/domain"
	"github.com/janaka/web-analyzer/internal/service/linkcheck"
)

type Fetcher interface {
	Get(ctx context.Context, url string) (*http.Response, error)
	Head(ctx context.Context, url string) (*http.Response, error)
}

type Checker interface {
	CheckLinks(ctx context.Context, baseURL string, links []string) (internal, external, broken int)
}

type Analyzer struct {
	fetch  Fetcher
	check  Checker
	limits Limits
}

type Limits struct {
	MaxLinks int
}

func New(fetcher Fetcher, checker Checker, maxLinks int) *Analyzer {
	return &Analyzer{fetch: fetcher, check: checker, limits: Limits{MaxLinks: maxLinks}}
}

func (a *Analyzer) Analyze(ctx context.Context, url string) (*domain.Analysis, error) {
	resp, err := a.fetch.Get(ctx, url)
	res := &domain.Analysis{URL: url}
	if err != nil {
		res.ErrorMessage = err.Error()
		return res, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()
	res.HTTPStatus = resp.StatusCode

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		b, _ := io.ReadAll(io.LimitedReader{R: resp.Body, N: 4096})
		res.ErrorMessage = fmt.Sprintf("non-2xx status %d: %s", resp.StatusCode, string(b))
		return res, nil
	}

	body, err := io.ReadAll(io.LimitedReader{R: resp.Body, N: 5 << 20})
	if err != nil {
		res.ErrorMessage = err.Error()
		return res, fmt.Errorf("read: %w", err)
	}

	// parse
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		res.ErrorMessage = err.Error()
		return res, fmt.Errorf("parse: %w", err)
	}

	res.HTMLVersion = DetectHTMLVersion(body)
	res.Title = ExtractTitle(doc)
	h := CountHeadings(doc)
	res.H1, res.H2, res.H3, res.H4, res.H5, res.H6 = h[1], h[2], h[3], h[4], h[5], h[6]
	res.HasLoginForm = HasLoginForm(doc)

	links := ExtractLinks(doc, url)
	if a.limits.MaxLinks > 0 && len(links) > a.limits.MaxLinks {
		links = links[:a.limits.MaxLinks]
	}
	in, ex, broken := a.check.CheckLinks(ctx, url, links)
	res.InternalLinks, res.ExternalLinks, res.BrokenLinks = in, ex, broken

	return res, nil
}

// Helper factory for default analyzer
func NewDefaultAnalyzer(fetcher Fetcher, workers, maxLinks int) *Analyzer {
	return New(fetcher, linkcheck.New(workers), maxLinks)
}
