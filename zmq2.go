// A Go interface to ZeroMQ version 2.
//
// For ZeroMQ version 3, see: http://github.com/pebbe/zmq3
//
// Requires ZeroMQ version 2.1 or 2.2
//
// The following functions return ErrorNotImplemented in 0MQ version 2.1:
//
// (*Socket)GetRcvtimeo, (*Socket)GetSndtimeo, (*Socket)SetRcvtimeo, (*Socket)SetSndtimeo
//
// http://www.zeromq.org/
package zmq2

/*
#cgo !windows pkg-config: libzmq
#cgo windows CFLAGS: -I/usr/local/include
#cgo windows LDFLAGS: -L/usr/local/lib -lzmq
#include <zmq.h>
#include "zmq2.h"
#include <stdlib.h>
#include <string.h>
void my_free (void *data, void *hint) {
    free (data);
}
int my_msg_init_data (zmq_msg_t *msg, void *data, size_t size) {
    return zmq_msg_init_data (msg, data, size, my_free, NULL);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var (
	ErrorNotImplemented = errors.New("Not implemented, requires 0MQ version 2.2")
)

var (
	ctx           unsafe.Pointer
	old           []unsafe.Pointer
	nr_of_threads int
)

func init() {
	var err error
	nr_of_threads = 1
	ctx, err = C.zmq_init(C.int(nr_of_threads))
	if ctx == nil {
		panic("Init of ZeroMQ context failed: " + errget(err).Error())
	}
}

//. Util

func errget(err error) error {
	errno, ok := err.(syscall.Errno)
	if ok && errno >= C.ZMQ_HAUSNUMERO {
		return errors.New(C.GoString(C.zmq_strerror(C.int(errno))))
	}
	return err
}

// Report 0MQ library version.
func Version() (major, minor, patch int) {
	var maj, min, pat C.int
	C.zmq_version(&maj, &min, &pat)
	return int(maj), int(min), int(pat)
}

// Get 0MQ error message string.
func Error(e int) string {
	return C.GoString(C.zmq_strerror(C.int(e)))
}

//. Context

// Returns the size of the 0MQ thread pool.
func GetIoThreads() (int, error) {
	return nr_of_threads, nil
}

/*
This function specifies the size of the ØMQ thread pool to handle I/O operations.
If your application is using only the inproc transport for messaging you may set
this to zero, otherwise set it to at least one.

This function creates a new context without closing the old one. Use it before
creating any sockets.

Default value   1
*/
func SetIoThreads(n int) error {
	if n != nr_of_threads {
		c, err := C.zmq_init(C.int(n))
		if c == nil {
			return errget(err)
		}
		old = append(old, ctx) // keep a reference, to prevent garbage collection
		ctx = c
		nr_of_threads = n
	}
	return nil
}

//. Sockets

// Specifies the type of a socket, used by NewSocket()
type Type int

const (
	// Constants for NewSocket()
	// See: http://api.zeromq.org/2-2:zmq-socket#toc3
	REQ    = Type(C.ZMQ_REQ)
	REP    = Type(C.ZMQ_REP)
	DEALER = Type(C.ZMQ_DEALER)
	ROUTER = Type(C.ZMQ_ROUTER)
	PUB    = Type(C.ZMQ_PUB)
	SUB    = Type(C.ZMQ_SUB)
	XPUB   = Type(C.ZMQ_XPUB)
	XSUB   = Type(C.ZMQ_XSUB)
	PUSH   = Type(C.ZMQ_PUSH)
	PULL   = Type(C.ZMQ_PULL)
	PAIR   = Type(C.ZMQ_PAIR)
)

/*
Socket type as string.
*/
func (t Type) String() string {
	switch t {
	case REQ:
		return "REQ"
	case REP:
		return "REP"
	case DEALER:
		return "DEALER"
	case ROUTER:
		return "ROUTER"
	case PUB:
		return "PUB"
	case SUB:
		return "SUB"
	case XPUB:
		return "XPUB"
	case XSUB:
		return "XSUB"
	case PUSH:
		return "PUSH"
	case PULL:
		return "PULL"
	case PAIR:
		return "PAIR"
	}
	return "<INVALID>"
}

