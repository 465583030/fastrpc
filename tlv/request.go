package tlv

import (
	"bufio"
	"fmt"
	"sync"
)

// Request is a TLV request.
type Request struct {
	// Value is request value.
	Value []byte

	name    []byte
	sizeBuf [4]byte
}

// Reset resets the given request.
func (req *Request) Reset() {
	req.name = req.name[:0]
	req.Value = req.Value[:0]
}

// SetName sets request name.
func (req *Request) SetName(name string) {
	req.name = append(req.name[:0], name...)
}

// SetNameBytes set request name.
func (req *Request) SetNameBytes(name []byte) {
	req.name = append(req.name[:0], name...)
}

// Name returns request name.
func (req *Request) Name() []byte {
	return req.name
}

// Write appends p to the request value.
//
// It implements io.Writer.
func (req *Request) Write(p []byte) (int, error) {
	req.Value = append(req.Value, p...)
	return len(p), nil
}

// SwapValue swaps the given value with the request's value.
//
// It is forbidden accessing the swapped value after the call.
func (req *Request) SwapValue(value []byte) []byte {
	v := req.Value
	req.Value = value
	return v
}

// WriteRequest writes the request to bw.
//
// It implements fastrpc.RequestWriter
func (req *Request) WriteRequest(bw *bufio.Writer) error {
	if err := writeBytes(bw, req.name, req.sizeBuf[:]); err != nil {
		return fmt.Errorf("cannot write request name: %s", err)
	}
	if err := writeBytes(bw, req.Value, req.sizeBuf[:]); err != nil {
		return fmt.Errorf("cannot write request value: %s", err)
	}
	return nil
}

// ReadRequest reads the request from br.
func (req *Request) ReadRequest(br *bufio.Reader) error {
	var err error
	req.name, err = readBytes(br, req.name[:0], req.sizeBuf[:])
	if err != nil {
		return fmt.Errorf("cannot read request name: %s", err)
	}
	req.Value, err = readBytes(br, req.Value[:0], req.sizeBuf[:])
	if err != nil {
		return fmt.Errorf("cannot read request value: %s", err)
	}
	return nil
}

// AcquireRequest acquires new request.
func AcquireRequest() *Request {
	v := requestPool.Get()
	if v == nil {
		v = &Request{}
	}
	return v.(*Request)
}

// ReleaseRequest releases the given request.
func ReleaseRequest(req *Request) {
	req.Reset()
	requestPool.Put(req)
}

var requestPool sync.Pool
