package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// DeclareSolutionHandler handles sign in requests.
type DeclareSolutionHandler struct{}

// RequestType returns Declare Solution Request identifier.
func (*DeclareSolutionHandler) RequestType() data.MessageType {
	return data.MessageDeclareSolutionRequest
}

// BodyReader parses DeclareSolutionRequest json from ws.
func (*DeclareSolutionHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.DeclareSolutionRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes declare solution requests.
func (*DeclareSolutionHandler) Handle(server *web.Server, req *web.Request) {
	declareSolution, ok := req.Body.(*data.DeclareSolutionRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting DeclareSolutionRequest, found", req.Body)
		return
	}

	g, err := server.CheckCurrentPlayer(req)

	if err != nil {
		req.SendError(err)

		return
	}

	record, err := g.CheckSolution(declareSolution.Character, declareSolution.Room, declareSolution.Weapon)

	if err != nil {
		req.SendError(err)

		return
	}
	/*
		// first of all, req user state has changed, whether she/he found the solution or not

		umessage := req.UserIO.player.State()

		server.notifyPlayers(game, nil, data.MessageNotifyUserState, func(player *Player) interface{} {
			return umessage
		})

		// then, if the req user failed, the game could be ended, if she/he was the last but one to fail

		probablyWinner := game.Players[game.currentPlayer]

		if game.state == GameEnded && probablyWinner.PlayerID != req.UserIO.player.PlayerID {
			umessage = probablyWinner.State()

			server.notifyPlayers(game, nil, MessageNotifyUserState, func(player *Player) interface{} {
				return umessage
			})
		}

		// finally, notify updated game state

		message := NotifyGameState{
			State:         game.state,
			CurrentPlayer: game.Players[game.currentPlayer].PlayerID,
		}

		if game.state == GameEnded {
			message.Room = game.solutionRoom
			message.Character = game.solutionCharacter
			message.Weapon = game.solutionWeapon
		}
	*/
	server.NotifyPlayers(g, nil, data.MessageNotifyMoveRecord, func(player *game.Player) interface{} {
		return record.AsMessageFor(player)
	})
}
