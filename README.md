# Needle :syringe:

Needle is a lightweight dependency injection framework for Go, designed to make it easy to manage dependencies in your
applications. It supports various lifetimes such as singleton, scoped, thread-local, and transient, offering flexibility
and control over the lifecycle of your services.

## Features

- **Flexible Lifetimes**: Manage dependencies with different lifetimes: Singleton, Scoped, ThreadLocal, and Transient.
- **Thread-Safety**: Ensure thread-safety with built-in synchronization mechanisms.
- **Optional Configuration**: Customize resolution and registration with optional scope and thread ID settings.
- **Reflection-Based Injection**: Leverage reflection to dynamically resolve and inject dependencies.

## Installation

To install Needle, run:

```sh
go get github.com/goplexhq/needle
```

## Usage

### Registering Services

#### Basic Registration

Register a service with a singleton lifetime:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyService struct{}

func main() {
	err := needle.Register[MyService](needle.Singleton)
	if err != nil {
		fmt.Println("Error registering service:", err)
	}
}
```

#### Registration with Scope

Register a service with a scoped lifetime:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyService struct{}

func main() {
	err := needle.Register[MyService](needle.Scoped, needle.WithScope("request1"))
	if err != nil {
		fmt.Println("Error registering scoped service:", err)
	}
}
```

#### Thread-Local Registration

Register a service with a thread-local lifetime:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyService struct{}

func main() {
	err := needle.Register[MyService](needle.ThreadLocal) // defaults to current goroutine.
	if err != nil {
		fmt.Println("Error registering thread-local service:", err)
	}
}
```

### Resolving Services

#### Basic Resolution

Resolve a singleton service:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyService struct{}

func main() {
	val, err := needle.Resolve[MyService]()
	if err != nil {
		fmt.Println("Error resolving service:", err)
	} else {
		fmt.Println("Resolved service:", val)
	}
}
```

#### Resolution with Scope

Resolve a scoped service:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyService struct{}

func main() {
	val, err := needle.Resolve[MyService](needle.WithScope("request1"))
	if err != nil {
		fmt.Println("Error resolving scoped service:", err)
	} else {
		fmt.Println("Resolved service:", val)
	}
}
```

#### Resolution with Thread ID

Resolve a thread-local service:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyService struct{}

func main() {
	val, err := needle.Resolve[MyService](needle.WithThreadID("thread1"))
	// or
	// needle.Resolve[MyService]() to resolve the service registered in the current goroutine.
	if err != nil {
		fmt.Println("Error resolving thread-local service:", err)
	} else {
		fmt.Println("Resolved service:", val)
	}
}
```

### Injecting Dependencies

#### Injecting into Struct Fields

Annotate struct fields with `needle:"inject"` and inject dependencies using the global registry:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyDependency struct{}

type MyStruct struct {
	Dep *MyDependency `needle:"inject"` // field must be a pointer to struct
}

func main() {
	var myStruct MyStruct
	err := needle.InjectStructFields(&myStruct)
	if err != nil {
		fmt.Println("Error injecting dependencies:", err)
	} else {
		fmt.Println("Injected struct:", myStruct)
	}
}
```

#### Injecting with Scope and Thread ID

Inject dependencies with optional scope and thread ID settings:

```go
package main

import (
	"fmt"
	"github.com/goplexhq/needle"
)

type MyDependency struct{}

type MyStruct struct {
	Dep *MyDependency `needle:"inject"`
}

