package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// RevealHandler handles vote start requests.
type RevealHandler struct{}

// RequestType returns Vote Start Request identifier.
func (*RevealHandler) RequestType() data.MessageType {
	return data.MessageRevealRequest
}

// BodyReader parses RevealRequest json from ws.
func (*RevealHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.RevealRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes reveal requests.
func (*RevealHandler) Handle(server *web.Server, req *web.Request) {
	reveal, ok := req.Body.(*data.RevealRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting RevealRequest, found", req.Body)
		return
	}

	g, err := server.CheckAnsweringPlayer(req)

	if err != nil {
		req.SendError(err)

		return
	}

	record, err := g.Reveal(reveal.Card)

	if err != nil {
		req.SendError(err)

		return
	}

	req.SendMessage(data.MessageEmptyResponse, nil)

	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return record.AsMessageFor(player)
	})
}
