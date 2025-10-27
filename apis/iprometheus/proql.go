package iprometheus

import (
	"fmt"
	"time"
)

const (
	MetricTypeMatrix = "matrix"
	MetricTypeVector = "vector"
)

type ProqlOption struct {
	Level     string
	Node      string
	Namespace string
	PodName   string
	Metric    string
	Paras     []string

	Time      time.Time
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

func (m *ProqlOption) String() string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("level[%v],node[%v],namespace[%v],podname[%v]metric[%v]",
		m.Level, m.Node, m.Namespace, m.PodName, m.Metric)
}
