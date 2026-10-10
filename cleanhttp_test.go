// Copyright IBM Corp. 2015, 2026
// SPDX-License-Identifier: MPL-2.0

package cleanhttp

import (
	"net/http"
	"runtime"
	"testing"
	"time"
)

func TestDefaultTransport(t *testing.T) {
	tr := DefaultTransport()
	if tr == nil {
		t.Fatal("expected non-nil Transport")
	}

	if !tr.DisableKeepAlives {
		t.Errorf("expected DisableKeepAlives to be true, got %v", tr.DisableKeepAlives)
	}

	if tr.MaxIdleConnsPerHost != -1 {
		t.Errorf("expected MaxIdleConnsPerHost to be -1, got %d", tr.MaxIdleConnsPerHost)
	}

	if tr.MaxIdleConns != 100 {
		t.Errorf("expected MaxIdleConns to be 100, got %d", tr.MaxIdleConns)
	}

	if tr.IdleConnTimeout != 90*time.Second {
		t.Errorf("expected IdleConnTimeout to be 90s, got %v", tr.IdleConnTimeout)
	}

	if tr.TLSHandshakeTimeout != 10*time.Second {
		t.Errorf("expected TLSHandshakeTimeout to be 10s, got %v", tr.TLSHandshakeTimeout)
	}

	if tr.ExpectContinueTimeout != 1*time.Second {
		t.Errorf("expected ExpectContinueTimeout to be 1s, got %v", tr.ExpectContinueTimeout)
	}

	if !tr.ForceAttemptHTTP2 {
		t.Errorf("expected ForceAttemptHTTP2 to be true, got %v", tr.ForceAttemptHTTP2)
	}
}

func TestDefaultPooledTransport(t *testing.T) {
	orig := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(orig)

	tests := []struct {
		name      string
		procs     int
		wantConns int
	}{
		{
			name:      "single core floor",
			procs:     1,
			wantConns: 4,
		},
		{
			name:      "dual core floor",
			procs:     2,
			wantConns: 4,
		},
		{
			name:      "triple core boundary",
			procs:     3,
			wantConns: 4,
		},
		{
			name:      "quad core scaling",
			procs:     4,
			wantConns: 5,
		},
		{
			name:      "eight core scaling",
			procs:     8,
			wantConns: 9,
		},
		{
			name:      "sixteen core scaling",
			procs:     16,
			wantConns: 17,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runtime.GOMAXPROCS(tc.procs)
			tr := DefaultPooledTransport()
			if tr == nil {
				t.Fatal("expected non-nil Transport")
			}

			if tr.MaxIdleConnsPerHost != tc.wantConns {
				t.Errorf("GOMAXPROCS=%d: expected MaxIdleConnsPerHost=%d, got %d",
					tc.procs, tc.wantConns, tr.MaxIdleConnsPerHost)
			}

			if tr.DisableKeepAlives {
				t.Errorf("expected DisableKeepAlives to be false, got %v", tr.DisableKeepAlives)
			}

			if tr.MaxIdleConns != 100 {
				t.Errorf("expected MaxIdleConns to be 100, got %d", tr.MaxIdleConns)
			}

			if tr.IdleConnTimeout != 90*time.Second {
				t.Errorf("expected IdleConnTimeout to be 90s, got %v", tr.IdleConnTimeout)
			}

			if tr.TLSHandshakeTimeout != 10*time.Second {
				t.Errorf("expected TLSHandshakeTimeout to be 10s, got %v", tr.TLSHandshakeTimeout)
			}

			if tr.ExpectContinueTimeout != 1*time.Second {
				t.Errorf("expected ExpectContinueTimeout to be 1s, got %v", tr.ExpectContinueTimeout)
			}

			if !tr.ForceAttemptHTTP2 {
				t.Errorf("expected ForceAttemptHTTP2 to be true, got %v", tr.ForceAttemptHTTP2)
			}
		})
	}
}

func TestDefaultClient(t *testing.T) {
	client := DefaultClient()
	if client == nil {
		t.Fatal("expected non-nil Client")
	}

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected client.Transport to be *http.Transport, got %T", client.Transport)
	}

	if !tr.DisableKeepAlives {
		t.Errorf("expected DisableKeepAlives to be true, got %v", tr.DisableKeepAlives)
	}

	if tr.MaxIdleConnsPerHost != -1 {
		t.Errorf("expected MaxIdleConnsPerHost to be -1, got %d", tr.MaxIdleConnsPerHost)
	}
}

func TestDefaultPooledClient(t *testing.T) {
	orig := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(orig)

	runtime.GOMAXPROCS(2)
	client := DefaultPooledClient()
	if client == nil {
		t.Fatal("expected non-nil Client")
	}

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected client.Transport to be *http.Transport, got %T", client.Transport)
	}

	if tr.DisableKeepAlives {
		t.Errorf("expected DisableKeepAlives to be false, got %v", tr.DisableKeepAlives)
	}

	if tr.MaxIdleConnsPerHost != 4 {
		t.Errorf("expected MaxIdleConnsPerHost to be 4 on 2 cores, got %d", tr.MaxIdleConnsPerHost)
	}
}
