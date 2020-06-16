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

	umessage := req.UserIO.player.State()

	server.notifyPlayers(game, nil, MessageNotifyUserState, func (player *Player) interface{} {
		return umessage
	})

	message := NotifyGameState{
		State: game.state,
	}

	if game.state != GameEnded {
		message.CurrentPlayer = game.Players[game.currentPlayer].PlayerID

	} else {
		message.Room = game.solutionRoom
		message.Character = game.solutionCharacter
		message.Weapon = game.solutionWeapon
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
