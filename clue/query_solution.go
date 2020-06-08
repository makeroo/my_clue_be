package clue

import "log"

// HandleQuerySolutionRequest processes query solution requests.
func HandleQuerySolutionRequest(server *Server, req *Request) {
	querySolution, ok := req.Body.(*QuerySolutionRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting QuerySolutionRequest, found", req.Body)
		return
	}

	game, err := server.checkCurrentPlayer(req)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	message := NotifyGameState{
		State:     game.state,
		Character: game.queryCharacter,
		Room:      game.queryRoom,
		Weapon:    game.queryWeapon,
	}

	player, err := game.QuerySolution(querySolution.Character, querySolution.Weapon)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	message.AnsweringPlayer = game.Players[game.queryingPlayer].PlayerID

	if player != nil {
		message.PlayerPositions = append(message.PlayerPositions, PlayerPosition{
			PlayerID: player.PlayerID,
			Room:     player.Room,
			// map x/y are always 0
		})
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
