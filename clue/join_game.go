package clue

import "log"

// HandleJoinGameRequest processes join game requests.
func HandleJoinGameRequest(server *Server, req *Request) {
	joinGame, ok := req.Body.(*JoinGameRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting JoinGameRequest, found", req.Body)
		return
	}

	user := req.UserIO.user

	if user == nil {
		server.sendError(req, NotSignedIn)

		return
	}

	if req.UserIO.player != nil {
		// TODO: what about changing game inside a tab?
		// workaround: close and repone ws
		server.sendError(req, AlreadyPlaying)

		return
	}

	game := server.games[joinGame.GameID]

	if game == nil {
		server.sendError(req, UnknownGame)

		return
	}

	var rPlayer *Player = nil

	for _, player := range user.joinedGames {
		if player.Game.GameID == game.GameID {
			if player.UserIO != nil {
				// no more than 1 tab(ws) per game
				server.sendError(req, AlreadyPlaying)

				return
			}

			// recover an already running game
			// ie. user disconnected for some reason and know she/he has come back!

			player.UserIO = req.UserIO
			req.UserIO.player = player

			rPlayer = player

			break
		}
	}

	if rPlayer == nil {
		var err error

		rPlayer, err = game.AddPlayer(req.UserIO)

		if err != nil {
			server.sendError(req, err.Error())

			return
		}

		req.UserIO.player = rPlayer
		user.joinedGames = append(user.joinedGames, rPlayer)
	}

	players := make([]NotifyUserState, len(game.Players))

	for i, player := range game.Players {
		userState := NotifyUserState{
			PlayerID:  player.PlayerID,
			Character: player.Character,
			Online:    player.UserIO != nil,
			Name:      player.User.Name,
		}

		players[i] = userState
	}

	req.UserIO.send <- MessageFrame{
		Header: MessageHeader{
			Type:  MessageJoinGameResponse,
			ReqID: req.ReqID,
		},
		Body: JoinGameResponse{
			Players:  players,
			PlayerID: rPlayer.PlayerID,
		},
	}

	message := NotifyUserState{
		PlayerID:  req.UserIO.player.PlayerID,
		Name:      user.Name,
		Character: req.UserIO.player.Character,
		Online:    true,
	}

	server.notifyPlayers(game, req.UserIO.player, MessageNotifyUserState, func(player *Player) interface{} {
		return message
	})
}
