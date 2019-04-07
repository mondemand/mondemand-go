package mondemand

type MessageType string

const (
	StatMsg = MessageType("MonDemand::StatsMsg")
	StatsMsgType = MessageType("MonDemand::StatsMsg")
	PerfMsgType  = MessageType("MonDemand::PerfMsg")
	TraceMsgType = MessageType("MonDemand::TraceMsg")
)
