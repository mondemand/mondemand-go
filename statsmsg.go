package mondemand

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/lwes/lwes-go"
)

type MetricType string

const (
	GaugeType   = MetricType("gauge")
	CounterType = MetricType("counter")
	UnknownType = MetricType("unknown")
	statsetMetricsNum = 11
	statsetSeparator = `:`
)

var statsetMetricNames = []string{
	"count",
	"sum",
	"min",
	"max",
	"avg",
	"p50",
	"p75",
	"p90",
	"p95",
	"p98",
	"p99",
}

func getK(key string, idx int) string {
	return fmt.Sprint(key, idx)
}

type Metric struct {
	Typ   MetricType
	Key   string
	Value int64
}

type StatsMsg struct {
	ProgId  string
	Metrics []*Metric
	Context Context
}

func NewStatsMsg(progId string) *StatsMsg {
	return &StatsMsg{ProgId: progId,
		Metrics: make([]*Metric, 0, 10),
		Context: make(map[string]string),
	}
}

func (st *StatsMsg) AddContext(key, value string) {
	st.Context[key] = value
}

func (st *StatsMsg) AddMetric(tag MetricType, key string, value int64) {
	st.Metrics = append(st.Metrics, &Metric{
		Typ: tag, Key: key, Value: value})
}

func (st *StatsMsg) ToLwes() *lwes.LwesEvent {
	event := lwes.NewLwesEvent("MonDemand::StatsMsg")
	event.Set("prog_id", st.ProgId)

	event.Set("ctxt_num", uint16(len(st.Context)))
	idx := 0
	for key, value := range st.Context {
		event.Set(getK("ctxt_k", idx), key)
		event.Set(getK("ctxt_v", idx), value)
		idx = +1
	}

	event.Set("num", uint16(len(st.Metrics)))
	for idx, metric := range st.Metrics {
		event.Set(getK("t", idx), string(metric.Typ))
		event.Set(getK("k", idx), metric.Key)
		event.Set(getK("v", idx), metric.Value)
	}

	return event
}

func getMetricType(typ string) MetricType {
	switch typ {
	case "gauge":
		return GaugeType
	case "counter":
		return CounterType
	}
	return UnknownType
}

type statset [11]int64

func decodeIntOrZero(value string) int64 {
	if value == "" {
		return 0
	}
	parsed, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func newStatset(value string) *statset {
	s := statset{}
	values := strings.SplitN(value, statsetSeparator, statsetMetricsNum)
	for idx := range values {
		s[idx] = decodeIntOrZero(values[idx])
	}
	return &s
}

func decodeStatsSetToGauges(value string, key string) []*Metric {
	values := newStatset(value)
	metrics := make([]*Metric, 0)
	for idx, name := range statsetMetricNames {
		metrics= append(metrics, &Metric{
			Key:   fmt.Sprintf("%s_%s", key, name),
			Typ:   GaugeType,
			Value: values[idx],
		})
	}
	return metrics
}

func decodeMetrics(event *lwes.LwesEvent) []*Metric {
	numMetrics := int(event.Attrs["num"].(uint16))
	metrics := make([]*Metric, 0, numMetrics)
	for i := 0; i < numMetrics; i++ {
		key := event.Attrs[getK("k", i)].(string)
		typ := event.Attrs[getK("t", i)].(string)
		if typ == "statset" {
			value := event.Attrs[getK("v", i)].(string)
			metrics = append(metrics, decodeStatsSetToGauges(value, key)...)
		} else {

			value := event.Attrs[getK("v", i)].(int64)
			metrics = append(metrics, &Metric{
				Key:   key,
				Typ:   getMetricType(typ),
				Value: value,
			})
		}
	}
	return metrics
}

func DecodeStatsMsg(event *lwes.LwesEvent) StatsMsg {
	msg := StatsMsg{
		Context: DecodeContext(event),
		Metrics: decodeMetrics(event),
		ProgId:  event.Attrs["prog_id"].(string),
	}
	return msg
}
