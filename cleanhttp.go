package cleanhttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// DefaultTransport returns a new http.Transport with the same default values
// as http.DefaultTransport, but with idle connections and keepalives disabled.
func DefaultTransport() *http.Transport {
	transport := DefaultPooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// DefaultPooledTransport returns a new http.Transport with similar default
// values to http.DefaultTransport. Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
func DefaultPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 1,
	}
	return transport
}

// NoTLSVerifyTransport returns a new http.Transport with options identical to
// DefaultTransport however with TLS verification turned off. Intended only for
// use with internal resources that are not routed on the public Internet.
func NoTLSVerifyTransport() *http.Transport {

	transport := &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	DisableKeepAlives:   true,
	MaxIdleConnsPerHost: -1,
	}
	return transport
}

// DefaultClient returns a new http.Client with similar default values to
// http.Client, but with a non-shared Transport, idle connections disabled, and
// keepalives disabled.
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: DefaultTransport(),
	}
}

// DefaultPooledClient returns a new http.Client with the same default values
// as http.Client, but with a shared Transport. Do not use this function
// for transient clients as it can leak file descriptors over time. Only use
// this for clients that will be re-used for the same host(s).
func DefaultPooledClient() *http.Client {
	return &http.Client{
		Transport: DefaultPooledTransport(),
	}
}

// NoTLSVerifyClient returns a new http.Client with options identical to
// DefaultClient however with TLS verification turned off. Intended only for
// use with internal resources that are not available on the public Internet.
func NoTLSVerifyClient() *http.Client {
	return &http.Client{
		Transport: NoTLSVerifyTransport(),
	}
}
