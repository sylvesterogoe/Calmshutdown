package calmshutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Start intercepts OS Interrupt and Terminate signals and waits 
// for a specified delay duration before shutting down the server.
// An http.ErrServerClosed error is passed onto the shutdownError 
// channel if the shutdown is unsuccessful.
func Start(shutdownError chan error, server *http.Server, delay time.Duration) {

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quit

	logger := log.New(os.Stdout, "Graceful shutdown: ", log.Ldate|log.Ltime)
	logger.Printf("%s signal was caught", signal.String())

	ctx, cancel := context.WithTimeout(context.Background(), delay)
	defer cancel()

	shutdownError <- server.Shutdown(ctx)
}

// StartAndAwaitGoroutines waits for all user defined goroutines
// that the application spins up before shutting down the server.
// It takes a global application-wide waitgroup(globalWG) that
// keeps track of all such goroutines
// An http.ErrServerClosed error is passed onto the shutdownError 
// channel if the shutdown is unsuccessful.
func StartAndAwaitGoroutines(shutdownError chan error, server *http.Server,
	delay time.Duration, globalWG *sync.WaitGroup) {

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quit

	logger := log.New(os.Stdout, "Graceful shutdown: ", log.Ldate|log.Ltime)
	logger.Printf("%s signal was caught", signal.String())

	ctx, cancel := context.WithTimeout(context.Background(), delay)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		shutdownError <- err
	}

	logger.Print("Completing background tasks")

	globalWG.Wait()
	
	shutdownError <- nil
}
