package needle

import "sync"

//nolint:gochecknoglobals
var (
	globalRegistry *Registry
	once           sync.Once
)

// InitGlobalRegistry initializes the global registry if it has not been initialized already.
// This function is thread-safe and ensures that the registry is only initialized once.
func InitGlobalRegistry() {
	once.Do(func() {
		globalRegistry = NewRegistry()
	})
}

// ensureGlobalRegistryInitialized ensures that the global registry is initialized.
// This function should be called before any operations that require the global registry.
func ensureGlobalRegistryInitialized() {
	if globalRegistry == nil {
		InitGlobalRegistry()
	}
}

// RegisteredServices returns a list of names of all services registered in the global registry.
func RegisteredServices() []string {
	ensureGlobalRegistryInitialized()

	return globalRegistry.RegisteredServices()
}

// Reset clears all entries in the global registry.
func Reset() {
	ensureGlobalRegistryInitialized()

	globalRegistry.Reset()
}
