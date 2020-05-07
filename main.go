package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/makeroo/my_clue_be/clue"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "http service address")

	flag.Parse()

	server := clue.NewServer()

	server.Run()

	http.HandleFunc("/ws", server.Handle)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
