package debug

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wangweihong/eazycloud/pkg/log"
)

var debugHandler chan os.Signal

func SetupRuntimeDebugSignalHandler(outputDir string) {
	debugHandler = make(chan os.Signal, 2)

	signal.Notify(debugHandler, debugSignals...)

	go func() {
		for {
			sig := <-debugHandler
			switch sig {
			case syscall.SIGUSR1:
				log.Info("receive SIGUSR1 signal, start system prof collect")
				if err := StartProf(outputDir); err != nil {
					log.Warnf("collect system prof error:%v", err)
					continue
				}
				log.Infof("runtime data collect success, outdir:%v", outputDir)
			case syscall.SIGUSR2:
				log.Infof("receive SIGUSR2 signal, change dynamic debug to %v", !Dynamic)
				Dynamic = !Dynamic
			}
		}
	}()
}
