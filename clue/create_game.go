package clue

// HandleCreateGameRequest processes create game requests.
func HandleCreateGameRequest(server *Server, req *Request) {
	user := req.UserIO.user

	if user == nil {
		server.sendError(req, NotSignedIn)

		return
	}

	if len(user.joinedGames) >= server.maxGamesPerPlayer {
		server.sendError(req, TooManyGames)

		return
	}

	g := NewGame(server.randomGameToken(), server.rand)

	server.games[g.GameID] = g

	player, err := g.AddPlayer(req.UserIO)

	if err != nil {
		server.sendError(req, err.Error())

		return
	}

	req.UserIO.player = player
	user.joinedGames = append(user.joinedGames, player)

	req.UserIO.send <- MessageFrame{
		Header: MessageHeader{
			Type:  MessageCreateGameResponse,
			ReqID: req.ReqID,
		},
		Body: CreateGameResponse{GameID: g.GameID},
	}
}
