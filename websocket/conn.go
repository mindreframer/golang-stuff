// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package websocket implements the WebSocket protocol defined in RFC 6455.
//
// Overview
//
// The Conn type represents a WebSocket connection. WebSocket messages are
// represented by the io.Reader interface when receiving a message and by the
// io.WriteCloser interface when sending a message. An application receives a
// message by calling the Conn.NextReader method and reading the returned
// io.Reader to EOF. An application sends a message by calling the
// Conn.NextWriter method and writing the message to the returned
// io.WriteCloser. The application terminates the message by closing the
// io.WriteCloser.
//
// The following example shows how to use NextReader and NextWriter to echo
// messages:
//
//	for {
//      op, r, err := conn.NextReader()
//      if err != nil {
//			return
//      }
//		if op != websocket.OpBinary && op != websocket.OpText {
//          // Ignore if not a data message.
//			continue
//		}
//		w, err := conn.NextWriter(op)
//		if err != nil {
//			return err
//		}
//		if _, err := io.Copy(w, r); err != nil {
//          return err
//      }
//      if err := w.Close(); err != nil {
//          return err
//      }
//	}
//
// Concurrency
//
// A Conn supports a single concurrent caller to the write methods (NextWriter,
// SetWriteDeadline, WriteMessage) and a single concurrent caller to the read
// methods (NextReader, SetReadDeadline). The Close and WriteControl methods
// can be called concurrently with all other methods.
//
// Text
//
// Text messages in the WebSocket protocol are transmitted as UTF-8. It is the
// application's responsibility to ensure that text messages are valid UTF-8.
package websocket

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Close codes defined in RFC 6455, section 11.7.
const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseTLSHandshake            = 1015
)

// Opcodes defined in RFC 6455, section 11.8.
const (
	OpContinuation = 0
	OpText         = 1
	OpBinary       = 2
	OpClose        = 8
	OpPing         = 9
	OpPong         = 10
)

var (
	ErrCloseSent = errors.New("websocket: close sent")
	ErrReadLimit = errors.New("websocket: read limit exceeded")
)

var (
	errBadWriteOpCode      = errors.New("websocket: bad write opcode")
	errWriteTimeout        = errors.New("websocket: write timeout")
	errWriteClosed         = errors.New("websocket: write closed")
	errInvalidControlFrame = errors.New("websocket: invalid control frame")
)

const (
	maxFrameHeaderSize         = 2 + 8 + 4 // Fixed header + length + mask
	maxControlFramePayloadSize = 125
	finalBit                   = 1 << 7
	maskBit                    = 1 << 7
	writeWait                  = time.Second
)

func maskBytes(key [4]byte, pos int, b []byte) int {
	for i := range b {
		b[i] ^= key[pos&3]
		pos += 1
	}
	return pos & 3
}

func newMaskKey() [4]byte {
	n := rand.Uint32()
	return [4]byte{byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 32)}
}

// Conn represents a WebSocket connection.
type Conn struct {
	conn     net.Conn
	isServer bool

	// Write fields
	mu        chan bool // used as mutex to protect write to conn and closeSent
	closeSent bool      // true if close message was sent

	// Message writer fields.
	writeErr      error
	writeBuf      []byte // frame is constructed in this buffer.
	writePos      int    // end of data in writeBuf.
	writeOpCode   int    // op code for the current frame.
	writeSeq      int    // incremented to invalidate message writers.
	writeDeadline time.Time

	// Read fields
	readErr       error
	br            *bufio.Reader
	readRemaining int64 // bytes remaining in current frame.
	readFinal     bool  // true the current message has more frames.
	readSeq       int   // incremented to invalidate message readers.
	readLength    int64 // Message size.
	readLimit     int64 // Maximum message size.
	readMaskPos   int
	readMaskKey   [4]byte
	savedPong     []byte
}

