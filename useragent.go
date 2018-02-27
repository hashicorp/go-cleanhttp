package cleanhttp

import "net/http"

type userAgentRoundTripper struct {
	inner     http.RoundTripper
	userAgent string
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", rt.userAgent)
	return rt.inner.RoundTrip(req)
}

// UserAgent modifies the http.Client's Transport to inject
// a UserAgent header on the request.
func UserAgent(v string) func(*http.Client) {
	return func(c *http.Client) {
		c.Transport = &userAgentRoundTripper{
			inner:     c.Transport,
			userAgent: v,
		}
	}
}
