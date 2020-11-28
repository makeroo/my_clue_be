package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/web"
	"github.com/makeroo/my_clue_be/internal/platform/web/handlers"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "http service address")

	flag.Parse()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO: implement check origin
			// see https://github.com/gorilla/websocket/issues/367
			origin := r.Header.Values("Origin")
			log.Printf("origin to check: %s", origin)
			return true
		},
	}

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	server := web.New(&upgrader, seededRand)

	server.RegisterHandler(&handlers.SignInHandler{})
	server.RegisterHandler(&handlers.CreateGameHandler{})
	server.RegisterHandler(&handlers.JoinGameHandler{})
	server.RegisterHandler(&handlers.SelectCharHandler{})
	server.RegisterHandler(&handlers.VoteStartHandler{})
	server.RegisterHandler(&handlers.RollDicesHandler{})
	server.RegisterHandler(&handlers.MoveHandler{})
	server.RegisterHandler(&handlers.PassHandler{})
	server.RegisterHandler(&handlers.QuerySolutionHandler{})
	server.RegisterHandler(&handlers.RevealHandler{})
	server.RegisterHandler(&handlers.DeclareSolutionHandler{})

	server.Run()

	http.Handle("/clue/ws", logRequest(server))

	log.Println("My Cluedo B/E up and running")
	// TODO: log version
	// TODO: log config

	// TODO: handle sigterm, graceful shutdown

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func logRequest(server *web.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		server.Handle(w, r)
	})
}
