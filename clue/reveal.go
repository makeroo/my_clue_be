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

	if req.UserIO.player.PlayerID != game.queryingPlayer {
		server.sendError(req, NotYourTurn)

		return
	}

	matched, err := game.Reveal(reveal.Card)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	message := NotifyGameState{
		State: game.state,
	}

	if game.state == GameStateTrySolution {
		message.Matched = matched

	} else {
		message.AnsweringPlayer = game.Players[game.queryingPlayer].PlayerID
		message.Character = game.queryCharacter
		message.Room = game.queryRoom
		message.Weapon = game.queryWeapon
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