func main() {
	var myStruct MyStruct
	err := needle.InjectStructFields(&myStruct, needle.WithScope("request1"))
	if err != nil {
		fmt.Println("Error injecting dependencies with scope:", err)
	} else {
		fmt.Println("Injected struct with scope:", myStruct)
	}

	err = needle.InjectStructFields(&myStruct, needle.WithThreadID("thread1"))
	if err != nil {
		fmt.Println("Error injecting dependencies with thread ID:", err)
	} else {
		fmt.Println("Injected struct with thread ID:", myStruct)
	}
}
```

## API Reference

### Functions

- #### `InitGlobalRegistry()`

  Initializes the global registry if it hasn't been initialized already.

- #### `RegisteredServices() []string`

  Returns a list of names of all services registered in the global registry.

- #### `Reset()`

  Clears all entries in the global registry.

- #### `Register[T any](lifetime Lifetime, optFuncs ...ResolutionOptionFunc) error`

  Registers a type with the specified lifetime to the global registry.

- #### `RegisterSingletonInstance[T any](val *T) error`

  Registers a pre-initialized singleton instance to the global registry.

- #### `RegisterToRegistry[T any](registry *Registry, lifetime Lifetime, optFuncs ...ResolutionOptionFunc) error`

  Registers a type with the specified lifetime to the given registry.

- #### `RegisterSingletonInstanceToRegistry[T any](registry *Registry, val *T) error`

  Registers a pre-initialized singleton instance to the given registry.

- #### `RegisterInstance[T any](lifetime Lifetime, val *T, optFuncs ...ResolutionOptionFunc) error`

  Registers a pre-initialized instance with the specified lifetime to the global registry.

- #### `RegisterInstanceToRegistry[T any](reg *Registry, lifetime Lifetime, val *T, optFns ...ResolutionOptionFunc) error`

  Registers a pre-initialized instance with the specified lifetime to the given registry.

- #### `RegisterScopedInstance[T any](val *T, optFuncs ...ResolutionOptionFunc) error`

  Registers a pre-initialized scoped instance to the global registry.

- #### `RegisterScopedInstanceToRegistry[T any](registry *Registry, val *T, optFuncs ...ResolutionOptionFunc) error`

  Registers a pre-initialized scoped instance to the given registry.

- #### `RegisterThreadLocalInstance[T any](val *T, optFuncs ...ResolutionOptionFunc) error`

  Registers a pre-initialized thread-local instance to the global registry.

- #### `RegisterThreadLocalInstanceToRegistry[T any](registry *Registry, val *T, optFuncs ...ResolutionOptionFunc) error`

  Registers a pre-initialized thread-local instance to the given registry.

- #### `Resolve[T any](optFuncs ...ResolutionOptionFunc) (*T, error)`

  Resolves an instance of the specified type from the global registry.

- #### `ResolveFromRegistry[T any](registry *Registry, optFuncs ...ResolutionOptionFunc) (*T, error)`

  Resolves an instance of the specified type from the given registry.

- #### `InjectStructFields[Dest any](dest *Dest, optFuncs ...ResolutionOptionFunc) error`

  Injects dependencies into the fields of a struct using the global registry.

- #### `InjectStructFieldsFromRegistry[Dest any](registry *Registry, dest *Dest, optFuncs ...ResolutionOptionFunc) error`

  Injects dependencies into the fields of a struct using the specified registry.

### Types

- #### `type Registry struct{}`

  A thread-safe registry for storing service instances and their metadata.

- #### `type Lifetime string`

  Represents the lifetime of a service. Supported lifetimes:

    - `Transient`
    - `Scoped`
    - `ThreadLocal`
    - `Singleton`

### Optional Configuration Functions

- #### `WithScope(scope string) ResolutionOptionFunc`

  Sets a scope for resolving scoped dependencies. Required when the lifetime is Scoped.

- #### `WithThreadID(threadID string) ResolutionOptionFunc`

  Sets a thread ID for resolving thread-local dependencies. Optional and defaults to the current goroutine ID if not
  provided and the lifetime is ThreadLocal.

### Errors

- #### `ErrRegistered`

  Indicates that the service is already registered in the registry.

- #### `ErrNotRegistered`

  Indicates that the service is not registered in the registry.

- #### `ErrInvalidServiceType`

  Indicates that the service type is invalid (must be a struct type).

- #### `ErrInvalidDestType`

  Indicates that the destination type is invalid (expected a struct type).

- #### `ErrServiceTypeMismatch`

  Indicates that the resolved service type does not match the expected type.

- #### `ErrFieldPtr`

  Indicates that an injectable field is not a pointer.

- #### `ErrResolveField`

  Indicates that the framework is unable to resolve a service for a field.

- #### `ErrEmptyScope`

  Indicates that a scope is required but not provided.

- #### `ErrTransientInstance`

  Indicates that a transient lifetime does not support pre-initialized instances.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue on GitHub.

## License

Needle is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
