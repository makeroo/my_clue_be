package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// SelectCharHandler handles select character requests.
type SelectCharHandler struct{}

// RequestType returns Select Character Request identifier.
func (*SelectCharHandler) RequestType() data.MessageType {
	return data.MessageSelectCharRequest
}

// BodyReader parses SelectCharacterRequest json from ws.
func (*SelectCharHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.SelectCharacterRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes select char requests.
func (*SelectCharHandler) Handle(server *web.Server, req *web.Request) {
	selectCharacter, ok := req.Body.(*data.SelectCharacterRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting SelectCharacterRequest, found", req.Body)
		return
	}

	g, err := server.CheckStartedGame(req.UserIO)

	if err != nil {
		req.SendError(err)

		return
	}

	newUserState, err := server.SelectCharacter(req.UserIO, selectCharacter.Character)

	if err != nil {
		req.SendError(err)

		return
	}

	if newUserState == nil {
		return
	}

	server.NotifyPlayers(g, nil, data.MessageNotifyUserState, func(player *game.Player) interface{} {
		return newUserState
	})

	req.SendMessage(data.MessageEmptyResponse, nil)
}
