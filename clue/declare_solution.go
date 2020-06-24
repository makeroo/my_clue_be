package clue

import "log"

// HandleDeclareSolutionRequest processes declare solution requests.
func HandleDeclareSolutionRequest(server *Server, req *Request) {
	declareSolution, ok := req.Body.(*DeclareSolutionRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting DeclareSolutionRequest, found", req.Body)
		return
	}

	game, err := server.checkCurrentPlayer(req)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	err = game.CheckSolution(declareSolution.Character, declareSolution.Room, declareSolution.Weapon)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	// first of all, req user state has changed, whether she/he found the solution or not

	umessage := req.UserIO.player.State()

	server.notifyPlayers(game, nil, MessageNotifyUserState, func(player *Player) interface{} {
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

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
