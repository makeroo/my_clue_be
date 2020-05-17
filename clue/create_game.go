package clue

import "log"

func HandleCreateGameRequest(server *Server, req Request) {
	_, ok := req.Body.(*CreateGameRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting CreateGameRequest, found", req.Body)
		return
	}

	user := req.UserIO.user

	if user == nil {
		server.sendError(req.UserIO, NotSignedIn)

		return
	}

	if len(user.joinedGames) >= server.maxGamesPerPlayer {
		server.sendError(req.UserIO, TooManyGames)

		return
	}

	g := NewGame(server.randomGameToken(), server.rand)

	server.games[g.GameID] = g

	player, err := g.AddPlayer(req.UserIO)

	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	req.UserIO.player = player
	user.joinedGames = append(user.joinedGames, player)

	req.UserIO.send <- MessageFrame{
		Header: MessageHeader{
			Type: MessageCreateGameResponse,
		},
		Body: CreateGameResponse{GameID: g.GameID},
	}
}
