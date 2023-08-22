//go:build windows
// +build windows

package debug

import "github.com/wangweihong/eazycloud/pkg/log"

func installSignalHandler(outputDir string) {
	log.Warnf("system don't support runtime debug")
}
