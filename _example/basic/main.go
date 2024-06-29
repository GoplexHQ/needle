package main

import (
	"fmt"
	"log"

	"github.com/goplexhq/needle"
)

type Service struct{}

func (s *Service) SayHi(name string) {
	fmt.Println("Hello,", name)
}

type App struct {
	service *Service `needle:"inject"`
}

func main() {
	if err := needle.Register[Service](needle.Transient); err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	var app App
	if err := needle.InjectStructFields(&app); err != nil {
		log.Fatalf("failed to inject struct fields: %v", err)
	}

	app.service.SayHi("Needle!") // Output: Hello, Needle!
}
