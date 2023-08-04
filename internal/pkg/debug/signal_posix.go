package debug

import (
	"os"
	"syscall"
)

var debugSignals = []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2}
