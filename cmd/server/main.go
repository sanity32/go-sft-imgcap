package main

import (
	"fmt"

	"github.com/sanity32/go-sft-imgcap/internal/model"
	"github.com/sanity32/go-sft-imgcap/internal/server"
)

func main() {
	addr := ":18084"
	// syscall.Umask(0)
	fmt.Println("umasc called")
	serv := server.New()

	model.MainHashDir.Create()

	serv.Listen(addr)
}
