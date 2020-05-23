package clue

import "log"

// HandleMoveRequest processes move requests.
func HandleMoveRequest(server *Server, req *Request) {
	move, ok := req.Body.(*MoveRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting MoveRequest, found", req.Body)
		return
	}

	game, err := server.checkStartedGame(req)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	if err := game.Move(move.EnterRoom, move.MapX, move.MapY); err != nil {
		server.sendError(req, err.Error())

		return
	}

	message := NotifyGameState{
		State:         game.state,
		CurrentPlayer: game.Players[game.currentPlayer].PlayerID,

		Room: req.UserIO.player.Room,
		MapX: req.UserIO.player.MapX,
		MapY: req.UserIO.player.MapY,
	}

	server.notifyPlayers(game, nil, MessageNotifyGameState, func(player *Player) interface{} {
		return message
	})
}
