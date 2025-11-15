package internalhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContextInterceptor(t *testing.T) {
	t.Run("should inject service context into request headers", func(t *testing.T) {
		// Create a test server to capture the request
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the context header is present
			contextHeader := r.Header.Get(ContextHeader)
			if contextHeader == "" {
				t.Error("Context header should not be empty")
			}

			// Extract and verify the context
			extracted, err := ExtractContext(r.Header)
			if err != nil {
				t.Fatalf("ExtractContext failed: %v", err)
			}

			if extracted.Attributes["test-key"] != "test-value" {
				t.Errorf("Expected attribute 'test-value', got '%v'", extracted.Attributes["test-key"])
			}

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		// Create service context and add to request context
		serviceCtx := NewServiceContext()
		serviceCtx.WithAttribute("test-key", "test-value")

		ctx := WithServiceContext(context.Background(), serviceCtx)

		// Create request with the context
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			t.Fatalf("NewRequestWithContext failed: %v", err)
		}

		// Create interceptor and execute the request
		client := WrapClientWithContextInterceptor(nil)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("client.Do failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
