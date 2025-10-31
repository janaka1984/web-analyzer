package humanizer

import (
	"fmt"
	"strings"
)

// HTTPError converts a raw HTTP/network error string + status code
// into a concise, user-friendly message.
func HTTPError(raw string, status int) string {
	rawLower := strings.ToLower(raw)

	switch {
	case strings.Contains(rawLower, "403"), strings.Contains(rawLower, "access denied"):
		return "Access denied — this site blocks automated requests or requires login."
	case strings.Contains(rawLower, "timeout"), strings.Contains(rawLower, "deadline exceeded"):
		return "Request timed out — the site took too long to respond."
	case strings.Contains(rawLower, "x509"), strings.Contains(rawLower, "certificate"):
		return "SSL/TLS certificate error — the site’s certificate may be invalid."
	case strings.Contains(rawLower, "no such host"):
		return "Invalid domain — please check the website address."
	case strings.Contains(rawLower, "connection refused"):
		return "The site refused the connection or is temporarily unavailable."
	case status >= 400 && status < 500:
		return fmt.Sprintf("Client error (%d) — check the URL or permissions.", status)
	case status >= 500:
		return fmt.Sprintf("Server error (%d) — the target site may be down.", status)
	default:
		return "Unexpected error while analyzing the website."
	}
}
