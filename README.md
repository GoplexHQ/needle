# Needle :syringe:

Needle is a lightweight and flexible dependency injection framework for Go. It helps you manage dependencies and service
lifetimes in a clean and efficient manner.

## Features

- Simple API for registering and resolving dependencies
- Supports multiple lifetimes: Singleton, Transient, Scoped, ThreadLocal, and Pooled
- Thread-safe service registry
- Easy integration with existing Go projects

## Installation

Install Needle using `go get`:

```bash
go get github.com/goplexhq/needle
```

## Usage

### Registering Services

You can register services with different lifetimes using the `Register` and `RegisterInstance` functions.

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
    err := needle.RegisterInstance(instance)
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

Needle can inject dependencies into struct fields.

```go
type Dep struct {
    Name string
}

type MyService struct {
    Dep *Dep
}

func main() {
    err := needle.RegisterInstance(&Dep{Name: "myDep"})
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

### Custom Store

You can use a custom store to manage dependencies separately from the global store.

```go
type MyService struct{}

func main() {
    store := needle.NewStore()

    err := needle.RegisterToStore[MyService](store, needle.Singleton)
    if err != nil {
        log.Fatalf("failed to register service: %v", err)
    }

    service, err := needle.ResolveFromStore[MyService](store)
    if err != nil {
        log.Fatalf("failed to resolve service from store: %v", err)
    }

    fmt.Println(service)
}
```

## API Reference

### Functions

#### `InitGlobalStore()`

Initializes the global store if it hasn't been initialized already.

#### `RegisteredServices() []string`

Returns a list of names of all services registered in the global store.

#### `Reset()`

Clears all entries in the global store.

#### `Register[T any](lifetime Lifetime) error`

Registers a type with the specified lifetime to the global store.

#### `RegisterInstance[T any](val *T) error`

Registers a pre-initialized singleton instance to the global store.

#### `RegisterToStore[T any](store *Store, lifetime Lifetime) error`

Registers a type with the specified lifetime to the given store.

#### `RegisterInstanceToStore[T any](store *Store, val *T) error`

Registers a pre-initialized singleton instance to the given store.

#### `Resolve[T any]() (*T, error)`

Resolves an instance of the specified type from the global store.

#### `ResolveFromStore[T any](store *Store) (*T, error)`

Resolves an instance of the specified type from the given store.

#### `InjectStructFields[Dest any](dest *Dest) error`

Injects dependencies into the fields of a struct using the global store.

#### `InjectStructFieldsFromStore[Dest any](store *Store, dest *Dest) error`

Injects dependencies into the fields of a struct using the specified store.

### Types

#### `Store`

A thread-safe registry for storing service instances and their metadata.

#### `Lifetime`

Represents the lifetime of a service. Supported lifetimes:

- `Transient`
- `Scoped`
- `ThreadLocal`
- `Pooled`
- `Singleton`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

## License

This project is licensed under the MIT License.
