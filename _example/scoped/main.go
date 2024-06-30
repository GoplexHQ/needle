package main

import (
	"fmt"
	"log"

	"github.com/goplexhq/needle"
)

type Config struct {
	ConnectionString string
}

type Database struct {
	Config *Config `needle:"inject"`
}

const (
	scopeA = "Scope A"
	scopeB = "Scope B"
)

func main() {
	registerServices()
	runWorker(scopeA) // Output: [WORKER]: "Scope A" running at "a@scope.db"...
	runWorker(scopeB) // Output: [WORKER]: "Scope B" running at "b@scope.db"...
}

func registerServices() {
	cfgA := &Config{ConnectionString: "a@scope.db"}
	if err := needle.RegisterScopedInstance(cfgA, needle.WithScope(scopeA)); err != nil {
		log.Fatalf("Failed to register Config: %v", err)
	}

	cfgB := &Config{ConnectionString: "b@scope.db"}
	if err := needle.RegisterScopedInstance(cfgB, needle.WithScope(scopeB)); err != nil {
		log.Fatalf("Failed to register Config: %v", err)
	}
}

func runWorker(scope string) {
	var db Database

	if err := needle.InjectStructFields(&db, needle.WithScope(scope)); err != nil {
		log.Fatalf("Failed to inject struct fields: %v", err)
	}

	fmt.Printf("[WORKER]: %q running at %q...\n", scope, db.Config.ConnectionString)
}
