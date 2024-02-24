package main

import "github.com/sanity32/go-sft-imgcap/internal/server"

func main() {
	addr := ":8111"

	serv := server.New()

	serv.Listen(addr)
}
