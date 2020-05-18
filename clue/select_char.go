package clue

import "log"

// HandleSelectCharacterRequest processes select char requests.
func HandleSelectCharacterRequest(server *Server, req Request) {
	selectCharacter, ok := req.Body.(*SelectCharacterRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting SelectCharacterRequest, found", req.Body)
		return
	}

	game, err := server.checkStartedGame(req)

	notify, err := game.SelectCharacter(req.UserIO.player, selectCharacter.Character)
	if err != nil {
		server.sendError(req.UserIO, err.Error())

		return
	}

	if !notify {
		return
	}

	message := NotifyUserState{
		PlayerID:  req.UserIO.player.PlayerID,
		Character: selectCharacter.Character,
	}

	server.notifyPlayers(game, nil, MessageNotifyUserState, func(player *Player) interface{} {
		return message
	})
}
