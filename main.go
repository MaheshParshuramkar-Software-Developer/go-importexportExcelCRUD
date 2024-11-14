package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-importexportExcelCRUD/initializers"
	"github.com/go-importexportExcelCRUD/routes"
)

func main() {
	cfg := initializers.LoadConfig()
	initializers.InitializeConnections()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-notify
		fmt.Println("Received Graceful shutdown signal")
		initializers.StopServices()
	}()

	r := routes.SetupRouter(cfg)
	initializers.StartServer(r)

	// wait for graceful shutdown
	wg.Wait()

	// wait for queued task to complete
	time.Sleep(time.Second * 2)
	fmt.Println("Exiting ...")
}
