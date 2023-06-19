package main

import (
	"fmt"
	"strconv"
)

// TraceContext holds trace context for an incoming or outgoing request.
type TraceContext struct {
	Trace    [16]byte
	Span     [8]byte
	ParentID [8]byte
	Options  uint8
}

func FormatJaegerTraceparentHeader(tc TraceContext) string {

	// uber-trace-id:118c6c15301b9b3b3:56e66177e6e55a91:18c6c15301b9b3b3:1
	// elastic-apm-traceparent:00-f109f092a7d869fb4615784bacefcfd7-5bf936f4fcde3af0-01

	// Jaeger span value format {trace-id}:{span-id}:{parent-span-id}:{flags}
	traceId := fmt.Sprintf("%032x", tc.Trace[:])

	var i int
	for i = range traceId {
		if traceId[i] != '0' {
			break
		}
	}

	fmt.Printf("trace id: (0x)%v, (byte)%v, 0 index=%v\n", traceId, tc.Trace, i)
	if i < 16 {
		return fmt.Sprintf("%032x:%016x:%016x:%x", tc.Trace[:], tc.Span[:], tc.ParentID[:], tc.Options)
	} else {
		return fmt.Sprintf("%016x:%016x:%016x:%x", tc.Trace[8:], tc.Span[:], tc.ParentID[:], tc.Options)
	}
}

func main() {
	// tid := "f109f092a7d869fb4615784bacefcfd7"
	tid := "000000000000000018c6c15301b9b3b3"
	var tidbyte [16]byte
	j := 0
	for i := 0; i < len(tid); i += 2 {
		b, _ := strconv.ParseUint(string(tid[i:i+2]), 16, 8)
		tidbyte[j] = byte(b)
		j += 1
	}

	span := "56e66177e6e55a91"
	var spanbyte [8]byte
	j = 0
	for i := 0; i < len(span); i += 2 {
		b, _ := strconv.ParseUint(string(span[i:i+2]), 16, 8)
		spanbyte[j] = byte(b)
		j += 1
	}

	parentid := "18c6c15301b9b3b3"
	var parentidbyte [8]byte
	j = 0
	for i := 0; i < len(parentid); i += 2 {
		b, _ := strconv.ParseUint(string(parentid[i:i+2]), 16, 8)
		parentidbyte[j] = byte(b)
		j += 1
	}

	tc := TraceContext{Trace: tidbyte,
		Span:     spanbyte,
		ParentID: parentidbyte,
		Options:  0}

	fmt.Printf("uber-trace-id: %v\n", FormatJaegerTraceparentHeader(tc))
}
