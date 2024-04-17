package syncer_test

import (
	"testing"
	"time"

	"github.com/wangweihong/eazycloud/pkg/wait"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/workqueue"

	"github.com/wangweihong/eazycloud/pkg/syncer"
)

type Object struct {
	ID   string
	Name string
	Data string
}

func getObject(id string) *Object {
	return &Object{
		ID:   "123",
		Name: "123",
		Data: "213",
	}
}

func TestNewWorkequeueSyncer(t *testing.T) {
	Convey("onework", t, func() {
		stop := make(chan struct{}, 0)
		s := syncer.NewWorkequeueSyncer(func(key interface{}) error {
			id := key.(string)
			getObject(id)
			return nil
		}, workqueue.New(), 3*time.Second, 1, 3)
		s.Run(stop)
		s.Trigger("123")
		s.Trigger("123")

		go wait.Until(func() {
			s.Trigger("123")
		}, 300*time.Second, stop)

		go func() {
			select {
			case <-time.After(1 * time.Second):
				close(stop)
			}
		}()
		<-stop
	})
}
