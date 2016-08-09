// +build windows

package graceful

import "os"

// termSignals are the list of signals that will cause graceful program termination.
var termSignals = []os.Signal{os.Interrupt}
