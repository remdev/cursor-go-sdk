package connect

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// StreamReader reads Connect JSON stream envelopes from an HTTP response.
type StreamReader struct {
	resp   *http.Response
	parser *streamParser
	cancel func()
}

// Next returns the next message envelope or io.EOF when the stream ends cleanly.
func (r *StreamReader) Next() (map[string]any, error) {
	for {
		msg, done, err := r.parser.tryDequeue()
		if err != nil {
			return nil, err
		}
		if msg != nil {
			return msg, nil
		}
		if done {
			return nil, io.EOF
		}
		chunk := make([]byte, 8192)
		n, readErr := r.resp.Body.Read(chunk)
		if n > 0 {
			if err := r.parser.feed(chunk[:n]); err != nil {
				return nil, err
			}
			continue
		}
		if readErr == io.EOF {
			if err := r.parser.assertEnded(); err != nil {
				return nil, err
			}
			return nil, io.EOF
		}
		if readErr != nil {
			return nil, readErr
		}
	}
}

// Close releases the underlying HTTP response.
func (r *StreamReader) Close() error {
	if r.cancel != nil {
		r.cancel()
	}
	if r.resp != nil && r.resp.Body != nil {
		return r.resp.Body.Close()
	}
	return nil
}

type streamParser struct {
	buffer  []byte
	pending []map[string]any
	status  int
	ended   bool
}

func newStreamParser(status int) *streamParser {
	return &streamParser{status: status, pending: make([]map[string]any, 0)}
}

func (p *streamParser) feed(chunk []byte) error {
	p.buffer = append(p.buffer, chunk...)
	for !p.ended && len(p.buffer) >= 5 {
		flags := p.buffer[0]
		length := binary.BigEndian.Uint32(p.buffer[1:5])
		if len(p.buffer) < int(5+length) {
			return nil
		}
		payload := p.buffer[5 : 5+length]
		p.buffer = p.buffer[5+length:]
		var envelope map[string]any
		if len(payload) > 0 {
			if err := json.Unmarshal(payload, &envelope); err != nil {
				return fmt.Errorf("bridge returned invalid JSON in stream: %w", err)
			}
		}
		if flags&EndStreamFlag != 0 {
			p.ended = true
			if errObj, ok := envelope["error"].(map[string]any); ok && errObj != nil {
				return connectError(errObj, p.status, nil)
			}
			return nil
		}
		p.pending = append(p.pending, envelope)
	}
	return nil
}

func (p *streamParser) tryDequeue() (map[string]any, bool, error) {
	if len(p.pending) > 0 {
		msg := p.pending[0]
		p.pending = p.pending[1:]
		return msg, false, nil
	}
	if p.ended {
		return nil, true, nil
	}
	return nil, false, nil
}

func (p *streamParser) assertEnded() error {
	if !p.ended {
		return fmt.Errorf("unexpected end of Connect stream")
	}
	return nil
}

// EncodeStreamEnvelope builds a Connect JSON stream frame.
func EncodeStreamEnvelope(payload map[string]any, end bool) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	flags := byte(0)
	if end {
		flags = EndStreamFlag
	}
	buf := make([]byte, 5+len(body))
	buf[0] = flags
	binary.BigEndian.PutUint32(buf[1:5], uint32(len(body)))
	copy(buf[5:], body)
	return buf, nil
}
