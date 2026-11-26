package doh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
)

const (
	dnsMessageMediaType = "application/dns-message"
	maxDNSMessageSize   = 65535
)

// Resolver sends DNS wire-format messages to an HTTPS DNS-over-HTTPS endpoint.
// It does not select an upstream, install system DNS settings, or log queries.
type Resolver struct {
	endpoint string
	client   *http.Client
}

// New validates an HTTPS resolver endpoint. A nil client uses http.DefaultClient.
func New(endpoint string, client *http.Client) (*Resolver, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse DoH endpoint: %w", err)
	}
	if parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.Fragment != "" {
		return nil, fmt.Errorf("DoH endpoint must be an HTTPS URL without user information or a fragment")
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &Resolver{endpoint: parsed.String(), client: client}, nil
}

// Exchange sends one DNS wire-format query using the RFC 8484 POST method.
func (r *Resolver) Exchange(ctx context.Context, query []byte) ([]byte, error) {
	if r == nil || r.client == nil || r.endpoint == "" {
		return nil, fmt.Errorf("DoH resolver is not initialized")
	}
	if ctx == nil {
		return nil, fmt.Errorf("request context is required")
	}
	if len(query) < 12 || len(query) > maxDNSMessageSize {
		return nil, fmt.Errorf("DNS query must contain between 12 and %d bytes", maxDNSMessageSize)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint, bytes.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("create DoH request: %w", err)
	}
	req.Header.Set("Content-Type", dnsMessageMediaType)
	req.Header.Set("Accept", dnsMessageMediaType)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send DoH request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DoH endpoint returned HTTP status %d", resp.StatusCode)
	}
	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || contentType != dnsMessageMediaType {
		return nil, fmt.Errorf("DoH endpoint returned an invalid DNS message content type")
	}

	answer, err := io.ReadAll(io.LimitReader(resp.Body, maxDNSMessageSize+1))
	if err != nil {
		return nil, fmt.Errorf("read DoH response: %w", err)
	}
	if len(answer) < 12 || len(answer) > maxDNSMessageSize {
		return nil, fmt.Errorf("DoH response must contain between 12 and %d bytes", maxDNSMessageSize)
	}
	return answer, nil
}
