// Package graceful coordinates graceful program shutdown.
package graceful

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// Done is a channel that is closed when program termination
// has been requested.
var Done chan struct{}

var (
	mu                sync.Mutex
	shutdownRequested bool
	shutdownCallbacks []func()
)

// OnShutdown registers a function callback for when a program
// shutdown request has been received.
func OnShutdown(callback func()) {
	mu.Lock()
	shutdownCallbacks = append(shutdownCallbacks, callback)
	mu.Unlock()
}

// Shutdown requests program shutdown.
func Shutdown() {
	requestShutdown()
}

// requestShutdown performs a one-time shutdown sequence.
func requestShutdown() {
	for _, callback := range initateShutdown() {
		callback()
	}
}

// initiateShutdown starts a shutdown and returns the
// list of callbacks. It does not actually call the callbacks,
// because this could potentially cause a deadlock of the callback
// tries to shutdown the program.
func initateShutdown() []func() {
	mu.Lock()
	defer mu.Unlock()
	if shutdownRequested {
		return nil
	}
	shutdownRequested = true
	log.Println("shutdown requested")
	close(Done)
	callbacks := shutdownCallbacks
	shutdownCallbacks = nil
	return callbacks
}

func init() {
	Done = make(chan struct{})

	ch := make(chan os.Signal)
	signal.Notify(ch, termSignals...)

	go func() {
		for sig := range ch {
			log.Printf("signal caught: %s", sig)
			requestShutdown()
			return
		}
	}()
}
