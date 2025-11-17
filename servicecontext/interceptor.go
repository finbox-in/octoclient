package servicecontext

import (
	"net/http"
)

// ContextInterceptor is an HTTP client interceptor that automatically injects ServiceContext into requests
type ContextInterceptor struct {
	client *http.Client
}

// NewContextInterceptor creates a new context interceptor
func NewContextInterceptor(client *http.Client) *ContextInterceptor {
	return &ContextInterceptor{
		client: client,
	}
}

// RoundTrip implements the http.RoundTripper interface to intercept requests
func (ci * ContextInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	// Extract ServiceContext from the request context
	serviceCtx := GetServiceContextFromGoContext(req.Context())

	// If ServiceContext exists, inject it into headers
	if serviceCtx != nil {
		InjectContext(serviceCtx, req.Header)
	}

	// Use the underlying client's transport or default transport
	transport := ci.client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	return transport.RoundTrip(req)
}

// WrapClientWithContextInterceptor wraps an HTTP client with context injection capability
func WrapClientWithContextInterceptor(client *http.Client) *http.Client {
	if client == nil {
		client = &http.Client{}
	}

	// Create a new client with the interceptor as transport
	wrappedClient := &http.Client{
		Timeout:   client.Timeout,
		Jar:       client.Jar,
		Transport: NewContextInterceptor(client),
	}

	return wrappedClient
}
