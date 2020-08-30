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
	return data.MessageVoteStartRequest
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
	/*
		message := NotifyGameState{
			State: game.state,
		}

		var skipPlayer *Player = nil
		currentPlayer := game.Players[game.currentPlayer]

		if game.state == GameStateTrySolution {
			message.Revealed = game.Revealed

			skipPlayer = currentPlayer

			if currentPlayer.UserIO != nil {
				messageWithRevealedCard := message

				messageWithRevealedCard.RevealedCard = game.RevealedCard

				currentPlayer.UserIO.send <- MessageFrame{
					Header: MessageHeader{
						Type: MessageNotifyGameState,
					},
					Body: messageWithRevealedCard,
				}
			}

		} else {
			message.AnsweringPlayer = game.Players[game.answeringPlayer].PlayerID
			message.Character = game.queryCharacter
			message.Room = game.queryRoom
			message.Weapon = game.queryWeapon
		}
	*/
	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return record.AsMessageFor(player)
	})
}
