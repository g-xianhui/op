package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/g-xianhui/crypt"
	"math/big"
	"net"
	"os"
	"strconv"
)

func auth(conn net.Conn) (secret []byte, err error) {
	challenge, err := readPack(conn)
	if err != nil {
		return
	}
	log(DEBUG, "challenge: %X\n", challenge)

	a := crypt.Randomkey()
	A := crypt.DHExchange(a)
	err = writePack(conn, A.Bytes())
	if err != nil {
		return
	}
	log(DEBUG, "A: %X\n", A.Bytes())

	B, err := readPack(conn)
	if err != nil {
		return
	}
	log(DEBUG, "B: %X\n", B)

	z := new(big.Int)
	z.SetBytes(B)
	s := crypt.DHSecret(a, z)
	log(DEBUG, "secret: %X\n", s.Bytes())

	mac := hmac.New(sha256.New, s.Bytes())
	mac.Write(challenge)
	err = writePack(conn, mac.Sum(nil))
	if err != nil {
		return
	}

	if len(s.Bytes()) < 16 {
		err = errors.New("secret length less than 16 bytes")
		return
	}

	secret = s.Bytes()[:16]
	return
}

func login(conn net.Conn, accountName string) (session uint32, secret []byte, err error) {
	writePack(conn, []byte("login"))

	secret, err = auth(conn)
	if err != nil {
		return
	}

	err = writeEncrypt(conn, []byte(accountName), secret)
	if err != nil {
		log(ERROR, "aes failed\n")
		return
	}

	data, err := readEncrypt(conn, secret)
	if err != nil {
		return
	}
	n, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		return
	}
	session = uint32(n)

	return
}

func reconnect(accountName string, secret []byte) (session uint32, conn net.Conn, err error) {
	conn, err = net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("err:", err.Error())
		os.Exit(1)
	}

	cmd := fmt.Sprintf("reconnect:%s", accountName)
	err = writePack(conn, []byte(cmd))
	if err != nil {
		return
	}

	challenge, err := readPack(conn)
	if err != nil {
		return
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write(challenge)
	err = writePack(conn, mac.Sum(nil))
	if err != nil {
		return
	}

	data, err := readEncrypt(conn, secret)
	if err != nil {
		return
	}
	n, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		return
	}
	session = uint32(n)
	return
}
