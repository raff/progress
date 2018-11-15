// +build ignore

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/raff/progress"
	"github.com/raff/slots"
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

	ui.Handle("q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("r", func(ui.Event) {
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("<Resize>", func(e ui.Event) {
		payload := e.Payload.(ui.Resize)
		ui.Body.Width = payload.Width
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	p := progress.New(*tasks, *border,
		progress.Header(2),
		progress.Messages(10))

	p.SetHeader("Example application\nPress `q` to quit")

	slots := slots.New(*tasks)

	tpool := pool.New(*tasks)
	tpool.Run()

	defer func() {
		ui.Clear()
		ui.Close()
		tpool.Stop()

		fmt.Println(strings.Join(p.Messages(), "\n"))
	}()

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
				p.AddMessage(fmt.Sprintf("Task %v Done!", n))

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
