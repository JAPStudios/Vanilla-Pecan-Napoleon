package doh

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExchangePostsDNSMessage(t *testing.T) {
	query := []byte{0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != dnsMessageMediaType {
			t.Errorf("Content-Type = %q", got)
		}
		if got := r.Header.Get("Accept"); got != dnsMessageMediaType {
			t.Errorf("Accept = %q", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		if !bytes.Equal(body, query) {
			t.Errorf("request body = %v, want %v", body, query)
		}
		w.Header().Set("Content-Type", dnsMessageMediaType)
		_, _ = w.Write(query)
	}))
	defer server.Close()

	resolver, err := New(server.URL, server.Client())
	if err != nil {
		t.Fatal(err)
	}
	got, err := resolver.Exchange(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, query) {
		t.Fatalf("response = %v, want %v", got, query)
	}
}

func TestNewRejectsUnsafeEndpoint(t *testing.T) {
	for _, endpoint := range []string{
		"http://resolver.example/dns-query",
		"https://user:password@resolver.example/dns-query",
		"https://resolver.example/dns-query#fragment",
		"/dns-query",
	} {
		if _, err := New(endpoint, nil); err == nil {
			t.Errorf("New(%q) succeeded; want error", endpoint)
		}
	}
}

func TestExchangeRejectsInvalidQuery(t *testing.T) {
	resolver, err := New("https://resolver.example/dns-query", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, query := range [][]byte{nil, make([]byte, 11), make([]byte, maxDNSMessageSize+1)} {
		if _, err := resolver.Exchange(context.Background(), query); err == nil {
			t.Errorf("Exchange with %d-byte query succeeded; want error", len(query))
		}
	}
}

func TestExchangeRejectsInvalidResponses(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		contentType string
		body        []byte
	}{
		{name: "HTTP error", status: http.StatusBadGateway, contentType: dnsMessageMediaType, body: make([]byte, 12)},
		{name: "wrong media type", status: http.StatusOK, contentType: "text/plain", body: make([]byte, 12)},
		{name: "short message", status: http.StatusOK, contentType: dnsMessageMediaType, body: make([]byte, 11)},
		{name: "oversized message", status: http.StatusOK, contentType: dnsMessageMediaType, body: make([]byte, maxDNSMessageSize+1)},
	}
	query := make([]byte, 12)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.status)
				_, _ = w.Write(tc.body)
			}))
			defer server.Close()
			resolver, err := New(server.URL, server.Client())
			if err != nil {
				t.Fatal(err)
			}
			if _, err := resolver.Exchange(context.Background(), query); err == nil {
				t.Fatal("Exchange succeeded; want error")
			}
		})
	}
}

func TestExchangePropagatesCancellation(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
	}))
	defer server.Close()
	resolver, err := New(server.URL, server.Client())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := resolver.Exchange(ctx, make([]byte, 12))
		done <- err
	}()
	<-started
	cancel()
	if err := <-done; err == nil || !strings.Contains(err.Error(), "DoH") {
		t.Fatalf("Exchange error = %v, want DoH cancellation error", err)
	}
}
