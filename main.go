package main

import (
	"github.com/krissukoco/go-gin-chat/server"
)

func main() {
	srv, err := server.NewDefaultServer()
	if err != nil {
		panic(err)
	}
	srv.Start()
}