func newConn(conn net.Conn, isServer bool, readBufSize, writeBufSize int) *Conn {
	mu := make(chan bool, 1)
	mu <- true

	return &Conn{
		isServer:    isServer,
		br:          bufio.NewReaderSize(conn, readBufSize),
		conn:        conn,
		mu:          mu,
		readFinal:   true,
		writeBuf:    make([]byte, writeBufSize+maxFrameHeaderSize),
		writeOpCode: -1,
		writePos:    maxFrameHeaderSize,
	}
}

// Close closes the underlying network connection without sending or waiting for a close frame.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// Write methods

func (c *Conn) write(opCode int, deadline time.Time, bufs ...[]byte) error {
	<-c.mu
	defer func() { c.mu <- true }()

	if c.closeSent {
		return ErrCloseSent
	} else if opCode == OpClose {
		c.closeSent = true
	}

	c.conn.SetWriteDeadline(deadline)
	for _, buf := range bufs {
		if len(buf) > 0 {
			n, err := c.conn.Write(buf)
			if n != len(buf) {
				// Close on partial write.
				c.conn.Close()
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteControl writes a control message with the given deadline. The allowed
// opCodes are OpClose, OpPing and OpPong.
func (c *Conn) WriteControl(opCode int, data []byte, deadline time.Time) error {
	if opCode != OpClose && opCode != OpPing && opCode != OpPong {
		return errBadWriteOpCode
	}
	if len(data) > maxControlFramePayloadSize {
		return errInvalidControlFrame
	}

	b0 := byte(opCode) | finalBit
	b1 := byte(len(data))
	if !c.isServer {
		b1 |= maskBit
	}

	buf := make([]byte, 0, maxFrameHeaderSize+maxControlFramePayloadSize)
	buf = append(buf, b0, b1)

	if c.isServer {
		buf = append(buf, data...)
	} else {
		key := newMaskKey()
		buf = append(buf, key[:]...)
		buf = append(buf, data...)
		maskBytes(key, 0, buf[6:])
	}

	d := time.Hour * 1000
	if !deadline.IsZero() {
		d = deadline.Sub(time.Now())
		if d < 0 {
			return errWriteTimeout
		}
	}

	timer := time.NewTimer(d)
	select {
	case <-c.mu:
		timer.Stop()
	case <-timer.C:
		return errWriteTimeout
	}
	defer func() { c.mu <- true }()

	if c.closeSent {
		return ErrCloseSent
	} else if opCode == OpClose {
		c.closeSent = true
	}

	c.conn.SetWriteDeadline(deadline)
	n, err := c.conn.Write(buf)
	if n != 0 && n != len(buf) {
		c.conn.Close()
	}
	return err
}

// NextWriter returns a writer for the next message to send. The allowed
// opCodes are OpText, OpBinary, OpClose and OpPing. The writer's Close method
// flushes the complete message to the network.
//
// There can be at most one open writer on a connection. NextWriter closes the
// previous writer if the application has not already done so.
//
// The NextWriter method and the writers returned from the method cannot be
// accessed by more than one goroutine at a time.
func (c *Conn) NextWriter(opCode int) (io.WriteCloser, error) {
	if c.writeErr != nil {
		return nil, c.writeErr
	}

	if c.writeOpCode != -1 {
		if err := c.flushFrame(true, nil); err != nil {
			return nil, err
		}
	}

	if opCode != OpText && opCode != OpBinary && opCode != OpClose && opCode != OpPing {
		return nil, errBadWriteOpCode
	}

	c.writeOpCode = opCode
	return messageWriter{c, c.writeSeq}, nil
}

func (c *Conn) flushFrame(final bool, extra []byte) error {
	length := c.writePos - maxFrameHeaderSize + len(extra)

	// Check for invalid control frames.
	if (c.writeOpCode == OpClose || c.writeOpCode == OpPing) &&
		(!final || length > maxControlFramePayloadSize) {
		c.writeSeq += 1
		c.writeOpCode = -1
		c.writePos = maxFrameHeaderSize
		return errInvalidControlFrame
	}

	b0 := byte(c.writeOpCode)
	if final {
		b0 |= finalBit
	}
	b1 := byte(0)
	if !c.isServer {
		b1 |= maskBit
	}

	// Assume that the frame starts at beginning of c.writeBuf.
	framePos := 0
	if c.isServer {
		// Adjust up if mask not included in the header.
		framePos = 4
	}

	switch {
	case length >= 65536:
		c.writeBuf[framePos] = b0
		c.writeBuf[framePos+1] = b1 | 127
		binary.BigEndian.PutUint64(c.writeBuf[framePos+2:], uint64(length))
	case length > 125:
		framePos += 6
		c.writeBuf[framePos] = b0
		c.writeBuf[framePos+1] = b1 | 126
		binary.BigEndian.PutUint16(c.writeBuf[framePos+2:], uint16(length))
	default:
		framePos += 8
		c.writeBuf[framePos] = b0
		c.writeBuf[framePos+1] = b1 | byte(length)
	}

	if !c.isServer {
		key := newMaskKey()
		copy(c.writeBuf[maxFrameHeaderSize-4:], key[:])
		maskBytes(key, 0, c.writeBuf[maxFrameHeaderSize:c.writePos])
		if len(extra) > 0 {
			c.writeErr = errors.New("websocket: internal error, extra used in client mode")
			return c.writeErr
		}
	}

	// Write the buffers to the connection.
	c.writeErr = c.write(c.writeOpCode, c.writeDeadline, c.writeBuf[framePos:c.writePos], extra)

	// Setup for next frame.
	c.writePos = maxFrameHeaderSize
	c.writeOpCode = OpContinuation
	if final {
		c.writeSeq += 1
		c.writeOpCode = -1
	}
	return c.writeErr
}

type messageWriter struct {
	c   *Conn
	seq int
}

func (w messageWriter) err() error {
	c := w.c
	if c.writeSeq != w.seq {
		return errWriteClosed
	}
	if c.writeErr != nil {
		return c.writeErr
	}
	return nil
}

func (w messageWriter) ncopy(max int) (int, error) {
	n := len(w.c.writeBuf) - w.c.writePos
	if n <= 0 {
		if err := w.c.flushFrame(false, nil); err != nil {
			return 0, err
		}
		n = len(w.c.writeBuf) - w.c.writePos
	}
	if n > max {
		n = max
	}
	return n, nil
}

func (w messageWriter) write(final bool, p []byte) (int, error) {
	if err := w.err(); err != nil {
		return 0, err
	}

	if len(p) > 2*len(w.c.writeBuf) && w.c.isServer {
		// Don't buffer large messages.
		err := w.c.flushFrame(final, p)
		if err != nil {
			return 0, err
		}
		return len(p), nil
	}

	nn := len(p)
	for len(p) > 0 {
		n, err := w.ncopy(len(p))
		if err != nil {
			return 0, err
		}
		copy(w.c.writeBuf[w.c.writePos:], p[:n])
		w.c.writePos += n
		p = p[n:]
	}
	return nn, nil
}

func (w messageWriter) Write(p []byte) (int, error) {
	return w.write(false, p)
}

func (w messageWriter) WriteString(p string) (int, error) {
	if err := w.err(); err != nil {
		return 0, err
	}

	nn := len(p)
	for len(p) > 0 {
		n, err := w.ncopy(len(p))
		if err != nil {
			return 0, err
		}
		copy(w.c.writeBuf[w.c.writePos:], p[:n])
		w.c.writePos += n
		p = p[n:]
	}
	return nn, nil
}

func (w messageWriter) ReadFrom(r io.Reader) (nn int64, err error) {
	if err := w.err(); err != nil {
		return 0, err
	}
	for {
		if w.c.writePos == len(w.c.writeBuf) {
			err = w.c.flushFrame(false, nil)
			if err != nil {
				break
			}
		}
		var n int
		n, err = r.Read(w.c.writeBuf[w.c.writePos:])
		w.c.writePos += n
		nn += int64(n)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return nn, err
}

func (w messageWriter) Close() error {
	if err := w.err(); err != nil {
		return err
	}
	return w.c.flushFrame(true, nil)
}

// WriteMessage is a helper method for getting a writer using NextWriter,
// writing the message and closing the writer.
func (c *Conn) WriteMessage(opCode int, data []byte) error {
	wr, err := c.NextWriter(opCode)
	if err != nil {
		return err
	}
	w := wr.(messageWriter)
	if _, err := w.write(true, data); err != nil {
		return err
	}
	if c.writeSeq == w.seq {
		if err := c.flushFrame(true, nil); err != nil {
			return err
		}
	}
	return nil
}

// SetWriteDeadline sets the deadline for future calls to NextWriter and the
// io.WriteCloser returned from NextWriter. If the deadline is reached, the
// call will fail with a timeout instead of blocking. A zero value for t means
// Write will not time out. Even if Write times out, it may return n > 0,
// indicating that some of the data was successfully written.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline = t
	return nil
}

// Read methods

func (c *Conn) advanceFrame() (int, error) {

	// 1. Skip remainder of previous frame.

	if c.readRemaining > 0 {
		if _, err := io.CopyN(ioutil.Discard, c.br, c.readRemaining); err != nil {
			return -1, err
		}
	}

	// 2. Read and parse first two bytes of frame header.

	var b [8]byte
	if err := c.read(b[:2]); err != nil {
		return -1, err
	}

	final := b[0]&finalBit != 0
	opCode := int(b[0] & 0xf)
	reserved := int((b[0] >> 4) & 0x7)
	mask := b[1]&maskBit != 0
	c.readRemaining = int64(b[1] & 0x7f)

	if reserved != 0 {
		return -1, c.handleProtocolError("unexpected reserved bits " + strconv.Itoa(reserved))
	}

	switch opCode {
	case OpClose, OpPing, OpPong:
		if c.readRemaining > maxControlFramePayloadSize {
			return -1, c.handleProtocolError("control frame length > 125")
		}
		if !final {
			return -1, c.handleProtocolError("control frame not final")
		}
	case OpText, OpBinary:
		if !c.readFinal {
			return -1, c.handleProtocolError("message start before final message frame")
		}
		c.readFinal = final
	case OpContinuation:
		if c.readFinal {
			return -1, c.handleProtocolError("continuation after final message frame")
		}
		c.readFinal = final
	default:
		return -1, c.handleProtocolError("unknown opcode " + strconv.Itoa(opCode))
	}

	// 3. Read and parse frame length.

	switch c.readRemaining {
	case 126:
		if err := c.read(b[:2]); err != nil {
			return -1, err
		}
		c.readRemaining = int64(binary.BigEndian.Uint16(b[:2]))
	case 127:
		if err := c.read(b[:8]); err != nil {
			return -1, err
		}
		c.readRemaining = int64(binary.BigEndian.Uint64(b[:8]))
	}

	// 4. Handle frame masking.

	if mask != c.isServer {
		return -1, c.handleProtocolError("incorrect mask flag")
	}

	if mask {
		c.readMaskPos = 0
		if err := c.read(c.readMaskKey[:]); err != nil {
			return -1, err
		}
	}

	// 5. For text and binary messages, enforce read limit and return.

	if opCode == OpContinuation || opCode == OpText || opCode == OpBinary {

		c.readLength += c.readRemaining
		if c.readLimit > 0 && c.readLength > c.readLimit {
			c.WriteControl(OpClose, FormatCloseMessage(CloseMessageTooBig, ""), time.Now().Add(writeWait))
			return -1, ErrReadLimit
		}

		return opCode, nil
	}

	// 6. Read control frame payload.

	payload := make([]byte, c.readRemaining)
	c.readRemaining = 0
	if err := c.read(payload); err != nil {
		return -1, err
	}
	maskBytes(c.readMaskKey, 0, payload)

	// 7. Process control frame payload.

	switch opCode {
	case OpPong:
		c.savedPong = payload
	case OpPing:
		c.WriteControl(OpPong, payload, time.Now().Add(writeWait))
	case OpClose:
		c.WriteControl(OpClose, []byte{}, time.Now().Add(writeWait))
		if len(payload) < 2 {
			return -1, io.EOF
		} else {
			closeCode := binary.BigEndian.Uint16(payload)
			switch closeCode {
			case CloseNormalClosure, CloseGoingAway:
				return -1, io.EOF
			default:
				return -1, errors.New("websocket: close " +
					strconv.Itoa(int(closeCode)) + " " +
					string(payload[2:]))
			}
		}
	}

	return opCode, nil
}

func (c *Conn) handleProtocolError(message string) error {
	c.WriteControl(OpClose, FormatCloseMessage(CloseProtocolError, message), time.Now().Add(writeWait))
	return errors.New("websocket: " + message)
}

func (c *Conn) read(buf []byte) error {
	var err error
	for len(buf) > 0 && err == nil {
		var nn int
		nn, err = c.br.Read(buf)
		buf = buf[nn:]
	}
	if err == io.EOF {
		if len(buf) == 0 {
			err = nil
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	return err
}

// NextReader returns the next message received from the peer. The returned
// opCode is one of OpText, OpBinary or OpPong. The connection automatically
// handles ping messages received from the peer. NextReader returns an error
// upon receiving a close message from the peer.
//
// There can be at most one open reader on a connection. NextReader discards
// the previous message if the application has not already consumed it.
//
// The NextReader method and the readers returned from the method cannot be
// accessed by more than one goroutine at a time.
func (c *Conn) NextReader() (opCode int, r io.Reader, err error) {

	c.readSeq += 1
	c.readLength = 0

	if c.savedPong != nil {
		r := bytes.NewReader(c.savedPong)
		c.savedPong = nil
		return OpPong, r, nil
	}

	for c.readErr == nil {
		var opCode int
		opCode, c.readErr = c.advanceFrame()
		switch opCode {
		case OpText, OpBinary:
			return opCode, messageReader{c, c.readSeq}, nil
		case OpPong:
			r := bytes.NewReader(c.savedPong)
			c.savedPong = nil
			return OpPong, r, nil
		case OpContinuation:
			// do nothing
		}
	}
	return -1, nil, c.readErr
}

type messageReader struct {
	c   *Conn
	seq int
}

func (r messageReader) Read(b []byte) (n int, err error) {

	if r.seq != r.c.readSeq {
		return 0, io.EOF
	}

	for r.c.readErr == nil {

		if r.c.readRemaining > 0 {
			if int64(len(b)) > r.c.readRemaining {
				b = b[:r.c.readRemaining]
			}
			r.c.readErr = r.c.read(b)
			r.c.readMaskPos = maskBytes(r.c.readMaskKey, r.c.readMaskPos, b)
			r.c.readRemaining -= int64(len(b))
			return len(b), r.c.readErr
		}

		if r.c.readFinal {
			r.c.readSeq += 1
			return 0, io.EOF
		}

		var opCode int
		opCode, r.c.readErr = r.c.advanceFrame()

		if opCode == OpText || opCode == OpBinary {
			r.c.readErr = errors.New("websocket: internal error, unexpected text or binary in Reader")
		}
	}
	return 0, r.c.readErr
}

// SetReadDeadline sets the deadline for future calls to NextReader and the
// io.Reader returned from NextReader. If the deadline is reached, the call
// will fail with a timeout instead of blocking. A zero value for t means that
// the methods will not time out.
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetReadLimit sets the maximum size for a message read from the peer. If a
// message exceeds the limit, the connection sends a close frame to the peer
// and returns ErrReadLimit to the application.
func (c *Conn) SetReadLimit(limit int64) {
	c.readLimit = limit
}

// FormatCloseMessage formats closeCode and text as a WebSocket close message.
func FormatCloseMessage(closeCode int, text string) []byte {
	buf := make([]byte, 2+len(text))
	binary.BigEndian.PutUint16(buf, uint16(closeCode))
	copy(buf[2:], text)
	return buf
}
