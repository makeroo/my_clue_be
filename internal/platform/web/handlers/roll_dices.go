package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// RollDicesHandler handles roll dices requests.
type RollDicesHandler struct{}

// RequestType returns Roll Dices Request identifier.
func (*RollDicesHandler) RequestType() data.MessageType {
	return data.MessageRollDicesRequest
}

// BodyReader does nothing: roll dices request does not a payload.
func (*RollDicesHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	return nil, nil
}

// Handle processes roll dices requests.
func (*RollDicesHandler) Handle(server *web.Server, req *web.Request) {
	g, err := server.CheckCurrentPlayer(req)

	if err != nil {
		req.SendError(err)

		return
	}

	record, err := g.RollDices()

	if err != nil {
		req.SendError(err)

		return
	}

	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return record.AsMessageFor(player)
	})
}
