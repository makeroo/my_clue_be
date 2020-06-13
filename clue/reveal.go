package clue

import "log"

// HandleRevealRequest processes reveal requests.
func HandleRevealRequest(server *Server, req *Request) {
	reveal, ok := req.Body.(*RevealRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting RevealRequest, found", req.Body)
		return
	}

	game, err := server.checkStartedGame(req)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	if req.UserIO.player.PlayerID != game.Players[game.answeringPlayer].PlayerID {
		server.sendError(req, NotYourTurn)

		return
	}

	err = game.Reveal(reveal.Card)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

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

	server.notifyPlayers(game, skipPlayer, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
