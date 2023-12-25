package waitgroup

import (
	"github.com/wangweihong/eazycloud/pkg/util/randutil"
	"runtime/debug"

	"context"
	"fmt"

	"sync"
)

type WaitGroupRoutineFunc struct {
	Name string
	Ctx  context.Context
	Call func() WaitGroupResult //异步函数体
}

// WARN:
func NewWaitGroupHandleFunc(ctx context.Context, name string, call func() WaitGroupResult) WaitGroupRoutineFunc {
	if name == "" {
		name = "async-routine"
	}
	name = name + "-" + randutil.RandNumSets(6)
	return WaitGroupRoutineFunc{
		Name: name,
		Ctx:  ctx,
		Call: call,
	}
}

func NewWaitGroup(ctx context.Context) *Group {
	return &Group{
		ctx:     ctx,
		wg:      sync.WaitGroup{},
		results: make(map[string]WaitGroupResult),
		retLock: sync.Mutex{},
	}
}

// Group allows to start a group of goroutines and wait for their completion.
type Group struct {
	ctx     context.Context
	wg      sync.WaitGroup
	results map[string]WaitGroupResult
	retLock sync.Mutex
	debug   bool
}

func (g *Group) Debug() {
	g.debug = true
}

func (g *Group) Wait() {
	g.wg.Wait()
	g.PrintResults()
}

// Start starts f in a new goroutine in the group.
func (g *Group) Start(f WaitGroupRoutineFunc) {
	g.wg.Add(1)
	go func() {
		ret := NewWaitGroupResult(nil, nil)
		defer g.wg.Done()
		defer g.setResult(f.Name, &ret)
		defer g.handleWaitGroupCrash(&ret)
		ret = f.Call()
	}()
}

func (g *Group) setResult(name string, ret *WaitGroupResult) {
	g.retLock.Lock()
	defer g.retLock.Unlock()
	g.results[name] = *ret
}

func (g *Group) GetResults() map[string]WaitGroupResult {
	g.retLock.Lock()
	defer g.retLock.Unlock()
	return g.results
}

func (g *Group) GetResultFast() ( /*total*/ int /*success*/, int /*fail*/, int) {
	g.retLock.Lock()
	defer g.retLock.Unlock()

	total := len(g.results)
	var fail int
	var success int
	for _, v := range g.results {
		if v.Error != nil {
			fail += 1
		} else {
			success += 1
		}
	}
	return total, success, fail
}

func (g *Group) PrintResults() {
	g.retLock.Lock()
	defer g.retLock.Unlock()
	if g.debug {
		fmt.Printf("has start %v waitgroup routines\n", len(g.results))
		for k, v := range g.results {
			fmt.Printf("waitgroup routine %v, results:%v\n", k, v)
		}
	}
}

func (g *Group) handleWaitGroupCrash(st *WaitGroupResult) {
	if x := recover(); x != nil {
		st.Error = fmt.Errorf("runtime panic:%v, stack:%v", x, string(debug.Stack()))
	}
}

func (g *Group) ConvertResultToBatchOutput() BatchOutput {
	g.retLock.Lock()
	defer g.retLock.Unlock()

	var bo BatchOutput
	for _, v := range g.results {
		bo.Total += 1
		if v.Error != nil {
			bo.Fail += 1
		} else {
			bo.Success += 1
		}
		bo.Results = append(bo.Results, SetOutput(v.Data, v.Error))
	}
	return bo
}

type WaitGroupResult struct {
	Error error
	Data  interface{}
}

func NewWaitGroupResult(data interface{}, err error) WaitGroupResult {
	return WaitGroupResult{
		Data:  data,
		Error: err,
	}
}
