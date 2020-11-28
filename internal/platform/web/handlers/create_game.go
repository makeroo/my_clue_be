package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// CreateGameHandler handles create game requests.
type CreateGameHandler struct{}

// RequestType returns Create Game Request identfier.
func (*CreateGameHandler) RequestType() data.MessageType {
	return data.MessageCreateGameRequest
}

// BodyReader does nothing, create game request doesn't have a payload.
func (*CreateGameHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	return nil, nil
}

// Handle processes create game requests.
func (*CreateGameHandler) Handle(server *web.Server, req *web.Request) {
	g, player, err := server.NewGame(req.UserIO)

	if err != nil {
		req.SendError(err)
	}

	req.SendMessage(data.MessageCreateGameResponse, data.CreateGameResponse{
		GameID: g.ID(),
		MyID:   player.ID(),
	})
}
