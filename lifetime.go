package needle

type Lifetime string

const (
	Transient   Lifetime = "TRANSIENT"    // A new instance of the dependency is created each time it is requested.
	Scoped      Lifetime = "SCOPED"       // A single instance of the dependency is created per scope, i.e. web request.
	ThreadLocal Lifetime = "THREAD_LOCAL" // A single instance of the dependency is created per thread.
	Singleton   Lifetime = "SINGLETON"    // A single instance is created and shared for the application's entire lifetime.
)

func (Lifetime) Values() []Lifetime {
	return []Lifetime{
		Transient,
		Scoped,
		ThreadLocal,
		Singleton,
	}
}

func (lifetime Lifetime) Valid() bool {
	for _, name := range lifetime.Values() {
		if lifetime == name {
			return true
		}
	}

	return false
}

func (lifetime Lifetime) String() string {
	return string(lifetime)
}
