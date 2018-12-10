package mondemand

import (
	"fmt"
	"github.com/lwes/lwes-go"
	"net"
	"strings"
	"time"
)

type Timeline struct {
	Label      string
	Start, End time.Time
}

type PerfMsg struct {
	Context     Context
	Timelines   []*Timeline
	CallerLabel string
	PerfId      string
	ReceiptTime int64
	SenderIP    net.IP
	SenderPort  uint16
}

func decodeTimeline(event *lwes.LwesEvent) []*Timeline {
	numTimelines := int(event.Attrs["num"].(uint16))
	timelines := make([]*Timeline, 0, numTimelines)
	for i := 0; i < numTimelines; i++ {
		label := event.Attrs[fmt.Sprint("label", i)].(string)
		start := event.Attrs[fmt.Sprint("start", i)].(int64)
		end := event.Attrs[fmt.Sprint("end", i)].(int64)
		timelines = append(timelines, &Timeline{label, toTime(start), toTime(end)})
	}
	return timelines
}

func DecodePerfMsg(lwe *lwes.LwesEvent) PerfMsg {
	msg := PerfMsg{
		Context:   DecodeContext(lwe),
		Timelines: decodeTimeline(lwe),
	}

	msg.CallerLabel = lwe.Attrs["caller_label"].(string)
	msg.PerfId = lwe.Attrs["id"].(string)
	msg.ReceiptTime = lwe.Attrs["ReceiptTime"].(int64)
	msg.SenderIP = lwe.Attrs["SenderIP"].(net.IP)
	msg.SenderPort = lwe.Attrs["SenderPort"].(uint16)

	return msg
}

func toTime(msec int64) time.Time {
	return time.Unix(msec/1000, (msec%1000)*1000000)
}

func (msg *PerfMsg) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PerfMsg[%s] (at %s, from %s:%d){   }",
		msg.PerfId,
		toTime(msg.ReceiptTime).Format("2006-01-02T15:04:05.000Z07:00"),
		msg.SenderIP, msg.SenderPort,
	))
	sb.WriteString(fmt.Sprintln("{\t%v\n", msg.CallerLabel))
	for _, tl := range msg.Timelines {
		start, end := tl.Start, tl.End
		sb.WriteString(fmt.Sprintf("\t%s\t%s\n\t |%s %s|\n", tl.Label, end.Sub(start),
			start.Format("2006-01-02T15:04:05.000Z07:00"),
			end.Format("2006-01-02T15:04:05.000Z07:00"),
		))
	}
	sb.WriteString(fmt.Sprintln("}"))
	return sb.String()
}
