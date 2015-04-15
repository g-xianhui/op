package main

import (
	"encoding/binary"
	"errors"
	"io"
)

var _ = io.EOF

const MAX_CLIENT_BUF = 1 << 16

type msg struct {
	t       uint32
	session uint32
	data    []byte
}

type cacheBuf struct {
	buf    []byte
	curPos int
}

// read a package from conn, package is format like: | length(2) | data(lenght-2) |
func readPack(conn io.Reader, last *cacheBuf) (pack []byte, err error) {
	if last.curPos >= MAX_CLIENT_BUF {
		err = errors.New("client buf overflow")
		return
	}

	n, err := conn.Read(last.buf[last.curPos:])
	if err != nil {
		return
	}

	last.curPos += n
	if last.curPos < 2 {
		return
	}

	packLen := int(binary.BigEndian.Uint16(last.buf[:2]))
	if last.curPos < packLen {
		return
	}

	pack = make([]byte, packLen-2)
	copy(pack, last.buf[2:packLen])

	copy(last.buf, last.buf[packLen:last.curPos])
	last.curPos -= packLen

	return
}

// send a package to conn
func sendPack(conn io.Writer, pack []byte) error {
	l := len(pack)
	if l > MAX_CLIENT_BUF {
		return errors.New("msg is too big")
	}
	lbuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lbuf, uint16(l+2))
	if _, err := conn.Write(lbuf); err != nil {
		return err
	}
	if _, err := conn.Write(pack); err != nil {
		return err
	}
	return nil
}

func packMsg(m *msg) []byte {
	l := len(m.data) + 8
	pack := make([]byte, l)
	binary.BigEndian.PutUint32(pack, m.t)
	binary.BigEndian.PutUint32(pack[4:], m.session)
	copy(pack[8:], m.data)
	return pack
}

func unpackMsg(pack []byte) *msg {
	m := &msg{}
	m.t = binary.BigEndian.Uint32(pack[:4])
	m.session = binary.BigEndian.Uint32(pack[4:8])
	m.data = pack[8:]
	return m
}

func netio(agent *Agent) {
	last := &cacheBuf{}
	last.buf = make([]byte, MAX_CLIENT_BUF)

	for {
		pack, err := readPack(agent.conn, last)
		if err != nil {
			if err != io.EOF {
				log("read from client err: %s", err)
			} else {
				log("client end")
			}
			break
		}

		m := unpackMsg(pack)
		agent.outside <- m
	}
}
