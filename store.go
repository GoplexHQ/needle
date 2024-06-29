package needle

import "sync"

// Store represents a thread-safe registry for storing service instances and their metadata.
type Store struct {
	entries map[string]entry
	lock    sync.RWMutex
}

// NewStore creates and returns a new instance of Store.
func NewStore() *Store {
	return &Store{ //nolint:exhaustruct
		entries: make(map[string]entry),
	}
}

// set adds or updates a service entry in the store.
func (s *Store) set(name string, value entry) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.entries[name] = value
}

// get retrieves a service entry from the store by name.
// Returns the entry and a boolean indicating whether the entry was found.
func (s *Store) get(name string) (entry, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v, ok := s.entries[name]

	return v, ok
}

// has checks if a service entry exists in the store by name.
func (s *Store) has(name string) bool {
	_, ok := s.get(name)

	return ok
}

// RegisteredServices returns a list of names of all registered services.
// Service names are stored in the following form "<pkg>.<service>".
func (s *Store) RegisteredServices() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	names := make([]string, 0, len(s.entries))

	for n := range s.entries {
		names = append(names, n)
	}

	return names
}

// Reset clears all entries in the store.
func (s *Store) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	clear(s.entries)
}
