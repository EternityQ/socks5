package socks5

import (
	"errors"
	"io"
)

type ClientAuthMessage struct {
	Version  byte
	NMethods byte
	Methods  []Method
}

type ClientPasswordMessage struct {
	Username string
	Password string
}

type Method = byte

const Version = 0x05

var (
	ErrPasswordCheckerNotSet = errors.New("error password checker not set")
	ErrPasswordAuthFailure   = errors.New("error authenticating username or password")
)

const (
	PasswordMethodVersion = 0x01
	PasswordAuthSuccess   = 0x00
	PasswordAUthFailure   = 0xFF
)

const (
	MethodNoAuth       Method = 0x00
	MethodGssApi       Method = 0x01
	MethodPassword     Method = 0x02
	MethodNoAcceptable Method = 0xFF
)

// NewClientAuthMessage ...
func NewClientAuthMessage(conn io.Reader) (*ClientAuthMessage, error) {

	//Red version nMethods
	buf := make([]byte, 2)

	_, err := io.ReadFull(conn, buf)

	if err != nil {
		return nil, err
	}

	// Validate version
	if buf[0] != Version {
		return nil, ErrVersionNotSupported
	}

	//Read methods
	nMethods := buf[1]
	buf = make([]byte, nMethods)
	_, err = io.ReadFull(conn, buf)

	if err != nil {
		return nil, err
	}

	return &ClientAuthMessage{
		Version:  Version,
		NMethods: nMethods,
		Methods:  buf,
	}, nil

}

func NewServerAuthMessage(conn io.Writer, method Method) error {
	buf := []byte{Version, method}
	_, err := conn.Write(buf)
	return err
}

func NewClientPasswordMessage(conn io.Reader) (*ClientPasswordMessage, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	version, usernameLen := buf[0], buf[1]

	if version != PasswordMethodVersion {
		return nil, ErrMethodVersionNotSupported
	}

	buf = make([]byte, usernameLen+1)

	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	username, passwordLen := string(buf[:usernameLen]), buf[len(buf)-1]

	if len(buf) < int(passwordLen) {
		buf = make([]byte, passwordLen)
	}

	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}

	password := string(buf)

	return &ClientPasswordMessage{
		Username: username,
		Password: password,
	}, nil
}

func WriteServerPasswordMessage(conn io.Writer, status byte) error {
	_, err := conn.Write([]byte{PasswordMethodVersion, status})
	return err
}
