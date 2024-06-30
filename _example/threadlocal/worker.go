package main

import (
	"log/slog"

	"github.com/goplexhq/needle"
)

type Worker struct {
	id int
}

func (w *Worker) Work(start, end int) {
	calculator, err := needle.Resolve[Calculator]()
	if err != nil {
		panic(err)
	}

	calculator.CalculateSum(start, end)

	slog.Info("Worker calculated sum of numbers",
		"worker", w.id,
		"calculator", calculator.id,
		"start", start,
		"end", end)
}

func (w *Worker) Result() int {
	calculator, err := needle.Resolve[Calculator]()
	if err != nil {
		panic(err)
	}

	return calculator.sum
}
