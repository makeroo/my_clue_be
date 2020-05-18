package clue

import "log"

// HandleQuerySolutionRequest processes query solution requests.
func HandleQuerySolutionRequest(server *Server, req Request) {
	querySolution, ok := req.Body.(*QuerySolutionRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting QuerySolutionRequest, found", req.Body)
		return
	}

	game, err := server.checkCurrentPlayer(req)

	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	if err := game.QuerySolution(querySolution.Character, querySolution.Weapon); err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	message := NotifyGameState{
		State:           game.state,
		AnsweringPlayer: game.Players[game.queryingPlayer].PlayerID,
		Character:       game.queryCharacter,
		Room:            game.queryRoom,
		Weapon:          game.queryWeapon,
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
