package main

import (
	"fmt"
	"log"

	"github.com/goplexhq/needle"
)

type Config struct {
	AppName string
	Version string
}

type Logger struct {
	Prefix string
}

func (l *Logger) Log(message string) {
	fmt.Printf("[%s] %s\n", l.Prefix, message)
}

type App struct {
	Config *Config `needle:"inject"`
	Logger *Logger `needle:"inject"`
}

func (a *App) Start() {
	a.Logger.Log(fmt.Sprintf("Starting %s v%s...", a.Config.AppName, a.Config.Version))
}

func main() {
	err := needle.RegisterSingletonInstance(&Config{AppName: "MyApp", Version: "1.0.0"})
	if err != nil {
		log.Fatalf("failed to register config: %v", err)
	}

	err = needle.RegisterSingletonInstance(&Logger{Prefix: "INFO"})
	if err != nil {
		log.Fatalf("failed to register logger: %v", err)
	}

	var app App

	err = needle.InjectStructFields(&app)
	if err != nil {
		log.Fatalf("failed to inject struct fields: %v", err)
	}

	app.Start() // Output: [INFO] Starting MyApp v1.0.0...
}
