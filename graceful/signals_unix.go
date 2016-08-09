// +build darwin linux

package graceful

import (
	"os"
	"syscall"
)

// termSignals are the list of signals that will cause graceful program termination.
var termSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