// Used by  (*Socket)Send() and (*Socket)Recv()
type Flag int

const (
	// Flags for (*Socket)Send(), (*Socket)Recv()
	// For Send, see: http://api.zeromq.org/2-2:zmq-send#toc2
	// For Recv, see: http://api.zeromq.org/2-2:zmq-recv#toc2
	NOBLOCK = Flag(C.ZMQ_NOBLOCK)
	SNDMORE = Flag(C.ZMQ_SNDMORE)
)

/*
Socket flag as string.
*/
func (f Flag) String() string {
	ff := make([]string, 0)
	if f&NOBLOCK != 0 {
		ff = append(ff, "NOBLOCK")
	}
	if f&SNDMORE != 0 {
		ff = append(ff, "SNDMORE")
	}
	if len(ff) == 0 {
		return "<NONE>"
	}
	return strings.Join(ff, "|")
}

// Used by (soc *Socket)GetEvents()
type State uint32

const (
	// Flags for (*Socket)GetEvents()
	// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc22
	POLLIN  = State(C.ZMQ_POLLIN)
	POLLOUT = State(C.ZMQ_POLLOUT)
)

// Uses by Device()
type Dev int

const (
	// Constants for Device()
	// See: http://api.zeromq.org/2-2:zmq-device#toc2
	QUEUE     = Dev(C.ZMQ_QUEUE)
	FORWARDER = Dev(C.ZMQ_FORWARDER)
	STREAMER  = Dev(C.ZMQ_STREAMER)
)

/*
Dev as string
*/
func (d Dev) String() string {
	switch d {
	case QUEUE:
		return "QUEUE"
	case FORWARDER:
		return "FORWARDER"
	case STREAMER:
		return "STREAMER"
	}
	return "<INVALID>"
}

/*
Socket state as string.
*/
func (s State) String() string {
	ss := make([]string, 0)
	if s&POLLIN != 0 {
		ss = append(ss, "POLLIN")
	}
	if s&POLLOUT != 0 {
		ss = append(ss, "POLLOUT")
	}
	if len(ss) == 0 {
		return "<NONE>"
	}
	return strings.Join(ss, "|")
}

/*
Socket functions starting with `Set` or `Get` are used for setting and
getting socket options.
*/
type Socket struct {
	soc unsafe.Pointer
}

/*
Socket as string.
*/
func (soc Socket) String() string {
	t, _ := soc.GetType()
	i, err := soc.GetIdentity()
	if err == nil && i != "" {
		return fmt.Sprintf("Socket(%v,%q)", t, i)
	}
	return fmt.Sprintf("Socket(%v,%p)", t, soc.soc)
}

/*
Create 0MQ socket.

WARNING:
The Socket is not thread safe. This means that you cannot access the same Socket
from different goroutines without using something like a mutex.

For a description of socket types, see: http://api.zeromq.org/2-2:zmq-socket#toc3
*/
func NewSocket(t Type) (soc *Socket, err error) {
	soc = &Socket{}
	s, e := C.zmq_socket(ctx, C.int(t))
	if s == nil {
		err = errget(e)
	} else {
		soc.soc = s
		runtime.SetFinalizer(soc, (*Socket).Close)
	}
	return
}

// If not called explicitly, the socket will be closed on garbage collection
func (soc *Socket) Close() error {
	if i, err := C.zmq_close(soc.soc); int(i) != 0 {
		return errget(err)
	}
	soc.soc = unsafe.Pointer(nil)
	return nil
}

