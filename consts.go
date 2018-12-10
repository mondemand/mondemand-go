package mondemand

type MessageType string

const (
	StatMsg = MessageType("MonDemand::StatsMsg")
	PerMsg  = MessageType("MonDemand::PerfMsg")
)
