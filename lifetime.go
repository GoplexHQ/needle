package needle

type Lifetime string

const (
	Transient   Lifetime = "TRANSIENT"    // A new instance of the dependency is created each time it is requested.
	Scoped      Lifetime = "SCOPED"       // A single instance of the dependency is created per scope, such as per web request or session.
	ThreadLocal Lifetime = "THREAD_LOCAL" // A single instance of the dependency is created per thread.
	Pooled      Lifetime = "POOLED"       // Instances are maintained in a pool and reused to manage the overhead of creating and destroying instances frequently.
	Singleton   Lifetime = "SINGLETON"    // A single instance of the dependency is created and shared throughout the entire lifetime of the application.
)

func (Lifetime) Values() []Lifetime {
	return []Lifetime{
		Transient,
		Scoped,
		ThreadLocal,
		Pooled,
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
