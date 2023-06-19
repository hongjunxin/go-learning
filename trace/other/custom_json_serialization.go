package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
)

// json 序列化自定义转换器

type traceID struct {
	Low  uint64 `json:"lo"`
	High uint64 `json:"hi"`
}

type spanID uint64

type span struct {
	TraceID traceID `protobuf:"bytes,1,opt,name=trace_id,json=traceId,proto3,customtype=TraceID" json:"trace_id"`
	SpanID  spanID  `protobuf:"bytes,2,opt,name=span_id,json=spanId,proto3,customtype=SpanID" json:"span_id"`
}

func main() {
	s := &span{TraceID: traceID{Low: 4637188659961372800, High: 0}, SpanID: 4637188659961372800}
	jsonByte, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json.Marshal() failed")
		os.Exit(1)
	}
	fmt.Printf("json.Marshal(): %v\n", string(jsonByte))

	sback := &span{}
	json.Unmarshal(jsonByte, sback)
	fmt.Printf("json.Unmarshal(): %v\n", sback)
}

// MarshalJSON converts span id into a base64 string enclosed in quotes.
// Used by protobuf JSON serialization.
// Example: {1} => "AAAAAAAAAAE=".
func (s spanID) MarshalJSON() ([]byte, error) {
	var b [8]byte
	s.MarshalTo(b[:]) // can only error on incorrect buffer size
	v := make([]byte, 12+2)
	base64.StdEncoding.Encode(v[1:13], b[:])
	v[0], v[13] = '"', '"'
	return v, nil
}

// MarshalTo converts span ID into a binary representation. Called by protobuf serialization.
func (s *spanID) MarshalTo(data []byte) (n int, err error) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(*s))
	return marshalBytes(data, b[:])
}

func marshalBytes(dst []byte, src []byte) (n int, err error) {
	if len(dst) < len(src) {
		return 0, fmt.Errorf("buffer is too short")
	}
	return copy(dst, src), nil
}

// UnmarshalJSON inflates span id from base64 string, possibly enclosed in quotes.
// User by protobuf JSON serialization.
//
// There appears to be a bug in gogoproto, as this function is only called for numeric values.
// https://github.com/gogo/protobuf/issues/411#issuecomment-393856837
func (s *spanID) UnmarshalJSON(data []byte) error {
	str := string(data)
	if l := len(str); l > 2 && str[0] == '"' && str[l-1] == '"' {
		str = str[1 : l-1]
	}
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return fmt.Errorf("cannot unmarshal SpanID from string '%s': %v", string(data), err)
	}
	return s.Unmarshal(b)
}

// Unmarshal inflates span ID from a binary representation. Called by protobuf serialization.
func (s *spanID) Unmarshal(data []byte) error {
	var err error
	*s, err = SpanIDFromBytes(data)
	return err
}

// SpanIDFromBytes creates a SpandID from list of bytes
func SpanIDFromBytes(data []byte) (spanID, error) {
	if len(data) != 8 {
		return spanID(0), fmt.Errorf("invalid length for SpanID")
	}
	return NewSpanID(binary.BigEndian.Uint64(data)), nil
}

// NewSpanID creates a new SpanID from a 64bit unsigned int.
func NewSpanID(v uint64) spanID {
	return spanID(v)
}
