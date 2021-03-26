package wl

import (
	"bytes"
	"errors"
	"github.com/yalue/native_endian"
	"syscall"
)

// Event is the Wayland event (e.g. a response) from the compositor
type Event struct {
	Pid    ProxyId
	Opcode uint32
	Data   []byte
	scms   []syscall.SocketControlMessage
	off    int
	err    error
}

/*
	Okay, so if you pass a file descriptor across a
	UNIX domain socket, you may actually receive it
	on an earlier call to recvmsg. If we don't do
	anything about this we end up getting a file
	descriptor on the wrong wayland protocol message.

	See https://keithp.com/blogs/fd-passing/

	This is a hacky solution:
	- have a map from client fd to list of receive fds
	- when we receive a fd over the socket add it to
	  the clients list of fds
	- when we call event.FD() take the earliest FD
	  from the list

	UPDATE: Scratch the above. The actual solution is to
	store a lists of incoming fds in the Context
	and modify the generator to set FDs like:
		ev.Fd = p.Context().NextFD()
*/

// Error unable to read message header is returned when it is not possible to read enough bytes from the unix socket, use InternalError to get the underlying cause
var ErrReadHeader = errors.New("unable to read message header")

// Error size of message header is wrong is returned when the returned size of message heaer is not 8 bytes
var ErrSizeOfHeaderWrong = errors.New("size of message header is wrong")

// Error unsufficient control msg buffer is returned when the oobn is bigger than the control message buffer
var ErrControlMsgBuffer = errors.New("unsufficient control msg buffer")

// Error control message parse error is returned when the unix socket control message cannot be parsed, use InternalError to get the underlying cause
var ErrControlMsgParseError = errors.New("control message parse error")

// Error invalid message size is returned when the payload message size read from the unix socket is incorrect
var ErrInvalidMsgSize = errors.New("invalid message size")

// Error cannot read message is returned when the payload message cannot be read, use InternalError to get the underlying cause
var ErrReadPayload = errors.New("cannot read message")

func (c *Context) readEvent() (*Event, error) {
	buf := bytePool.Take(8)
	control := bytePool.Take(24)

	n, oobn, _, _, err := c.conn.ReadMsgUnix(buf[:], control)
	if err != nil {
		return nil, combinedError{ErrReadHeader, err}
	}
	if n != 8 {
		return nil, ErrSizeOfHeaderWrong
	}
	ev := new(Event)
	if oobn > 0 {
		if oobn > len(control) {
			return nil, ErrControlMsgBuffer
		}
		scms, err := syscall.ParseSocketControlMessage(control)
		if err != nil {
			return nil, combinedError{ErrControlMsgParseError, err}
		}
		ev.scms = scms
	}

	ev.Pid = ProxyId(native_endian.NativeEndian().Uint32(buf[0:4]))
	ev.Opcode = uint32(native_endian.NativeEndian().Uint16(buf[4:6]))
	size := uint32(native_endian.NativeEndian().Uint16(buf[6:8]))

	// subtract 8 bytes from header
	data := bytePool.Take(int(size) - 8)
	n, err = c.conn.Read(data)
	if err != nil {
		return nil, combinedError{ErrReadPayload, err}
	}
	if n != int(size)-8 {
		return nil, ErrInvalidMsgSize
	}
	ev.Data = data

	bytePool.Give(buf)
	bytePool.Give(control)

	return ev, nil
}

// Error no socket control messages
var ErrNoControlMsgs = errors.New("no socket control messages")

// Error unable to parse unix rights
var ErrUnableToParseUnixRights = errors.New("unable to parse unix rights")

func (ev *Event) FD() (uintptr, error) {
	if ev.scms == nil {
		return 0, ErrNoControlMsgs
	}
	fds, err := syscall.ParseUnixRights(&ev.scms[0])
	if err != nil {
		return 0, ErrUnableToParseUnixRights
	}
	//TODO: is this required??????????????
	ev.scms = append(ev.scms, ev.scms[1:]...)
	return uintptr(fds[0]), nil
}

// Error unable to read unsigned int is returned when the buffer is too short to contain a specific unsigned int
var ErrUnableToParseUint32 = errors.New("unable to read unsigned int")

func (ev *Event) Uint32() uint32 {
	buf := ev.next(4)
	if len(buf) != 4 {
		ev.err = ErrUnableToParseUint32
		return 0
	}
	return native_endian.NativeEndian().Uint32(buf)
}

// Event Proxy decodes Proxy by it's Id from the Event
func (ev *Event) Proxy(c *Context) Proxy {
	id := ev.Uint32()
	if id == 0 {
		return nil
	} else {
		return c.LookupProxy(ProxyId(id))
	}
}

// Error unable to parse string is returned when the buffer is too short to contain a specific string
var ErrUnableToParseString = errors.New("unable to parse string")

// Event String decodes a string from the Event
func (ev *Event) String() string {
	l := int(ev.Uint32())
	buf := ev.next(l)
	if len(buf) != l {
		ev.err = ErrUnableToParseString
		return ""
	}
	ret := string(bytes.TrimRight(buf, "\x00"))
	//padding to 32 bit boundary
	if (l & 0x3) != 0 {
		ev.next(4 - (l & 0x3))
	}
	return ret
}

// Event Int32 decodes an Int32 from the Event
func (ev *Event) Int32() int32 {
	return int32(ev.Uint32())
}

// Event Float32 decodes a Float32 from the Event
func (ev *Event) Float32() float32 {
	return float32(FixedToFloat(ev.Int32()))
}

// Event Array decodes an Array from the Event
func (ev *Event) Array() []int32 {
	l := int(ev.Uint32())
	arr := make([]int32, l/4)
	for i := range arr {
		arr[i] = ev.Int32()
	}
	return arr
}

func (ev *Event) next(n int) []byte {
	ret := ev.Data[ev.off : ev.off+n]
	ev.off += n
	return ret
}
