package zmq2

/*
#include <zmq.h>
#include "zmq2.h"
#include <stdint.h>
*/
import "C"

import (
	"time"
	"unsafe"
)

func (soc *Socket) getString(opt C.int, bufsize int) (string, error) {
	value := make([]byte, bufsize)
	size := C.size_t(bufsize)
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value[0]), &size); i != 0 {
		return "", errget(err)
	}
	return string(value[:int(size)]), nil
}

func (soc *Socket) getInt(opt C.int) (int, error) {
	value := C.int(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return int(value), nil
}

func (soc *Socket) getInt64(opt C.int) (int64, error) {
	value := C.int64_t(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return int64(value), nil
}

func (soc *Socket) getUInt64(opt C.int) (uint64, error) {
	value := C.uint64_t(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return uint64(value), nil
}

func (soc *Socket) getUInt32(opt C.int) (uint32, error) {
	value := C.uint32_t(0)
	size := C.size_t(unsafe.Sizeof(value))
	if i, err := C.zmq_getsockopt(soc.soc, opt, unsafe.Pointer(&value), &size); i != 0 {
		return 0, errget(err)
	}
	return uint32(value), nil
}

// ZMQ_TYPE: Retrieve socket type
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc3
func (soc *Socket) GetType() (Type, error) {
	v, err := soc.getInt(C.ZMQ_TYPE)
	return Type(v), err
}

// ZMQ_RCVMORE: More message data parts to follow
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc4
func (soc *Socket) GetRcvmore() (bool, error) {
	v, err := soc.getInt64(C.ZMQ_RCVMORE)
	return v != 0, err
}

// ZMQ_HWM: Retrieve high water mark
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc5
func (soc *Socket) GetHwm() (uint64, error) {
	return soc.getUInt64(C.ZMQ_HWM)
}

// ZMQ_RCVTIMEO: Maximum time before a socket operation returns with EAGAIN
//
// Returns time.Duration(-1) for infinite
//
// Returns ErrorNotImplemented in 0MQ version 2.1
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc6
func (soc *Socket) GetRcvtimeo() (time.Duration, error) {
	major, minor, _ := Version()
	if major == 2 && minor < 2 {
		return 0, ErrorNotImplemented
	}
	v, err := soc.getInt(C.ZMQ_RCVTIMEO)
	if v < 0 {
		return time.Duration(-1), err
	}
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_SNDTIMEO: Maximum time before a socket operation returns with EAGAIN
//
// Returns time.Duration(-1) for infinite
//
// Returns ErrorNotImplemented in 0MQ version 2.1
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc7
func (soc *Socket) GetSndtimeo() (time.Duration, error) {
	major, minor, _ := Version()
	if major == 2 && minor < 2 {
		return 0, ErrorNotImplemented
	}
	v, err := soc.getInt(C.ZMQ_SNDTIMEO)
	if v < 0 {
		return time.Duration(-1), err
	}
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_SWAP: Retrieve disk offload size
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc8
func (soc *Socket) GetSwap() (int64, error) {
	return soc.getInt64(C.ZMQ_SWAP)
}

// ZMQ_AFFINITY: Retrieve I/O thread affinity
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc9
func (soc *Socket) GetAffinity() (uint64, error) {
	return soc.getUInt64(C.ZMQ_AFFINITY)
}

// ZMQ_IDENTITY: Retrieve socket identity
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc10
func (soc *Socket) GetIdentity() (string, error) {
	return soc.getString(C.ZMQ_IDENTITY, 256)
}

// ZMQ_RATE: Retrieve multicast data rate
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc11
func (soc *Socket) GetRate() (int64, error) {
	return soc.getInt64(C.ZMQ_RATE)
}

// ZMQ_RECOVERY_IVL: Get multicast recovery interval
//
// Note: return time is time.Duration
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc12
func (soc *Socket) GetRecoveryIvl() (time.Duration, error) {
	v, e := soc.getInt64(C.ZMQ_RECOVERY_IVL)
	return time.Duration(v) * time.Second, e
}

// ZMQ_RECOVERY_IVL_MSEC: Get multicast recovery interval in milliseconds
//
// Note: return time is time.Duration
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc13
func (soc *Socket) GetRecoveryIvlMsec() (time.Duration, error) {
	v, e := soc.getInt64(C.ZMQ_RECOVERY_IVL_MSEC)
	if v == -1 {
		return -1, e
	}
	return time.Duration(v) * time.Millisecond, e
}

// ZMQ_MCAST_LOOP: Control multicast loop-back
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc14
func (soc *Socket) GetMcastLoop() (bool, error) {
	v, e := soc.getInt64(C.ZMQ_MCAST_LOOP)
	if v == 0 {
		return false, e
	}
	return true, e
}

// ZMQ_SNDBUF: Retrieve kernel transmit buffer size
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc15
func (soc *Socket) GetSndbuf(value int) (uint64, error) {
	return soc.getUInt64(C.ZMQ_SNDBUF)
}

// ZMQ_RCVBUF: Retrieve kernel receive buffer size
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc16
func (soc *Socket) GetRcvbuf(value int) (uint64, error) {
	return soc.getUInt64(C.ZMQ_RCVBUF)
}

// ZMQ_LINGER: Retrieve linger period for socket shutdown
//
// Returns time.Duration(-1) for infinite
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc17
func (soc *Socket) GetLinger() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_LINGER)
	if v < 0 {
		return time.Duration(-1), err
	}
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_RECONNECT_IVL: Retrieve reconnection interval
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc18
func (soc *Socket) GetReconnectIvl() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_RECONNECT_IVL)
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_RECONNECT_IVL_MAX: Retrieve maximum reconnection interval
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc19
func (soc *Socket) GetReconnectIvlMax() (time.Duration, error) {
	v, err := soc.getInt(C.ZMQ_RECONNECT_IVL_MAX)
	return time.Duration(v) * time.Millisecond, err
}

// ZMQ_BACKLOG: Retrieve maximum length of the queue of outstanding connections
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc20
func (soc *Socket) GetBacklog() (int, error) {
	return soc.getInt(C.ZMQ_BACKLOG)
}

// ZMQ_FD: Retrieve file descriptor associated with the socket
// see socketget_unix.go and socketget_windows.go

// ZMQ_EVENTS: Retrieve socket event state
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc22
func (soc *Socket) GetEvents() (State, error) {
	v, err := soc.getUInt32(C.ZMQ_EVENTS)
	return State(v), err
}
