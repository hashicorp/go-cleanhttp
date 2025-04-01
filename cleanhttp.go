// Copyright (c) HashiCorp, Inc.
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
	transport := DefaultPooledTransportWithMin(0)
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// DefaultPooledTransport returns a new http.Transport with similar default
// values to http.DefaultTransport. Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
func DefaultPooledTransport() *http.Transport {
	return DefaultPooledTransportWithMin(0)
}

// DefaultPooledTransportWithMin returns a new http.Transport with a minimum
// value for MaxIdleConnsPerHost. Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
func DefaultPooledTransportWithMin(minIdleConnsPerHost int) *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost: func() int {
			maxIdleConnsPerHost := runtime.GOMAXPROCS(0) + 1
			if maxIdleConnsPerHost < minIdleConnsPerHost {
				return minIdleConnsPerHost
			}
			return maxIdleConnsPerHost
		}(),
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

// DefaultPooledClient returns a new http.Client with similar default values to
// http.Client, but with a shared Transport. Do not use this function for
// transient clients as it can leak file descriptors over time. Only use this
// for clients that will be re-used for the same host(s).
func DefaultPooledClient() *http.Client {
	return DefaultPooledClientWithMin(0)
}

// DefaultPooledClient returns a new http.Client with similar default values to
// http.Client, but with a shared Transport and minimum value for MaxIdleConnsPerHost.
// Do not use this function for transient clients as it can leak file descriptors
// over time. Only use this for clients that will be re-used for the same host(s).
func DefaultPooledClientWithMin(minIdleConnsPerHost int) *http.Client {
	return &http.Client{
		Transport: DefaultPooledTransportWithMin(minIdleConnsPerHost),
	}
}
