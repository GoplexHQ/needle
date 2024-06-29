package needle

import "sync"

var (
	globalStore *Store
	once        sync.Once
)

// InitGlobalStore initializes the global store if it has not been initialized already.
// This function is thread-safe and ensures that the store is only initialized once.
func InitGlobalStore() {
	once.Do(func() {
		globalStore = NewStore()
	})
}

// ensureGlobalStoreInitialized ensures that the global store is initialized.
// This function should be called before any operations that require the global store.
func ensureGlobalStoreInitialized() {
	if globalStore == nil {
		InitGlobalStore()
	}
}

// RegisteredServices returns a list of names of all services registered in the global store.
func RegisteredServices() []string {
	ensureGlobalStoreInitialized()

	return globalStore.RegisteredServices()
}

// Reset clears all entries in the global store.
func Reset() {
	ensureGlobalStoreInitialized()

	globalStore.Reset()
}
