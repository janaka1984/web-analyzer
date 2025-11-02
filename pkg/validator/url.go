package validator

import (
	"errors"
	"net/url"
	"strings"
)

// NormalizeURL ensures a valid absolute URL (used for form input)
func NormalizeURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("empty url")
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", errors.New("invalid url")
	}
	return u.String(), nil
}

// NormalizeLink standardizes a link to avoid duplicates
// e.g https://example.com, https://example.com/ and https://example.com?utm_source=abc
func NormalizeLink(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	// Remove fragment (#section)
	u.Fragment = ""

	// Remove UTM and referral parameters
	q := u.Query()
	for k := range q {
		if strings.HasPrefix(strings.ToLower(k), "utm_") || k == "ref" {
			q.Del(k)
		}
	}
	u.RawQuery = q.Encode()

	// Remove trailing slash (but not root)
	if strings.HasSuffix(u.Path, "/") && len(u.Path) > 1 {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	return u.String()
}
