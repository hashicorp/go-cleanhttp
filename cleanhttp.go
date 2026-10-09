// Copyright IBM Corp. 2015, 2026
// SPDX-License-Identifier: MPL-2.0

package cleanhttp

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

// DefaultTransport returns a new http.Transport with similar default values to
// http.DefaultTransport, but with idle connections and keepalives disabled.
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
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true, //nolint:staticcheck
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost(),
	}
	return transport
}

// defaultMaxIdleConnsPerHost returns the default MaxIdleConnsPerHost value.
// Single and dual core environments underestimate concurrency requirements
// when using runtime.GOMAXPROCS(0) + 1. A baseline minimum floor of 4 prevents
// bottlenecking connection pools.
func defaultMaxIdleConnsPerHost() int {
	return calcMaxIdleConnsPerHost(runtime.GOMAXPROCS(0))
}

func calcMaxIdleConnsPerHost(procs int) int {
	const minIdleConns = 4
	return max(minIdleConns, procs+1)
}

// DefaultClient returns a new http.Client with similar default values to
// http.Client, but with a non-shared Transport, idle connections disabled, and
// keepalives disabled.
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: DefaultTransport(),
	}
}

// DefaultPooledClient returns a new http.Client with similar default values to
// http.Client, but with a shared Transport. Do not use this function for
// transient clients as it can leak file descriptors over time. Only use this
// for clients that will be re-used for the same host(s).
func DefaultPooledClient() *http.Client {
	return &http.Client{
		Transport: DefaultPooledTransport(),
	}
}
