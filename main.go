package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/clue"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "http service address")

	flag.Parse()

	upgrader := websocket.Upgrader{
		// TODO: debug code, remove checkorigin
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	server := clue.NewServer(&upgrader, seededRand)

	server.Run()

	http.HandleFunc("/ws", server.Handle)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
