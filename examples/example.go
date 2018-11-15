// +build ignore

package main

import (
        "flag"
	"fmt"
	"math/rand"
	"time"

	ui "github.com/gizak/termui"
	"github.com/raff/slots"
	"github.com/raff/progress"
	"github.com/stefantalpalaru/pool"
)

func main() {
        tasks := flag.Int("tasks", 10, "number of concurrent tasks")
        total := flag.Int("total", 100, "number of total runs")
        border := flag.Bool("border", true, "border/no border")

        flag.Parse()

	rand.Seed(time.Now().Unix())

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("<Resize>", func(e ui.Event) {
		payload := e.Payload.(ui.Resize)
		ui.Body.Width = payload.Width
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	p := progress.New(*tasks, *border)
	slots := slots.New(*tasks)

	tpool := pool.New(*tasks)
	tpool.Run()
	defer tpool.Stop()

	go func() {
		for c := 0; c < *total; c++ {
			tpool.Add(func(params ...interface{}) interface{} {
				n := params[0].(int)
				z := rand.Intn(10)

				item := slots.Take()

				for i := 0; i < z; i++ {
					if i > 0 {
						time.Sleep(time.Second)
					}

					p.Set(item, fmt.Sprintf("Task %v Sleep %v", n, z-i), progress.PercInt(i, z))
				}

				p.Set(item, fmt.Sprintf("Task %v Done!", n), 100)

				slots.Release(item)
				return nil
			}, c)
		}

		status := tpool.Status()
		fmt.Println(status.Submitted, "submitted jobs,", status.Running, "running,", status.Completed, "completed.")
		fmt.Println()
		tpool.Wait()
		ui.StopLoop()
	}()

	ui.Loop()
}
