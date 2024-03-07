package syncer

import (
	"time"

	"github.com/wangweihong/eazycloud/pkg/wait"
)

type Syncer struct {
	period     time.Duration
	syncAction func()
}

func New(
	internal time.Duration,
	action func(),
) *Syncer {
	return &Syncer{
		period:     internal,
		syncAction: action,
	}
}

func (u *Syncer) Run(stop <-chan struct{}) {
	go func() {
		wait.Until(u.syncAction, u.period, stop)
	}()
}
