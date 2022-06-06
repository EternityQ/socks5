package socks5

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	ErrVersionNotSupported       = errors.New("protocol version not supported")
	ErrMethodVersionNotSupported = errors.New("sub-negotiation method version not supported")
	ErrCommandNotSupported       = errors.New("request command not supported")
	ErrInvalidReservedField      = errors.New("invalid reserved field")
	ErrAddressTypeNotSupported   = errors.New("address type not supported")
)

const (
	ReservedField = 0x00
)

type Server interface {
	Run() error
}

type SOCKES5Server struct {
	IP   string
	Port int
	Conf *Config
}

type Config struct {
	AuthMethod    Method
	PasswordCheck func(username, password string) bool
}

func initConfig(config *Config) error {
	if config.AuthMethod == MethodPassword && config.PasswordCheck == nil {
		return ErrPasswordCheckerNotSet
	}
	return nil
}

func (s SOCKES5Server) Run() error {

	log.Println(`
			   _____ ____  ________ _______ ______
			  / ___// __ \/ ____/ //_/ ___// ____/
			  \__ \/ / / / /   / ,<  \__ \/___ \  
			 ___/ / /_/ / /___/ /| |___/ /___/ /  
			/____/\____/\____/_/ |_/____/_____/         `)

	if err := initConfig(s.Conf); err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	listen, err := net.Listen("tcp", address)
	log.Println("Socks5 service starts successfully and listens to port", s.Port)

	if err != nil {
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("connection failure form %s:%s", conn.RemoteAddr(), err)
			continue
		}
		log.Println("Client connection: ", conn.RemoteAddr())
		go func() {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {

				}
			}(conn)
			err := handleConnection(conn, s.Conf)
			if err != nil {
				log.Printf("connection failure form %s:%s", conn.RemoteAddr(), err)
			}
		}()
	}
}

func handleConnection(conn net.Conn, config *Config) error {

	if err := authentication(conn, config); err != nil {
		return err
	}
	targetConn, err := request(conn)
	if err != nil {
		return err
	}

	return forward(conn, targetConn)
}

func request(conn io.ReadWriter) (io.ReadWriteCloser, error) {

	message, err := NewClientRequestMessage(conn)
	if err != nil {
		return nil, err
	}

	targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", message.Address, message.Port))
	if err != nil {
		return nil, WriteRequestFailureMessage(conn, ReplyConnectionRefused)
	}
	addrValue := targetConn.LocalAddr()
	addr := addrValue.(*net.TCPAddr)
	return targetConn, WriteRequestSuccessMessage(conn, addr.IP, uint16(addr.Port))
}

func authentication(conn io.ReadWriter, config *Config) error {

	clientMessage, err := NewClientAuthMessage(conn)
	if err != nil {
		return err
	}

	var selectMethod = MethodNoAcceptable
	for _, method := range clientMessage.Methods {
		if method == config.AuthMethod {
			selectMethod = method
			break
		}
	}

	switch selectMethod {
	case MethodNoAuth:
		if err := NewServerAuthMessage(conn, MethodNoAuth); err != nil {
			return err
		}
	case MethodPassword:
		if err := NewServerAuthMessage(conn, MethodPassword); err != nil {
			return err
		}
		clientPasswordMessage, err := NewClientPasswordMessage(conn)
		if err != nil {
			return err
		}
		if !config.PasswordCheck(clientPasswordMessage.Username, clientPasswordMessage.Password) {
			_ = WriteServerPasswordMessage(conn, PasswordAUthFailure)
			return ErrPasswordAuthFailure
		}
		if err := WriteServerPasswordMessage(conn, PasswordAuthSuccess); err != nil {
			return err
		}

	case MethodNoAcceptable:
		if err := NewServerAuthMessage(conn, MethodNoAcceptable); err != nil {
			return err
		} else {
			return errors.New("method not support")
		}
	}

	return nil
}

func forward(conn io.ReadWriter, targetConn io.ReadWriteCloser) error {
	defer func(targetConn io.ReadWriteCloser) {
		err := targetConn.Close()
		if err != nil {

		}
	}(targetConn)
	go func() {
		_, err := io.Copy(targetConn, conn)
		if err != nil {

		}
	}()
	_, err := io.Copy(conn, targetConn)
	return err
}
