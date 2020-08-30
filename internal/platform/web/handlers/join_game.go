package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// JoinGameHandler handles join game requests.
type JoinGameHandler struct{}

// RequestType returns Join Game Request identifier.
func (*JoinGameHandler) RequestType() data.MessageType {
	return data.MessageJoinGameRequest
}

// BodyReader parses JoinGameRequest json from ws.
func (*JoinGameHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.JoinGameRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes join game requests.
func (*JoinGameHandler) Handle(server *web.Server, req *web.Request) {
	joinGame, ok := req.Body.(*data.JoinGameRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting JoinGameRequest, found", req.Body)
		return
	}

	resp, err := server.JoinGame(joinGame.GameID, req.UserIO)

	if err != nil {
		req.SendError(err)

		return
	}

	req.SendMessage(data.MessageJoinGameResponse, resp)

	server.CompleteJoin(req.UserIO)
}
