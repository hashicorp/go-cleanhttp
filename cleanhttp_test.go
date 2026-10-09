// Copyright IBM Corp. 2015, 2026
// SPDX-License-Identifier: MPL-2.0

package cleanhttp

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func TestCalcMaxIdleConnsPerHost(t *testing.T) {
	cases := []struct {
		name     string
		procs    int
		expected int
	}{
		{
			name:     "negative procs clamped to floor",
			procs:    -1,
			expected: 4,
		},
		{
			name:     "zero procs clamped to floor",
			procs:    0,
			expected: 4,
		},
		{
			name:     "single core clamped to floor",
			procs:    1,
			expected: 4,
		},
		{
			name:     "dual core clamped to floor",
			procs:    2,
			expected: 4,
		},
		{
			name:     "three cores boundary",
			procs:    3,
			expected: 4,
		},
		{
			name:     "four cores scales above floor",
			procs:    4,
			expected: 5,
		},
		{
			name:     "eight cores scales above floor",
			procs:    8,
			expected: 9,
		},
		{
			name:     "sixteen cores scales above floor",
			procs:    16,
			expected: 17,
		},
		{
			name:     "thirty-two cores scales above floor",
			procs:    32,
			expected: 33,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := calcMaxIdleConnsPerHost(tc.procs)
			if actual != tc.expected {
				t.Fatalf("calcMaxIdleConnsPerHost(%d) = %d; want %d", tc.procs, actual, tc.expected)
			}
		})
	}
}

func TestDefaultPooledTransport_GOMAXPROCSFloor(t *testing.T) {
	origProcs := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(origProcs)

	cases := []struct {
		name     string
		setProcs int
		expected int
	}{
		{
			name:     "gomaxprocs 1 enforces floor",
			setProcs: 1,
			expected: 4,
		},
		{
			name:     "gomaxprocs 2 enforces floor",
			setProcs: 2,
			expected: 4,
		},
		{
			name:     "gomaxprocs 4 computes procs plus one",
			setProcs: 4,
			expected: 5,
		},
		{
			name:     "gomaxprocs 8 computes procs plus one",
			setProcs: 8,
			expected: 9,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runtime.GOMAXPROCS(tc.setProcs)
			tr := DefaultPooledTransport()
			if tr.MaxIdleConnsPerHost != tc.expected {
				t.Fatalf("DefaultPooledTransport().MaxIdleConnsPerHost = %d; want %d with GOMAXPROCS=%d",
					tr.MaxIdleConnsPerHost, tc.expected, tc.setProcs)
			}
		})
	}
}

func TestDefaultTransport(t *testing.T) {
	tr := DefaultTransport()
	if !tr.DisableKeepAlives {
		t.Fatal("expected DisableKeepAlives to be true")
	}
	if tr.MaxIdleConnsPerHost != -1 {
		t.Fatalf("expected MaxIdleConnsPerHost to be -1, got %d", tr.MaxIdleConnsPerHost)
	}
}

func TestDefaultPooledTransport_Configuration(t *testing.T) {
	tr := DefaultPooledTransport()
	if tr.DisableKeepAlives {
		t.Fatal("expected DisableKeepAlives to be false")
	}
	if tr.MaxIdleConns != 100 {
		t.Fatalf("expected MaxIdleConns to be 100, got %d", tr.MaxIdleConns)
	}
	if tr.IdleConnTimeout != 90*time.Second {
		t.Fatalf("expected IdleConnTimeout to be 90s, got %v", tr.IdleConnTimeout)
	}
	if tr.TLSHandshakeTimeout != 10*time.Second {
		t.Fatalf("expected TLSHandshakeTimeout to be 10s, got %v", tr.TLSHandshakeTimeout)
	}
	if tr.ExpectContinueTimeout != 1*time.Second {
		t.Fatalf("expected ExpectContinueTimeout to be 1s, got %v", tr.ExpectContinueTimeout)
	}
	if !tr.ForceAttemptHTTP2 {
		t.Fatal("expected ForceAttemptHTTP2 to be true")
	}
	if tr.MaxIdleConnsPerHost < 4 {
		t.Fatalf("expected MaxIdleConnsPerHost >= 4, got %d", tr.MaxIdleConnsPerHost)
	}
	if tr.Proxy == nil {
		t.Fatal("expected non-nil Proxy")
	}
	if tr.DialContext == nil {
		t.Fatal("expected non-nil DialContext")
	}
}

func TestDefaultClient(t *testing.T) {
	client := DefaultClient()
	if client == nil {
		t.Fatal("expected non-nil client")
	}

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected client.Transport to be *http.Transport, got %T", client.Transport)
	}
	if !tr.DisableKeepAlives {
		t.Fatal("expected DisableKeepAlives to be true on DefaultClient")
	}
	if tr.MaxIdleConnsPerHost != -1 {
		t.Fatalf("expected MaxIdleConnsPerHost to be -1, got %d", tr.MaxIdleConnsPerHost)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "ok"); err != nil {
			t.Errorf("unexpected error writing response: %v", err)
		}
	}))
	defer ts.Close()

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error making request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("unexpected error closing response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error reading body: %v", err)
	}
	if string(body) != "ok" {
		t.Fatalf("expected response body 'ok', got %q", string(body))
	}
}

func TestDefaultPooledClient(t *testing.T) {
	client := DefaultPooledClient()
	if client == nil {
		t.Fatal("expected non-nil client")
	}

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected client.Transport to be *http.Transport, got %T", client.Transport)
	}
	if tr.DisableKeepAlives {
		t.Fatal("expected DisableKeepAlives to be false on DefaultPooledClient")
	}
	if tr.MaxIdleConnsPerHost < 4 {
		t.Fatalf("expected MaxIdleConnsPerHost >= 4, got %d", tr.MaxIdleConnsPerHost)
	}

	requestsHandled := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestsHandled++
		if _, err := fmt.Fprint(w, "pooled-ok"); err != nil {
			t.Errorf("unexpected error writing response: %v", err)
		}
	}))
	defer ts.Close()

	for i := 0; i < 3; i++ {
		resp, err := client.Get(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error on request %d: %v", i, err)
		}
		body, err := io.ReadAll(resp.Body)
		if closeErr := resp.Body.Close(); closeErr != nil {
			t.Fatalf("unexpected error closing body on request %d: %v", i, closeErr)
		}
		if err != nil {
			t.Fatalf("unexpected error reading body on request %d: %v", i, err)
		}
		if string(body) != "pooled-ok" {
			t.Fatalf("expected response 'pooled-ok', got %q", string(body))
		}
	}

	if requestsHandled != 3 {
		t.Fatalf("expected 3 requests handled, got %d", requestsHandled)
	}
}
