package golang

import (
	"fmt"
	"runtime"
	"time"
)

func TestTicker() {
	fmt.Println("1 goroutine", runtime.NumGoroutine())
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("return ticker loop")
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	fmt.Println("2 goroutine", runtime.NumGoroutine())
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	fmt.Println("Ticker stopped")
	done <- true

	fmt.Println("3 goroutine", runtime.NumGoroutine())
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("4 goroutine", runtime.NumGoroutine())
}
