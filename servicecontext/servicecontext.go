package servicecontext

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// ServiceContext represents the context data that will be propagated through service calls
type ServiceContext struct {
	Attributes map[string]any `json:"attributes"`
	Baggage    map[string]any `json:"baggage"`
}

// NewServiceContext creates a new ServiceContext with empty initilised maps
func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		Attributes: make(map[string]any),
		Baggage:    make(map[string]any),
	}
}

type ContextKeyType string

const (
	ContextHeader                    = "X-Service-Context"
	ServiceContextKey ContextKeyType = "service_context"
)

// WithAttribute adds an attribute to the service context
func (sc *ServiceContext) WithAttribute(key string, value any) {
	if sc.Attributes == nil {
		sc.Attributes = make(map[string]any)
	}
	sc.Attributes[key] = value
}

// WithBaggage adds baggage to the service context
func (sc *ServiceContext) WithBaggage(key, value string) {
	if sc.Baggage == nil {
		sc.Baggage = make(map[string]any)
	}
	sc.Baggage[key] = value
}

// InjectContext injects the service context into HTTP headers
func InjectContext(serviceCtx *ServiceContext, headers http.Header) error {
	if serviceCtx == nil {
		return nil
	}

	contextJSON, err := json.Marshal(serviceCtx)
	if err != nil {
		return err
	}

	headers.Set(ContextHeader, base64.StdEncoding.EncodeToString(contextJSON))
	return nil
}

// ExtractContext extracts the service context from HTTP headers
func ExtractContext(headers http.Header) (*ServiceContext, error) {
	contextHeader := headers.Get(ContextHeader)
	if contextHeader == "" {
		return NewServiceContext(), nil
	}

	decoded, err := base64.StdEncoding.DecodeString(contextHeader)
	if err != nil {
		return NewServiceContext(), err
	}

	var serviceCtx ServiceContext
	if err := json.Unmarshal(decoded, &serviceCtx); err != nil {
		return NewServiceContext(), err
	}

	return &serviceCtx, nil
}

func GetServiceContextFromGoContext(ctx context.Context) *ServiceContext {
	if serviceCtx, ok := ctx.Value(ServiceContextKey).(*ServiceContext); ok {
		return serviceCtx
	}

	return nil
}

// WithServiceContext adds ServiceContext to Go context
func WithServiceContext(ctx context.Context, serviceCtx *ServiceContext) context.Context {
	return context.WithValue(ctx, ServiceContextKey, serviceCtx)
}
