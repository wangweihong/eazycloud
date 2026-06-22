package libprometheus

import (
	"context"
	"time"

	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/wangweihong/eazycloud/apis/iprometheus"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
)

func GetMetric(parent context.Context, conf *iprometheus.Config, expr string, ts time.Time) (iprometheus.Metric, error) {
	var resp iprometheus.Metric

	p, err := newPrometheusClient(conf)
	if err != nil {
		return resp, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	value, _, err := p.client.Query(ctx, expr, ts)
	if err != nil {
		return resp, errors.WithStack(err)
	}

	resp.MetricData = parseQueryResp(value)
	return resp, nil
}

func GetMetricOverTime(parent context.Context, conf *iprometheus.Config, expr string, start, end time.Time, step time.Duration) (iprometheus.Metric, error) {
	resp := iprometheus.Metric{}

	p, err := newPrometheusClient(conf)
	if err != nil {
		return resp, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	timeRange := apiv1.Range{Start: start, End: end, Step: step}
	value, _, err := p.client.QueryRange(ctx, expr, timeRange)
	if err != nil {
		return resp, errors.WithStack(err)
	}
	resp.MetricData = parseQueryRangeResp(value)
	return resp, nil
}

func AlertList(parent context.Context, conf *iprometheus.Config, timeout ...int) (*apiv1.AlertsResult, error) {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()

	alertRet, err := p.client.Alerts(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &alertRet, nil
}

func AlertManagers(parent context.Context, conf *iprometheus.Config, timeout ...int) (*apiv1.AlertManagersResult, error) {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()

	alertRet, err := p.client.AlertManagers(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &alertRet, nil
}

func Config(parent context.Context, conf *iprometheus.Config, timeout ...int) (*apiv1.ConfigResult, error) {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()

	alertRet, err := p.client.Config(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &alertRet, nil
}

func Rules(parent context.Context, conf *iprometheus.Config, timeout ...int) (*apiv1.RulesResult, error) {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()

	ret, err := p.client.Rules(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &ret, nil
}

func Targets(parent context.Context, conf *iprometheus.Config, timeout ...int) (*apiv1.TargetsResult, error) {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()

	ret, err := p.client.Targets(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &ret, nil
}

func Series(parent context.Context, conf *iprometheus.Config, matches []string, startTime time.Time, endTime time.Time, timeout ...int) ([]model.LabelSet, error) {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()

	ret, _, err := p.client.Series(ctx, matches, startTime, endTime)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ret, nil
}

func LabelValues(parent context.Context, conf *iprometheus.Config, label string, timeout ...int) (model.LabelValues, error) {
	// p, err := newPrometheusClient(conf)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }
	// defer p.Close()

	// ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	// defer cancel()
	// ret,_, err := p.client.LabelValues(ctx, label)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// return ret, nil
	return nil, errors.Errorf("not support yet")
}

func DeleteSeries(parent context.Context, conf *iprometheus.Config, matches []string, t time.Time) error {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, defaultTimeout*time.Second)
	defer cancel()
	end := time.Unix(time.Now().Unix(), 0)
	start := time.Unix(time.Now().Unix()-100*24*60*60, 0)
	err = p.client.DeleteSeries(ctx, matches, start, end)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Flags(parent context.Context, conf *iprometheus.Config) (*apiv1.FlagsResult, error) {
	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, defaultTimeout*time.Second)
	defer cancel()
	ret, err := p.client.Flags(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &ret, nil
}

func GetNamedMetrics(parent context.Context, conf *iprometheus.Config, metrics []string, opt *iprometheus.ProqlOption, timeout ...int) ([]iprometheus.Metric, error) {
	var resp []iprometheus.Metric

	p, err := newPrometheusClient(conf)
	if err != nil {
		return resp, errors.WithStack(err)
	}
	defer p.Close()

	wg := waitgroup.RunGenericConcurrently(parent, metrics, func(ctx context.Context, metric string) waitgroup.GenericResult[iprometheus.Metric] {
		parsedResp := iprometheus.Metric{MetricName: metric}

		ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
		defer cancel()

		value, _, err := p.client.Query(ctx, makeProql(metric, opt), opt.Time)
		if err != nil {
			parsedResp.Error = err.Error()
		} else {
			parsedResp.MetricData = parseQueryResp(value)
		}
		return waitgroup.NewGenericResult(parsedResp, nil)
	})

	for _, v := range wg.GetResults() {
		resp = append(resp, v.Data)
	}

	return resp, nil
}

func GetNamedMetricsTimeout(parent context.Context, conf *iprometheus.Config, metrics []string, opt *iprometheus.ProqlOption, timeout int64) ([]iprometheus.Metric, error) {
	var resp []iprometheus.Metric

	p, err := newPrometheusClient(conf)
	if err != nil {
		return resp, errors.WithStack(err)
	}
	defer p.Close()

	wg := waitgroup.RunGenericConcurrently(parent, metrics, func(ctx context.Context, metric string) waitgroup.GenericResult[iprometheus.Metric] {
		parsedResp := iprometheus.Metric{MetricName: metric}

		ctx, cancel := context.WithTimeout(parent, time.Duration(timeout)*time.Second)
		defer cancel()
		value, _, err := p.client.Query(ctx, makeProql(metric, opt), opt.Time)
		if err != nil {
			parsedResp.Error = err.Error()
		} else {
			parsedResp.MetricData = parseQueryResp(value)
		}
		return waitgroup.NewGenericResult(parsedResp, nil)
	})

	for _, v := range wg.GetResults() {
		resp = append(resp, v.Data)
	}
	return resp, nil
}

func GetNamedMetricsHistory(parent context.Context, conf *iprometheus.Config, metrics []string, opt *iprometheus.ProqlOption, timeout ...int) ([]iprometheus.Metric, error) {
	var resp []iprometheus.Metric

	p, err := newPrometheusClient(conf)
	if err != nil {
		return resp, errors.WithStack(err)
	}

	timeRange := apiv1.Range{Start: opt.StartTime, End: opt.EndTime, Step: opt.Duration}

	wg := waitgroup.RunGenericConcurrently(parent, metrics, func(ctx context.Context, metric string) waitgroup.GenericResult[iprometheus.Metric] {
		parsedResp := iprometheus.Metric{MetricName: metric}

		ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
		defer cancel()
		value, _, err := p.client.QueryRange(ctx, makeProql(metric, opt), timeRange)
		if err != nil {
			parsedResp.Error = err.Error()

		} else {
			parsedResp.MetricData = parseQueryRangeResp(value)
		}

		return waitgroup.NewGenericResult(parsedResp, nil)
	})

	for _, v := range wg.GetResults() {
		resp = append(resp, v.Data)
	}
	return resp, nil
}

func GetAutoMetricsHistory(parent context.Context, conf *iprometheus.Config, metrics []string, opt *iprometheus.ProqlOption, timeout ...int) ([]iprometheus.Metric, error) {
	var resp []iprometheus.Metric

	p, err := newPrometheusClient(conf)
	if err != nil {
		return resp, errors.WithStack(err)
	}

	timeRange := apiv1.Range{Start: opt.StartTime, End: opt.EndTime, Step: opt.Duration}

	wg := waitgroup.RunGenericConcurrently(parent, metrics, func(ctx context.Context, metric string) waitgroup.GenericResult[iprometheus.Metric] {
		parsedResp := iprometheus.Metric{MetricName: metric}

		ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
		defer cancel()
		value, _, err := p.client.QueryRange(ctx, metric, timeRange)
		if err != nil {
			parsedResp.Error = err.Error()
		} else {
			parsedResp.MetricData = parseQueryRangeResp(value)
		}

		return waitgroup.NewGenericResult(parsedResp, nil)
	})
	for _, v := range wg.GetResults() {
		resp = append(resp, v.Data)
	}
	return resp, nil
}

func GetAutoMetricsRealTime(parent context.Context, conf *iprometheus.Config, metrics string, timeout ...int) (*iprometheus.Metric, error) {
	resp := &iprometheus.Metric{MetricName: metrics}

	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer p.Close()

	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()
	value, _, err := p.client.Query(ctx, metrics, time.Now())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp.MetricData = parseQueryResp(value)
	return resp, nil
}

func GetAutoMetricsHistory2(parent context.Context, conf *iprometheus.Config, metric string, timeRange apiv1.Range, timeout ...int) (*iprometheus.Metric, error) {
	resp := &iprometheus.Metric{MetricName: metric}

	p, err := newPrometheusClient(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// time range: end.Sub(start)/step > 11000
	ctx, cancel := context.WithTimeout(parent, getTimeOut(timeout...)*time.Second)
	defer cancel()
	value, _, err := p.client.QueryRange(ctx, metric, timeRange)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp.MetricData = parseQueryRangeResp(value)
	return resp, nil
}

func parseQueryRangeResp(value model.Value) iprometheus.MetricData {
	resp := iprometheus.MetricData{MetricType: iprometheus.MetricTypeMatrix}
	data, _ := value.(model.Matrix)

	for _, v := range data {
		mv := iprometheus.MetricValue{
			Metadata: make(map[string]string),
		}

		for k, v := range v.Metric {
			mv.Metadata[string(k)] = string(v)
		}

		for _, k := range v.Values {
			mv.Series = append(mv.Series, iprometheus.Point{float64(k.Timestamp) / 1000, float64(k.Value)})
		}

		resp.MetricValues = append(resp.MetricValues, mv)
	}
	return resp
}

func parseQueryResp(value model.Value) iprometheus.MetricData {
	res := iprometheus.MetricData{MetricType: iprometheus.MetricTypeVector}
	data, _ := value.(model.Vector)
	for _, v := range data {
		mv := iprometheus.MetricValue{
			Metadata: make(map[string]string),
		}

		for k, v := range v.Metric {
			mv.Metadata[string(k)] = string(v)
		}

		mv.Sample = &iprometheus.Point{float64(v.Timestamp), float64(v.Value)} // timestamp  -> value
		res.MetricValues = append(res.MetricValues, mv)
	}
	return res
}
