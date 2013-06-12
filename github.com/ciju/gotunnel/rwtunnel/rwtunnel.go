package rwtunnel

import (
	l "github.com/ciju/gotunnel/log"
	"io"
)

func copyFromTo(a, b io.ReadWriteCloser) {
	defer func() {
		a.Close()
	}()
	io.Copy(a, b)
}

type RWTunnel struct {
	src, dst io.ReadWriteCloser
}

func (p *RWTunnel) Proxy() {
	// go copypaste(p.src, p.dst, false, "f")
	// go copypaste(p.dst, p.src, true, "b")
	go copyFromTo(p.src, p.dst)
	go copyFromTo(p.dst, p.src)
}

func NewRWTunnel(src, dst io.ReadWriteCloser) (p *RWTunnel) {
	b := &RWTunnel{src: src, dst: dst}
	b.Proxy()
	return b
}

// // only the actual host/port client should be able to close a
// // connection. ex: Keep-Alive and websockets.

func copypaste(in, out io.ReadWriteCloser, close_in bool, msg string) {
	var buf [512]byte

	defer func() {
		if close_in {
			l.Log("eof closing connection")
			in.Close()
			out.Close()
		}
	}()

	for {
		n, err := in.Read(buf[0:])
		// on readerror, only bail if no other choice.
		if err == io.EOF {
			l.Log("msg: ", msg)
			// fmt.Print(msg)
			// time.Sleep(1e9)
			l.Log("eof", msg)
			return
		}
		l.Log("-- read ", n)
		if err != nil {
			l.Log("something wrong while copying in ot out ", msg)
			l.Log("error: ", err)
			return
		}
		// if n < 1 {
		// 	fmt.Println("nothign to read")
		// 	return
		// }

		l.Log("-- wrote msg bytes", n)

		_, err = out.Write(buf[0:n])
		if err != nil {
			l.Log("something wrong while copying out to in ")
			// l.Fatal("something wrong while copying out to in", err)
			return
		}
	}
}
