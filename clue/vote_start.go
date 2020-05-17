package clue

import "log"

func HandleVoteStartRequest(server *Server, req Request) {
	voteStart, ok := req.Body.(*VoteStartRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting VoteStartRequest, found", req.Body)
		return
	}

	game, err := server.checkStartedGame(req)

	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	start, err := game.VoteStart(req.UserIO.player, voteStart.Vote)

	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	if !start {
		return
	}

	game.Start()

	playersOrder := []int{}

	for _, player := range game.Players {
		playersOrder = append(playersOrder, player.PlayerID)
	}

	server.notifyPlayers(game, nil, MessageNotifyGameStarted, func(player *Player) interface{} {
		return NotifyGameStarted{
			Deck:         player.Deck,
			PlayersOrder: playersOrder,
		}
	})

	newTurn := NotifyGameState{
		State:         game.state,
		CurrentPlayer: game.Players[game.currentPlayer].PlayerID,
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return newTurn
	})
}
