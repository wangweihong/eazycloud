package libkubernetes

import "time"

type callInfo struct {
	timeout time.Duration
}

type Option func(*callInfo)

func TimeoutOption(timeout time.Duration) Option {
	return func(c *callInfo) {
		if timeout < 0 {
			return
		}
		c.timeout = timeout
	}
}
