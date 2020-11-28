package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// PassHandler handles pass requests.
type PassHandler struct{}

// RequestType returns Pass Request identifier.
func (*PassHandler) RequestType() data.MessageType {
	return data.MessagePassRequest
}

// BodyReader does nothing, pass request doesn't have a payload.
func (*PassHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	return nil, nil
}

// Handle processes pass requests.
func (*PassHandler) Handle(server *web.Server, req *web.Request) {
	g, err := server.CheckCurrentPlayer(req)

	if err != nil {
		req.SendError(err)

		return
	}

	record, err := g.Pass()

	if err != nil {
		req.SendError(err)

		return
	}

	req.SendMessage(data.MessageEmptyResponse, nil)

	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return record.AsMessageFor(player)
	})
}
