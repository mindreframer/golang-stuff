package zmq2

/*
#include <zmq.h>
#include "zmq2.h"
#include <stdint.h>
#include <stdlib.h>
*/
import "C"

import (
	"time"
	"unsafe"
)

func (soc *Socket) setString(opt C.int, s string) error {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(cs), C.size_t(len(s))); i != 0 {
		return errget(err)
	}
	return nil
}

func (soc *Socket) setInt(opt C.int, value int) error {
	val := C.int(value)
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&val), C.size_t(unsafe.Sizeof(val))); i != 0 {
		return errget(err)
	}
	return nil
}

func (soc *Socket) setInt64(opt C.int, value int64) error {
	val := C.int64_t(value)
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&val), C.size_t(unsafe.Sizeof(val))); i != 0 {
		return errget(err)
	}
	return nil
}

func (soc *Socket) setUInt64(opt C.int, value uint64) error {
	val := C.uint64_t(value)
	if i, err := C.zmq_setsockopt(soc.soc, opt, unsafe.Pointer(&val), C.size_t(unsafe.Sizeof(val))); i != 0 {
		return errget(err)
	}
	return nil
}

// ZMQ_HWM: Set high water mark
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc3
func (soc *Socket) SetHwm(value uint64) error {
	return soc.setUInt64(C.ZMQ_HWM, value)
}

// ZMQ_SWAP: Set disk offload size
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc4
func (soc *Socket) SetSwap(value int64) error {
	return soc.setInt64(C.ZMQ_SWAP, value)
}

// ZMQ_AFFINITY: Set I/O thread affinity
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc5
func (soc *Socket) SetAffinity(value uint64) error {
	return soc.setUInt64(C.ZMQ_AFFINITY, value)
}

// ZMQ_IDENTITY: Set socket identity
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc6
func (soc *Socket) SetIdentity(value string) error {
	return soc.setString(C.ZMQ_IDENTITY, value)
}

// ZMQ_SUBSCRIBE: Establish message filter
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc7
func (soc *Socket) SetSubscribe(filter string) error {
	return soc.setString(C.ZMQ_SUBSCRIBE, filter)
}

// ZMQ_UNSUBSCRIBE: Remove message filter
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc8
func (soc *Socket) SetUnsubscribe(filter string) error {
	return soc.setString(C.ZMQ_UNSUBSCRIBE, filter)
}

// ZMQ_RCVTIMEO: Maximum time before a recv operation returns with EAGAIN
//
// Use -1 for infinite
//
// Returns ErrorNotImplemented in 0MQ version 2.1
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc9
func (soc *Socket) SetRcvtimeo(value time.Duration) error {
	major, minor, _ := Version()
	if major == 2 && minor < 2 {
		return ErrorNotImplemented
	}
	val := int(value / time.Millisecond)
	if value == -1 {
		val = -1
	}
	return soc.setInt(C.ZMQ_RCVTIMEO, val)
}

// ZMQ_SNDTIMEO: Maximum time before a send operation returns with EAGAIN
//
// Use -1 for infinite
//
// Returns ErrorNotImplemented in 0MQ version 2.1
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc10
func (soc *Socket) SetSndtimeo(value time.Duration) error {
	major, minor, _ := Version()
	if major == 2 && minor < 2 {
		return ErrorNotImplemented
	}
	val := int(value / time.Millisecond)
	if value == -1 {
		val = -1
	}
	return soc.setInt(C.ZMQ_SNDTIMEO, val)
}

// ZMQ_RATE: Set multicast data rate
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc11
func (soc *Socket) SetRate(value int64) error {
	return soc.setInt64(C.ZMQ_RATE, value)
}

// ZMQ_RECOVERY_IVL: Set multicast recovery interval
//
// Note: value is of type time.Duration
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc12
func (soc *Socket) SetRecoveryIvl(value time.Duration) error {
	val := int64(value / time.Second)
	return soc.setInt64(C.ZMQ_RECOVERY_IVL, val)
}

// ZMQ_RECOVERY_IVL_MSEC: Set multicast recovery interval in milliseconds
//
// Note: value is of type time.Duration
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc13
func (soc *Socket) SetRecoveryIvlMsec(value time.Duration) error {
	val := int64(value / time.Millisecond)
	if value < 0 {
		val = -1
	}
	return soc.setInt64(C.ZMQ_RECOVERY_IVL_MSEC, val)
}

// ZMQ_MCAST_LOOP: Control multicast loop-back
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc14
func (soc *Socket) SetMcastLoop(value bool) error {
	val := int64(0)
	if value {
		val = 1
	}
	return soc.setInt64(C.ZMQ_MCAST_LOOP, val)
}

// ZMQ_SNDBUF: Set kernel transmit buffer size
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc15
func (soc *Socket) SetSndbuf(value uint64) error {
	return soc.setUInt64(C.ZMQ_SNDBUF, value)
}

// ZMQ_RCVBUF: Set kernel receive buffer size
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc16
func (soc *Socket) SetRcvbuf(value uint64) error {
	return soc.setUInt64(C.ZMQ_RCVBUF, value)
}

// ZMQ_LINGER: Set linger period for socket shutdown
//
// Use -1 for infinite
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc17
func (soc *Socket) SetLinger(value time.Duration) error {
	val := int(value / time.Millisecond)
	if value == -1 {
		val = -1
	}
	return soc.setInt(C.ZMQ_LINGER, val)
}

// ZMQ_RECONNECT_IVL: Set reconnection interval
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc18
func (soc *Socket) SetReconnectIvl(value time.Duration) error {
	val := int(value / time.Millisecond)
	return soc.setInt(C.ZMQ_RECONNECT_IVL, val)
}

// ZMQ_RECONNECT_IVL_MAX: Set maximum reconnection interval
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc19
func (soc *Socket) SetReconnectIvlMax(value time.Duration) error {
	val := int(value / time.Millisecond)
	return soc.setInt(C.ZMQ_RECONNECT_IVL_MAX, val)
}

// ZMQ_BACKLOG: Set maximum length of the queue of outstanding connections
//
// See: http://api.zeromq.org/2-2:zmq-setsockopt#toc20
func (soc *Socket) SetBacklog(value int) error {
	return soc.setInt(C.ZMQ_BACKLOG, value)
}
