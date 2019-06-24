package mondemand

import (
	"github.com/lwes/lwes-go"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestDecodeTraceMsg(t *testing.T) {
	event := lwes.LwesEvent{
		Attrs: map[string]interface{}{
			ProgIdLabel: "programid",
			OwnerLabel: "trowner",
			TraceIdLabel: "trid",
			HostnameLabel: "host1",
			MessageLabel: "trace_message",
			"key1": "val2",
		},
	}

	trace := DecodeTraceMsg(&event)
	assert.Equal(t, "programid", trace.ProgId)
	assert.Equal(t, "trowner", trace.Owner)
	assert.Equal(t, "trid", trace.TraceId)
	assert.Equal(t, "host1", trace.Hostname)
	assert.Equal(t, "trace_message", trace.Message)
	assert.Equal(t, map[string]interface{}{
		"key1": "val2",
	}, trace.Extra)
}
