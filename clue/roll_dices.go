package clue

import "log"

func HandleRollDicestRequest(server *Server, req Request) {
	_, ok := req.Body.(*RollDicesRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting RollDicesRequest, found", req.Body)
		return
	}

	game, err := server.checkCurrentPlayer(req)

	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	if err := game.RollDices(); err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	message := NotifyGameState{
		State:         game.state,
		CurrentPlayer: game.Players[game.currentPlayer].PlayerID,
		Dice1:         game.dice1,
		Dice2:         game.dice2,
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
