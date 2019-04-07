package mondemand

import "github.com/lwes/lwes-go"

const (
	ProgIdLabel   = "mondemand.prog_id"
	TraceIdLabel  = "mondemand.trace_id"
	OwnerLabel    = "mondemand.owner"
	HostnameLabel = "mondemand.src_host"
	MessageLabel  = "mondemand.message"
)

type TraceMsg struct {
	ProgId  string
	TraceId string
	Owner string
	Hostname string
	Message string
	Extra map[string]string
}

func DecodeTraceMsg(event *lwes.LwesEvent) TraceMsg {
	msg := TraceMsg{
		Extra: make(map[string]string),
	}

	for k, v := range event.Attrs {
		switch k {
		case ProgIdLabel:
			msg.ProgId = v.(string)
		case TraceIdLabel:
			msg.TraceId = v.(string)
		case OwnerLabel:
			msg.Owner = v.(string)
		case HostnameLabel:
			msg.Hostname = v.(string)
		case MessageLabel:
			msg.Message = v.(string)
		default:
			msg.Extra[k] = v.(string)
		}
	}
	return msg
}
