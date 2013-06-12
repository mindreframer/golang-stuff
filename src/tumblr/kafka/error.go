// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import "errors"

// Error conditions that can occur in the Kafka client.
var (
	ErrIO           = errors.New("io or network")
	ErrClosed       = errors.New("already closed")
	ErrArg          = errors.New("invalid argument")
	ErrWire         = errors.New("invalid or dangerous data in wire format")
	ErrNotSupported = errors.New("not supported")
	ErrCompression  = errors.New("unknown compression")
	ErrChecksum     = errors.New("checksum")
	ErrNoBrokers    = errors.New("no brokers")
)

// KafkaError is a wrapper that unified errors emitted by the remote Kafka broker
type KafkaError error

// Known error conditions returned by the Kafka broker
var (
	KafkaErrUnknown          = KafkaError(errors.New("kafka: unknown error"))
	KafkaErrNoError          = KafkaError(nil)
	KafkaErrOffsetOutOfRange = KafkaError(errors.New("kafka: offset out of range"))
	KafkaErrInvalidMessage   = KafkaError(errors.New("kafka: invalid message"))
	KafkaErrWrongPartition   = KafkaError(errors.New("kafka: wrong partition"))
	KafkaErrInvalidFetchSize = KafkaError(errors.New("kafka: invalid fetch size"))
)

// KafkaErrorCode returns the integral protocol representation of the Kafka error corresponding to err.
func KafkaErrorCode(err KafkaError) ErrorCode {
	switch err {
	case KafkaErrUnknown:
		return ErrorCodeUnknown
	case KafkaErrNoError:
		return ErrorCodeNoError
	case KafkaErrOffsetOutOfRange:
		return ErrorCodeOffsetOutOfRange
	case KafkaErrInvalidMessage:
		return ErrorCodeInvalidMessage
	case KafkaErrWrongPartition:
		return ErrorCodeWrongPartition
	case KafkaErrInvalidFetchSize:
		return ErrorCodeInvalidFetchSize
	}
	panic("unknown kafka error")
}

// KafkaCodeError returns the error object corresponding to the integral Kafka protocol error code.
func KafkaCodeError(code ErrorCode) KafkaError {
	switch code {
	case ErrorCodeUnknown:
		return KafkaErrUnknown
	case ErrorCodeNoError:
		return KafkaErrNoError
	case ErrorCodeOffsetOutOfRange:
		return KafkaErrOffsetOutOfRange
	case ErrorCodeInvalidMessage:
		return KafkaErrInvalidMessage
	case ErrorCodeWrongPartition:
		return KafkaErrWrongPartition
	case ErrorCodeInvalidFetchSize:
		return KafkaErrInvalidFetchSize
	}
	panic("unknown kafka error code")
}

// ErrorCode represents a Kafka response error code
type ErrorCode int16

// List of known Kafka broker error codes
const (
	ErrorCodeUnknown ErrorCode = iota - 1
	ErrorCodeNoError
	ErrorCodeOffsetOutOfRange
	ErrorCodeInvalidMessage
	ErrorCodeWrongPartition
	ErrorCodeInvalidFetchSize
)

func isValidErrorCode(e ErrorCode) bool {
	return e >= ErrorCodeUnknown && e <= ErrorCodeInvalidFetchSize
}

// String returns a textual representation of the error code
func (x ErrorCode) String() string {
	switch x {
	case ErrorCodeUnknown:
		return "unknown"
	case ErrorCodeNoError:
		return "ok"
	case ErrorCodeOffsetOutOfRange:
		return "offset out of range"
	case ErrorCodeInvalidMessage:
		return "invalid message"
	case ErrorCodeWrongPartition:
		return "wrong partition"
	case ErrorCodeInvalidFetchSize:
		return "invalid fetch size"
	}
	return "Error code not implemented"
}
