package mocks

import (
	"github.com/jarcoal/httpmock"
)

// MockRegistrar defines the interface for mock registrars
type MockRegistrar interface {
	RegisterMocks()
	RegisterErrorMocks()
}

// Registry holds all registered mocks
type Registry struct {
	mocks map[string]MockRegistrar
}

// NewRegistry creates a new mock registry
func NewRegistry() *Registry {
	return &Registry{
		mocks: make(map[string]MockRegistrar),
	}
}

// Register adds a mock to the registry
func (r *Registry) Register(name string, mock MockRegistrar) {
	r.mocks[name] = mock
}

// ActivateMocks activates specific mocks or all if none specified
func (r *Registry) ActivateMocks(names ...string) {
	httpmock.Activate()

	// If no names provided, activate all
	if len(names) == 0 {
		for _, mock := range r.mocks {
			mock.RegisterMocks()
		}
		return
	}

	// Activate only specified mocks
	for _, name := range names {
		if mock, ok := r.mocks[name]; ok {
			mock.RegisterMocks()
		}
	}
}

// ActivateErrorMocks activates error mocks for specific resources or all if none specified
func (r *Registry) ActivateErrorMocks(names ...string) {
	httpmock.Activate()

	// If no names provided, activate all
	if len(names) == 0 {
		for _, mock := range r.mocks {
			mock.RegisterErrorMocks()
		}
		return
	}

	// Activate only specified mocks
	for _, name := range names {
		if mock, ok := r.mocks[name]; ok {
			mock.RegisterErrorMocks()
		}
	}
}

// DeactivateAndReset deactivates httpmock
func (r *Registry) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

// Global registry instance
var GlobalRegistry = NewRegistry()
