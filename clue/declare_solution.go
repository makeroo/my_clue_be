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

	message := NotifyGameState{
		State: game.state,
	}

	if game.state != GameEnded {
		message.CurrentPlayer = game.Players[game.currentPlayer].PlayerID
		message.Character = declareSolution.Character
		message.Room = declareSolution.Room
		message.Weapon = declareSolution.Weapon
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
