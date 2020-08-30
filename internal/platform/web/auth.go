package web

import (
	"fmt"

	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
)

// SignIn registers a new user.
func (server *Server) SignIn(userIO *UserIO, name string) (*User, string) {
	user := userIO.user

	if user != nil {
		if name != "" {
			user.name = name
		}

		return user, user.token
	}

	user = &User{
		token: server.randomUserToken(),
		name:  name,
	}

	userIO.user = user
	user.io = append(user.io, userIO)

	server.signedUsers[user.token] = user

	server.removeConnectedUser(userIO)

	fmt.Println("new user: name=", user.name)

	return user, user.token
}

// Authenticate checks provided token against known users.
func (server *Server) Authenticate(userIO *UserIO, name string, token string) (*User, error) {
	user := userIO.user

	if user != nil {
		// signin request from an already signed in user

		if user.token != token {
			return nil, game.TokenMismatch
		}

		if name != "" {
			user.name = name
		}

		// a new tab from a known user?

		for _, io := range user.io {
			if io == userIO {
				return user, nil
			}
		}

	} else {
		// signin request from a disconnected user?

		user = server.signedUsers[token]

		if user == nil {
			return nil, game.UnknownToken
		}

		if name != "" {
			user.name = name
		}

		// fmt.Println("user back online: token=", token)
	}

	server.removeConnectedUser(userIO)

	userIO.user = user
	user.io = append(user.io, userIO)

	server.broadcast(user, data.MessageNotifyUserState, func(me *gameUser, target *gameUser) interface{} {
		return me.State()
	})

	return user, nil
}
