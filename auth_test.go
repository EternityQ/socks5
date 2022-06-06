package socks5

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewClientAuthMessage(t *testing.T) {

	t.Run("should generate a message", func(t *testing.T) {
		b := []byte{Version, 2, MethodNoAuth, MethodGssApi}
		reader := bytes.NewReader(b)
		message, err := NewClientAuthMessage(reader)

		if err != nil {
			t.Errorf("wang error == nil bot got %s", err)
		}

		if message.Version != Version {
			t.Errorf("want socksSversion but got %d", message.Version)
		}

		if message.NMethods != 2 {
			t.Errorf("want nmethods = 2 but got %d", message.NMethods)
		}

		if !reflect.DeepEqual(message.Methods, []byte{MethodNoAuth, MethodGssApi}) {
			t.Errorf("want methods: %v, but got %v", []byte{MethodNoAuth, MethodGssApi}, message.Methods)
		}

	})

	t.Run("methods length is shorter than nMethods", func(t *testing.T) {
		b := []byte{Version, 2, MethodNoAuth}
		reader := bytes.NewReader(b)
		message, err := NewClientAuthMessage(reader)

		if err == nil || message != nil {
			t.Errorf("should get error != nil but got nile")
		}
	})

}

func TestNewServerAuthMessage(t *testing.T) {
	t.Run("should pass", func(t *testing.T) {
		var buf bytes.Buffer
		err := NewServerAuthMessage(&buf, MethodNoAuth)

		if err != nil {
			t.Fatalf("should get nil error but got %s", err)
		}

		got := buf.Bytes()

		if !reflect.DeepEqual(got, []byte{Version, MethodNoAuth}) {
			t.Fatalf("should send %v but send %v", []byte{Version, MethodNoAuth}, got)
		}
	})
}

func TestNewClientPasswordMessage(t *testing.T) {
	t.Run("valid password auth message", func(t *testing.T) {
		username, password := "admin", "123456"
		buf := bytes.Buffer{}
		buf.Write([]byte{PasswordMethodVersion, 5})
		buf.WriteString(username)
		buf.WriteByte(6)
		buf.WriteString(password)
		message, err := NewClientPasswordMessage(&buf)

		if err != nil {
			t.Errorf("want error = nil but got %s", err)
		}

		want := ClientPasswordMessage{Username: username, Password: password}
		if *message != want {
			t.Errorf("want message %#v but got %#v", want, *message)
		}

	})
}
