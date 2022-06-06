package main

import (
	"log"
	"socks5"
)

func main() {

	users := map[string]string{
		"zhang": "111111",
	}

	socks5Server := socks5.SOCKES5Server{
		IP:   "127.0.0.1",
		Port: 1080,
		Conf: &socks5.Config{
			AuthMethod: socks5.MethodPassword,
			PasswordCheck: func(username, password string) bool {
				wantPassword, ok := users[username]
				if !ok {
					return false
				}
				return wantPassword == password
			},
		},
	}

	err := socks5Server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
