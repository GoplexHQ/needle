package main

import (
	"log/slog"
	"sync"

	"github.com/goplexhq/needle"
)

type App struct {
	wg         sync.WaitGroup
	mu         sync.Mutex
	start      int
	end        int
	numWorkers int
	total      int
}

func (app *App) runWorker(id int, start, end int) {
	defer app.wg.Done()

	if err := needle.RegisterThreadLocalInstance(&Calculator{id: id}); err != nil {
		panic(err)
	}

	worker := Worker{id: id}
	worker.Work(start, end)

	app.mu.Lock()
	app.total += worker.Result()
	app.mu.Unlock()
}

func (app *App) run() {
	app.wg.Add(app.numWorkers)

	chunkSize := app.end / app.numWorkers

	for id := range app.numWorkers {
		go func() {
			app.runWorker(id, id*chunkSize+1, id*chunkSize+chunkSize)
		}()
	}

	app.wg.Wait()

	slog.Info("All workers finished", "start", app.start, "end", app.end, "total_sum", app.total)
}

func main() {
	app := App{
		start:      1,
		end:        200,
		numWorkers: 4,
		total:      0,
	}
	app.run()
}
