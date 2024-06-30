# Needle :syringe:

Needle is a lightweight and flexible dependency injection framework for Go. It helps you manage dependencies and service
lifetimes in a clean and efficient manner.

## Features

- Simple API for registering and resolving dependencies
- Supports multiple lifetimes: Singleton, Transient, Scoped, and ThreadLocal
- Thread-safe service registry
- Easy integration with existing Go projects

## Installation

Install Needle using `go get`:

```bash
go get github.com/goplexhq/needle
```

## Usage

### Registering Services

You can register services with different lifetimes using the `Register` and `RegisterSingletonInstance` functions.

#### Register a Singleton

```go
type MyService struct{}

func main() {
    err := needle.Register[MyService](needle.Singleton)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }
}
```

#### Register a Transient Service

```go
type MyService struct{}

func main() {
    err := needle.Register[MyService](needle.Transient)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }
}
```

#### Register a Pre-initialized Instance

```go
type MyService struct {
    Name string
}

func main() {
    instance := &MyService{Name: "myService"}
    err := needle.RegisterSingletonInstance(instance)
    if err != nil {
        log.Fatalf("failed to register service instance: %v", err)
    }
}
```

### Resolving Services

You can resolve services using the `Resolve` function.

```go
type MyService struct{}

func main() {
    err := needle.Register[MyService](needle.Singleton)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }
    
    service, err := needle.Resolve[MyService]()
    if err != nil {
        log.Fatalf("failed to resolve service: %v", err)
    }
    
    fmt.Println(service)
}
```

### Injecting Dependencies

Needle can inject dependencies into struct fields using the `needle:"inject"` tag.

```go
type Dep struct {
    Name string
}

type MyService struct {
    Dep *Dep `needle:"inject"`
}

func main() {
    err := needle.RegisterSingletonInstance(&Dep{Name: "myDep"})
    if err != nil {
        log.Fatalf("failed to register dependency: %v", err)
    }

    var myService MyService
    err = needle.InjectStructFields(&myService)
    if err != nil {
        log.Fatalf("failed to inject struct fields: %v", err)
    }

    fmt.Println(myService.Dep.Name) // Output: myDep
}
```

> **Note:** Fields that need dependency injection should be annotated with the `needle:"inject"` tag.
> This tells Needle to inject the corresponding dependency into the field.

### Custom Registry

You can use a custom registry to manage dependencies separately from the global registry.

```go
type MyService struct{}

func main() {
    registry := needle.NewRegistry()

    err := needle.RegisterToRegistry[MyService](registry, needle.Singleton)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }

    service, err := needle.ResolveFromRegistry[MyService](registry)
    if err != nil {
        log.Fatalf("failed to resolve service from registry: %v", err)
    }

    fmt.Println(service)
}
```

## API Reference

### Functions

#### `InitGlobalRegistry()`

Initializes the global registry if it hasn't been initialized already.

#### `RegisteredServices() []string`

Returns a list of names of all services registered in the global registry.

#### `Reset()`

Clears all entries in the global registry.

#### `Register[T any](lifetime Lifetime) error`

Registers a type with the specified lifetime to the global registry.

#### `RegisterSingletonInstance[T any](val *T) error`

Registers a pre-initialized singleton instance to the global registry.

#### `RegisterToRegistry[T any](registry *Registry, lifetime Lifetime) error`

Registers a type with the specified lifetime to the given registry.

#### `RegisterSingletonInstanceToRegistry[T any](registry *Registry, val *T) error`

Registers a pre-initialized singleton instance to the given registry.

#### `Resolve[T any]() (*T, error)`

Resolves an instance of the specified type from the global registry.

#### `ResolveFromRegistry[T any](registry *Registry) (*T, error)`

Resolves an instance of the specified type from the given registry.

#### `InjectStructFields[Dest any](dest *Dest) error`

Injects dependencies into the fields of a struct using the global registry.

#### `InjectStructFieldsFromRegistry[Dest any](registry *Registry, dest *Dest) error`

Injects dependencies into the fields of a struct using the specified registry.

### Types

#### `Registry`

A thread-safe registry for storing service instances and their metadata.

#### `Lifetime`

Represents the lifetime of a service. Supported lifetimes:

- `Transient`
- `Scoped`
- `ThreadLocal`
- `Singleton`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

## License

This project is licensed under the MIT License.
