// +build !windows

package zmq2

/*
#include <zmq.h>
*/
import "C"

// ZMQ_FD: Retrieve file descriptor associated with the socket
//
// See: http://api.zeromq.org/2-2:zmq-getsockopt#toc21
func (soc *Socket) GetFd() (int, error) {
	return soc.getInt(C.ZMQ_FD)
}
