package clue

import "log"

func HandlePassRequest(server *Server, req Request) {
	_, ok := req.Body.(*PassRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting PassRequest, found", req.Body)
		return
	}

	game, err := server.checkCurrentPlayer(req)

	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	err = game.Pass()

	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	message := NotifyGameState{
		State:         game.state,
		CurrentPlayer: game.currentPlayer,
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
