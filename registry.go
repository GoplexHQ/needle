package needle

import (
	"reflect"
	"sync"
)

// Registry represents a thread-safe registry for storing service instances and their metadata.
type Registry struct {
	registeredServices  map[string]serviceEntry
	transientServices   map[string]reflect.Value
	scopedServices      map[string]map[string]reflect.Value
	threadLocalServices map[string]map[string]reflect.Value
	singletonServices   map[string]reflect.Value
	lock                sync.RWMutex
}

// NewRegistry creates and returns a new instance of Registry.
func NewRegistry() *Registry {
	return &Registry{ //nolint:exhaustruct
		registeredServices:  make(map[string]serviceEntry),
		transientServices:   make(map[string]reflect.Value),
		scopedServices:      make(map[string]map[string]reflect.Value),
		threadLocalServices: make(map[string]map[string]reflect.Value),
		singletonServices:   make(map[string]reflect.Value),
	}
}

// set adds or updates a service entry in the registry.
func (r *Registry) set(name string, lifetime Lifetime, value reflect.Value, options *ResolutionOptions) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.registeredServices[name] = serviceEntry{
		name:     name,
		lifetime: lifetime,
		value:    nil,
	}

	switch lifetime {
	case Transient:
		r.transientServices[name] = value
	case Scoped:
		if r.scopedServices[options.scope] == nil {
			r.scopedServices[options.scope] = make(map[string]reflect.Value)
		}

		r.scopedServices[options.scope][name] = value
	case ThreadLocal:
		if r.threadLocalServices[options.threadID] == nil {
			r.threadLocalServices[options.threadID] = make(map[string]reflect.Value)
		}

		r.threadLocalServices[options.threadID][name] = value
	case Singleton:
		r.singletonServices[name] = value
	}
}

// get retrieves a service entry from the registry by name.
// Returns the entry and a boolean indicating whether the entry was found.
func (r *Registry) get(name string, options *ResolutionOptions) (serviceEntry, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	entry, found := r.registeredServices[name]
	if !found {
		return entry, false
	}

	var (
		value  reflect.Value
		exists bool
	)

	switch entry.lifetime {
	case Transient:
		value, exists = r.transientServices[name]
	case Scoped:
		scope, scopeFound := r.scopedServices[options.scope]
		if !scopeFound {
			return entry, false
		}

		value, exists = scope[name]
	case ThreadLocal:
		thread, threadFound := r.threadLocalServices[options.threadID]
		if !threadFound {
			return entry, false
		}

		value, exists = thread[name]
	case Singleton:
		value, exists = r.singletonServices[name]
	}

	if !exists {
		return entry, false
	}

	return entry.withValue(&value), true
}

// has checks if a service entry exists in the registry by name.
func (r *Registry) has(name string) (serviceEntry, bool) {
	r.lock.RLock()
	entry, found := r.registeredServices[name]
	r.lock.RUnlock()

	return entry, found
}

func (r *Registry) hasScoped(scope, name string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	s, scopeExists := r.scopedServices[scope]
	if !scopeExists {
		return false
	}

	_, found := s[name]

	return found
}

func (r *Registry) hasThreadLocal(thread, name string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	t, threadExists := r.threadLocalServices[thread]
	if !threadExists {
		return false
	}

	_, found := t[name]

	return found
}

// RegisteredServices returns a list of names of all registered services.
// Service names are registryd in the following form "<pkg>.<service>".
func (r *Registry) RegisteredServices() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	names := make([]string, 0, len(r.registeredServices))
	for name := range r.registeredServices {
		names = append(names, name)
	}

	return names
}

// Reset clears all entries in the registry.
func (r *Registry) Reset() {
	r.lock.Lock()
	defer r.lock.Unlock()

	clear(r.registeredServices)
	clear(r.transientServices)
	clear(r.scopedServices)
	clear(r.threadLocalServices)
	clear(r.singletonServices)
}
