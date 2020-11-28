package web

import (
	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
)

// UserIO collects data to handle ws I/O.
// There is an instance per websocket/browser tab.
// A user can have multiple tab/windows running different games.
// Each tab is binded to one game at most though.
// A use is not allowed to have more than one tab/window opened on the same game.
type UserIO struct {
	ws   *websocket.Conn
	send chan data.MessageFrame

	// user is defined after a sign in request
	user *User
	// player and game are defined after a create or join game request
	player *game.Player
	game   *serverGame
}

// User collects all the info to recognize a user and to allow her/him to play Clue.
type User struct {
	// Name is visible to all users.
	name string
	// Token is the secret that used to recognize a user.
	token string

	// io is a collection of all opened websockets of a user.
	io []*UserIO

	joinedGames []*gameUser
}