/*
Accept incoming connections on a socket.

For a description of endpoint, see: http://api.zeromq.org/2-2:zmq-bind#toc2
*/
func (soc *Socket) Bind(endpoint string) error {
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_bind(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Create outgoing connection from socket.

For a description of endpoint, see: http://api.zeromq.org/2-2:zmq-connect#toc2
*/
func (soc *Socket) Connect(endpoint string) error {
	s := C.CString(endpoint)
	defer C.free(unsafe.Pointer(s))
	if i, err := C.zmq_connect(soc.soc, s); int(i) != 0 {
		return errget(err)
	}
	return nil
}

/*
Receive a message part from a socket.

For a description of flags, see: http://api.zeromq.org/2-2:zmq-recv#toc2
*/
func (soc *Socket) Recv(flags Flag) (string, error) {
	b, err := soc.RecvBytes(flags)
	return string(b), err
}

/*
Receive a message part from a socket.

For a description of flags, see: http://api.zeromq.org/2-2:zmq-recv#toc2
*/
func (soc *Socket) RecvBytes(flags Flag) ([]byte, error) {
	var msg C.zmq_msg_t
	if i, err := C.zmq_msg_init(&msg); i != 0 {
		return []byte{}, errget(err)
	}
	defer C.zmq_msg_close(&msg)

	var size C.int
	var err error

	var i C.int
	i, err = C.zmq_recv(soc.soc, &msg, C.int(flags))
	if i == 0 {
		size = C.int(C.zmq_msg_size(&msg))
	} else {
		size = -1
	}

	if size < 0 {
		return []byte{}, errget(err)
	}
	if size == 0 {
		return []byte{}, nil
	}
	data := make([]byte, int(size))
	C.memcpy(unsafe.Pointer(&data[0]), C.zmq_msg_data(&msg), C.size_t(size))
	return data, nil
}

/*
Send a message part on a socket.

For a description of flags, see: http://api.zeromq.org/2-2:zmq-send#toc2
*/
func (soc *Socket) Send(data string, flags Flag) (int, error) {
	return soc.SendBytes([]byte(data), flags)
}

/*
Send a message part on a socket.

For a description of flags, see: http://api.zeromq.org/2-2:zmq-send#toc2
*/
func (soc *Socket) SendBytes(data []byte, flags Flag) (int, error) {
	datac := C.CString(string(data))
	var msg C.zmq_msg_t
	if i, err := C.my_msg_init_data(&msg, unsafe.Pointer(datac), C.size_t(len(data))); i != 0 {
		return -1, errget(err)
	}
	defer C.zmq_msg_close(&msg)
	n, err := C.zmq_send(soc.soc, &msg, C.int(flags))
	if n != 0 {
		return -1, errget(err)
	}
	return int(n), nil
}

/*
Start built-in ØMQ device

see: http://api.zeromq.org/2-2:zmq-device#toc2
*/
func Device(device Dev, frontend, backend *Socket) error {
	_, err := C.zmq_device(C.int(device), frontend.soc, backend.soc)
	return errget(err)
}

/*
Emulate the proxy that will be built-in in 0MQ version 3

See: http://api.zeromq.org/3-2:zmq-proxy
*/
func Proxy(frontend, backend, capture *Socket) error {
	items := NewPoller()
	items.Add(frontend, POLLIN)
	items.Add(backend, POLLIN)
	for {
		sockets, err := items.Poll(-1)
		if err != nil {
			return err
		}
		for _, socket := range sockets {
			for more := true; more; {
				msg, err := socket.Socket.RecvBytes(0)
				if err != nil {
					return err
				}
				more, err = socket.Socket.GetRcvmore()
				if err != nil {
					return err
				}
				fl := SNDMORE
				if !more {
					fl = 0
				}

				if capture != nil {
					_, err = capture.SendBytes(msg, fl)
					if err != nil {
						return err
					}
				}

				switch socket.Socket {
				case frontend:
					_, err = backend.SendBytes(msg, fl)
				case backend:
					_, err = frontend.SendBytes(msg, fl)
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
