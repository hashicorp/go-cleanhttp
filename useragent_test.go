package cleanhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserAgent(t *testing.T) {
	var actualUserAgent string
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		actualUserAgent = req.UserAgent()
	}))
	defer ts.Close()

	for i, c := range []struct {
		expected string
		new      func() *http.Client
	}{
		{"Go-http-client/1.1", DefaultPooledClient},
		{"Go-http-client/1.1", func() *http.Client { return DefaultClient() }},
		{"", func() *http.Client { return DefaultClient(UserAgent("")) }},
		{"foo/1", func() *http.Client { return DefaultClient(UserAgent("foo/1")) }},
	} {
		t.Run(fmt.Sprintf("%d %s", i, c.expected), func(t *testing.T) {
			cli := c.new()
			actualUserAgent = ""
			_, err := cli.Get(ts.URL)
			if err != nil {
				t.Fatal(err)
			}
			if actualUserAgent != c.expected {
				t.Fatalf("actual User-Agent '%s' is not '%s'", actualUserAgent, c.expected)
			}
		})
	}
}
