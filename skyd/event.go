package skyd

import (
	"bytes"
	"fmt"
	"github.com/ugorji/go-msgpack"
	"io"
	"time"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// An Event is a state change that occurs at a particular point in time.
type Event struct {
	Timestamp time.Time
	Data      map[int64]interface{}
}

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// NewEvent returns a new Event.
func NewEvent(timestamp string, data map[int64]interface{}) *Event {
	if data == nil {
		data = make(map[int64]interface{})
	}

	t, _ := time.Parse(time.RFC3339, timestamp)
	return &Event{
		Timestamp: t,
		Data:      data,
	}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Encoding
//--------------------------------------

// Encodes an event to MsgPack format.
func (e *Event) EncodeRaw(writer io.Writer) error {
	raw := []interface{}{ShiftTime(e.Timestamp), e.Data}
	encoder := msgpack.NewEncoder(writer)
	err := encoder.Encode(raw)
	return err
}

// Encodes an event to MsgPack format and returns the byte array.
func (e *Event) MarshalRaw() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := e.EncodeRaw(buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Decodes an event from MsgPack format.
func (e *Event) DecodeRaw(reader io.Reader) error {
	raw := make([]interface{}, 2)
	err := msgpack.NewDecoder(reader, nil).Decode(&raw)
	if err != nil {
		return err
	}

	// Convert the timestamp to int64.
	if timestamp, ok := normalize(raw[0]).(int64); ok {
		e.Timestamp = UnshiftTime(timestamp).UTC()
	} else {
		return fmt.Errorf("Unable to parse timestamp: '%v'", raw[0])
	}

	// Convert data to appropriate map.
	if raw[1] != nil {
		e.Data, err = e.decodeRawMap(raw[1].(map[interface{}]interface{}))
		if err != nil {
			return err
		}
	}

	return nil
}

// Decodes the map.
func (e *Event) decodeRawMap(raw map[interface{}]interface{}) (map[int64]interface{}, error) {
	m := make(map[int64]interface{})
	for k, v := range raw {
		if ki, ok := normalize(k).(int64); ok {
			m[ki] = normalize(v)
		} else {
			return nil, fmt.Errorf("Invalid property key: %v", k)
		}
	}
	return m, nil
}

func (e *Event) UnmarshalRaw(data []byte) error {
	return e.DecodeRaw(bytes.NewBuffer(data))
}

//--------------------------------------
// Comparator
//--------------------------------------

// Compares two events for equality.
func (e *Event) Equal(x *Event) bool {
	if !e.Timestamp.Equal(x.Timestamp) {
		return false
	}
	for k, v := range e.Data {
		if normalize(v) != normalize(x.Data[k]) {
			return false
		}
	}
	for k, v := range x.Data {
		if normalize(v) != normalize(e.Data[k]) {
			return false
		}
	}
	return true
}

//--------------------------------------
// Merging / Deduplication
//--------------------------------------

// Merges the data of another event into this event.
func (e *Event) Merge(a *Event) {
	if e.Data == nil && a.Data != nil {
		e.Data = make(map[int64]interface{})
	}
	for k, v := range a.Data {
		e.Data[k] = v
	}
}

// Merges the persistent data of another event into this event.
func (e *Event) MergePermanent(a *Event) {
	for k, v := range a.Data {
		if k > 0 {
			e.Data[k] = v
		}
	}
}

// Removes data in the event that is present in another event.
func (e *Event) Dedupe(a *Event) {
	for k, v := range a.Data {
		if normalize(e.Data[k]) == normalize(v) {
			delete(e.Data, k)
		}
	}
}
