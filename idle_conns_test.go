package cleanhttp

import (
	"runtime"
	"testing"
)

func TestDefaultMaxIdleConnsPerHostFloor(t *testing.T) {
	got := defaultMaxIdleConnsPerHost()
	if got < 8 {
		t.Fatalf("got %d, want >= 8", got)
	}
	// When GOMAXPROCS is high, scale with 2x procs.
	n := runtime.GOMAXPROCS(0) * 2
	if n > 8 && got != n {
		t.Fatalf("got %d want %d", got, n)
	}
	tr := DefaultPooledTransport()
	if tr.MaxIdleConnsPerHost != got {
		t.Fatalf("transport=%d helper=%d", tr.MaxIdleConnsPerHost, got)
	}
}
