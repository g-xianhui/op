package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/g-xianhui/crypt"
	"math/big"
	"math/rand"
	"net"
	"strconv"
)

func challengeCheck(challenge []byte, secret []byte, response []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(challenge)
	return hmac.Equal(response, mac.Sum(nil))
}

func dhsecret(conn net.Conn) (secret []byte, err error) {
	challenge := make([]byte, 8)
	binary.BigEndian.PutUint64(challenge, uint64(rand.Int63()))
	err = writePack(conn, challenge)
	if err != nil {
		return
	}

	B, err := readPack(conn)
	if err != nil {
		return
	}

	a := crypt.Randomkey()
	A := crypt.DHExchange(a)
	err = writePack(conn, A.Bytes())
	if err != nil {
		return
	}

	z := new(big.Int)
	z.SetBytes(B)
	s := crypt.DHSecret(a, z)

	response, err := readPack(conn)
	if err != nil {
		return
	}

	if !challengeCheck(challenge, s.Bytes(), response) {
		err = errors.New("challenge failed")
		return
	}

	secret = s.Bytes()[:16]
	return
}

func auth(conn net.Conn, secret []byte) (accountName string, session uint32, err error) {
	// success whatever for test
	data, err := readEncrypt(conn, secret)
	if err != nil {
		return
	}
	accountName = string(data)

	session = rand.Uint32()
	str := strconv.FormatInt(int64(session), 10)
	err = writeEncrypt(conn, []byte(str), secret)
	if err != nil {
		return
	}
	return
}

func login(conn net.Conn) error {
	secret, err := dhsecret(conn)
	if err != nil {
		return err
	}
	accountName, session, err := auth(conn, secret)
	if err != nil {
		return err
	}

	agent := agentcenter.findByAccount(accountName)
	if agent != nil {
		sendInnerMsg(agent, "refresh", &RefreshData{conn: conn, session: session, secret: secret})
	} else {
		agent, err := createAgent(conn, accountName, session, secret)
		if err != nil {
			return err
		}
		agentcenter.addByAccount(accountName, agent)
		agent.run()
	}
	return nil
}

// TODO send cache packages to client
func reconnect(conn net.Conn, accountName string) error {
	agent := agentcenter.findByAccount(accountName)
	if agent == nil {
		return errors.New(fmt.Sprintf("agent[%s] not found!", accountName))
	}

	challenge := make([]byte, 8)
	binary.BigEndian.PutUint64(challenge, uint64(rand.Int63()))
	err := writePack(conn, challenge)
	if err != nil {
		return err
	}

	response, err := readPack(conn)
	if err != nil {
		return err
	}

	if !challengeCheck(challenge, agent.secret, response) {
		return errors.New("challenge failed")
	}

	session := rand.Uint32()
	str := strconv.FormatInt(int64(session), 10)
	err = writeEncrypt(conn, []byte(str), agent.secret)
	if err != nil {
		return err
	}

	sendInnerMsg(agent, "refresh", &RefreshData{conn: conn, session: session})
	return nil
}
