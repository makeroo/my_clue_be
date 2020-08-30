package clue

import (
	"fmt"
	"log"
)

// HandleSignInRequest processes sign in requests.
func HandleSignInRequest(server *Server, req *Request) {
	signIn, ok := req.Body.(*SignInRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting SignInRequest, found", req.Body)
		return
	}

	user := req.UserIO.user

	if signIn.Token == "" {
		// this is a new user: generate a new token and return it

		if user != nil {
			req.UserIO.send <- MessageFrame{
				Header: MessageHeader{
					Type:  MessageSignInResponse,
					ReqID: req.ReqID,
				},
				Body: SignInResponse{
					RunningGames: user.RunningGames(),
				},
			}

			return
		}

		user = &User{
			Token: server.randomUserToken(),
			Name:  signIn.Name,
		}

		req.UserIO.user = user
		user.io = append(user.io, req.UserIO)

		server.signedUsers[user.Token] = user

		fmt.Println("new user: name=", user.Name, "token=", user.Token)

		req.UserIO.send <- MessageFrame{
			Header: MessageHeader{
				Type:  MessageSignInResponse,
				ReqID: req.ReqID,
			},
			Body: SignInResponse{
				Token: user.Token,
			},
		}

		server.removeConnectedUser(req.UserIO)

	} else if user != nil {
		// signin request from an already signed in user

		if user.Token != signIn.Token {
			server.sendError(req, TokenMismatch)

			return
		}

		// a new tab from a known user?

		for _, io := range user.io {
			if io == req.UserIO {
				req.UserIO.send <- MessageFrame{
					Header: MessageHeader{
						Type:  MessageSignInResponse,
						ReqID: req.ReqID,
					},
					Body: SignInResponse{
						RunningGames: user.RunningGames(),
					},
				}

				return
			}
		}

		if signIn.Name != "" {
			user.Name = signIn.Name
		}

		req.UserIO.user = user
		user.io = append(user.io, req.UserIO)

		server.removeConnectedUser(req.UserIO)

	} else {
		// signin request from a disconnected user?

		user = server.signedUsers[signIn.Token]

		if user == nil {
			server.sendError(req, UnknownToken)

			return
		}

		server.removeConnectedUser(req.UserIO)

		req.UserIO.user = user
		user.io = append(user.io, req.UserIO)

		user.Name = signIn.Name

		req.UserIO.send <- MessageFrame{
			Header: MessageHeader{
				Type:  MessageSignInResponse,
				ReqID: req.ReqID,
			},
			Body: SignInResponse{
				RunningGames: req.UserIO.user.RunningGames(),
			},
		}

		fmt.Println("user back online: token=", signIn.Token)
	}

	server.broadcast(user, MessageNotifyUserState, func(me *Player, target *Player) interface{} {
		return me.State()
	})
}
