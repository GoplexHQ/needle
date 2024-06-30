package needle

// Lifetime defines the lifecycle of a registered service.
type Lifetime string

const (
	Transient   Lifetime = "TRANSIENT"    // A new instance of the dependency is created each time it is requested.
	Scoped      Lifetime = "SCOPED"       // A single instance of the dependency is created per scope, i.e. web request.
	ThreadLocal Lifetime = "THREAD_LOCAL" // A single instance of the dependency is created per thread.
	Singleton   Lifetime = "SINGLETON"    // A single instance is created and shared for the application's entire lifetime.
)

// Values returns all possible values of Lifetime.
func (Lifetime) Values() []Lifetime {
	return []Lifetime{
		Transient,
		Scoped,
		ThreadLocal,
		Singleton,
	}
}

// Valid checks if the Lifetime value is valid.
func (lifetime Lifetime) Valid() bool {
	for _, name := range lifetime.Values() {
		if lifetime == name {
			return true
		}
	}

	return false
}

// String returns the string representation of the Lifetime.
func (lifetime Lifetime) String() string {
	return string(lifetime)
}
