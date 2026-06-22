package iprometheus

import (
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type Config struct {
	Address string `json:"address"`
}

type Metadata struct {
	Metric string `json:"metric"`
	Type   string `json:"type"`
	Help   string `json:"help"`
}

type MetricOne struct {
	MetricName string  `json:"metric_name"`
	Series     []Point `json:"series"`
	Sample     *Point  `json:"sample"`
	Error      string  `json:"error"`
}

type Metric struct {
	MetricData
	MetricName string `json:"metric_name"`
	Error      string `json:"error"`
}

func (metric *Metric) ToMetricOnes() []*MetricOne {
	var resp []*MetricOne

	if metric != nil {
		resp = make([]*MetricOne, 0, len(metric.MetricValues))
		for _, v := range metric.MetricValues {
			resp = append(resp, &MetricOne{
				MetricName: metric.MetricName,
				Series:     v.Series,
				Sample:     v.Sample,
				Error:      metric.Error,
			})
		}
	}

	return resp
}

func (metric *Metric) RangeAvg() *MetricOne {
	resp := &MetricOne{}
	if metric != nil {
		resp.MetricName = metric.MetricName
		resp.Error = metric.Error
		if len(metric.MetricValues) != 0 {
			resp.Series = metric.MetricValues[0].Series
			resp.Sample = metric.MetricValues[0].Sample
		}
	}
	return resp
}

func (metric *Metric) LastAvg() float64 {
	resp := 0.0
	if metric != nil {
		for _, v := range metric.MetricValues {
			if len(v.Series) != 0 {
				resp += v.Series[len(v.Series)-1].Value()
			}
		}
		resp = resp / float64(len(metric.MetricValues))
	}
	return resp
}

func (metric *Metric) LastSum() float64 {
	resp := 0.0
	if metric != nil {
		for _, v := range metric.MetricValues {
			if len(v.Series) != 0 {
				resp += v.Series[len(v.Series)-1].Value()
			}
		}
	}

	return resp
}

type MetricValue struct {
	Name     string            `json:"name"`
	Max      float64           `json:"max"`
	Metadata map[string]string `json:"metadata"`
	Sample   *Point            `json:"sample"`
	Series   []Point           `json:"series"`
}
type Point [2]float64

func (p Point) Timestamp() float64 {
	return p[0]
}

func (p Point) Value() float64 {
	return p[1]
}

type MetricData struct {
	MetricType   string        `json:"metric_type"`
	MetricValues []MetricValue `json:"metric_values"`
}

func (p Point) MarshalJSON() ([]byte, error) {
	t, err := jsoniter.Marshal(p.Timestamp())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	v, err := jsoniter.Marshal(strconv.FormatFloat(p.Value(), 'f', -1, 64))
	if err != nil {
		return nil, errors.WithStack(err)

	}
	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
}

func (p *Point) UnmarshalJSON(b []byte) error {
	var v []interface{}
	if err := jsoniter.Unmarshal(b, &v); err != nil {
		return errors.WithStack(err)

	}

	if v == nil {
		return nil
	}

	if len(v) != 2 {
		return errors.Errorf("unsupported array length")
	}

	ts, ok := v[0].(float64)
	if !ok {
		return errors.Errorf("failed to unmarshal [timestamp]")
	}
	valstr, ok := v[1].(string)
	if !ok {
		return errors.Errorf("failed to unmarshal [value]")
	}
	valf, err := strconv.ParseFloat(valstr, 64)
	if err != nil {
		return errors.WithStack(err)
	}

	p[0] = ts
	p[1] = valf
	return nil
}
