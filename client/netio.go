package main

import (
	"encoding/binary"
	"errors"
	"github.com/g-xianhui/crypt"
	"io"
)

var _ = io.EOF

const MAX_CLIENT_BUF = 1 << 16

func Readn(conn io.Reader, count int) ([]byte, int, error) {
	data := make([]byte, count)
	n := 0
	for {
		i, err := conn.Read(data[n:])
		if i > 0 {
			n += i
		}

		if n == count || err != nil {
			return data, n, err
		}
	}
}

// read a package from conn, package is format like: | length(2) | data(lenght) |
func readPack(conn io.Reader) ([]byte, error) {
	lenbuf, n, err := Readn(conn, 2)
	if n < 2 {
		if err == io.EOF && n == 0 {
			return nil, err
		} else {
			log(ERROR, "Readn: %s\n", err)
			return nil, errors.New("uncomplete head")
		}
	}
	textLen := int(binary.BigEndian.Uint16(lenbuf))

	data, n, err := Readn(conn, textLen)
	if n < textLen {
		log(ERROR, "Readn: %s\n", err)
		return nil, errors.New("uncomplete content")
	}
	return data, nil
}

// write a package to conn
func writePack(conn io.Writer, data []byte) error {
	l := len(data)
	if l >= MAX_CLIENT_BUF {
		return errors.New("msg is too big")
	}

	pack := make([]byte, l+2)
	binary.BigEndian.PutUint16(pack[:2], uint16(l))
	copy(pack[2:], data)
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
	if len(pack) < 8 {
		return nil
	}
	m := &msg{}
	m.t = binary.BigEndian.Uint32(pack[:4])
	m.session = binary.BigEndian.Uint32(pack[4:8])
	m.data = pack[8:]
	return m
}

func recv(agent *Agent) {
	for {
		pack, err := readEncrypt(agent.conn, agent.secret)
		if pack != nil {
			m := unpackMsg(pack)
			if m == nil {
				continue
			}
			agent.outside <- m
		}

		if err != nil {
			if err != io.EOF {
				log(ERROR, "read from server err: %s\n", err)
			} else {
				log(DEBUG, "server closed\n")
			}
			break
		}
	}
}

func readEncrypt(conn io.Reader, secret []byte) ([]byte, error) {
	ciphertext, err := readPack(conn)
	if err != nil {
		return nil, err
	}

	data, err := crypt.AesDecrypt(ciphertext, secret)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeEncrypt(conn io.Writer, data []byte, secret []byte) error {
	ciphertext, err := crypt.AesEncrypt(data, secret)
	if err != nil {
		return err
	}
	return writePack(conn, ciphertext)
}
