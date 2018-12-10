package mondemand

import (
	"fmt"
	"github.com/lwes/lwes-go"
)

type Context map[string]string

func DecodeContext(event *lwes.LwesEvent) Context {
	var num_ctxt int
	var ctx Context

	if val, ok := event.Attrs["ctxt_num"].(uint16); ok {
		num_ctxt = int(val)
		ctx = make(map[string]string, num_ctxt)
	} else {
		return ctx
	}

	for i := 0; i < num_ctxt; i++ {
		k := event.Attrs[fmt.Sprint("ctxt_k", i)].(string)
		v := event.Attrs[fmt.Sprint("ctxt_v", i)].(string)
		ctx[k] = v
	}
	return ctx
}
