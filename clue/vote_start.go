package clue

import "log"

// HandleVoteStartRequest processes vote start requests.
func HandleVoteStartRequest(server *Server, req *Request) {
	voteStart, ok := req.Body.(*VoteStartRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting VoteStartRequest, found", req.Body)
		return
	}

	game, err := server.checkStartedGame(req)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	start, err := game.VoteStart(req.UserIO.player, voteStart.Vote)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	if !start {
		return
	}

	game.Start()

	server.notifyPlayers(game, nil, MessageNotifyGameStarted, func(player *Player) interface{} {
		return game.GameStartedMessage(player)
	})

	newTurn := game.FullState(req.UserIO.player.PlayerID)

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return newTurn
	})
}
