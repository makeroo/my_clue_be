package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/web"
)

// SignInHandler handles sign in requests.
type SignInHandler struct{}

// RequestType returns Sign In Request identifier.
func (*SignInHandler) RequestType() data.MessageType {
	return data.MessageSignInRequest
}

// BodyReader parses SignInRequest json from ws.
func (*SignInHandler) BodyReader(ws *websocket.Conn) (interface{}, error) {
	body := data.SignInRequest{}
	err := ws.ReadJSON(&body)
	return &body, err
}

// Handle processes sign in requests.
func (*SignInHandler) Handle(server *web.Server, req *web.Request) {
	signIn, ok := req.Body.(*data.SignInRequest)

	if !ok {
		log.Println("ERROR request type mismatch, expecting SignInRequest, found", req.Body)
		return
	}

	if signIn.Token == "" {
		// this is a new user: generate a new token and return it

		user, token := server.SignIn(req.UserIO, signIn.Name)

		req.SendMessage(data.MessageSignInResponse, data.SignInResponse{
			Token:        token,
			RunningGames: server.RunningGames(user),
		})

		return

	}

	user, err := server.Authenticate(req.UserIO, signIn.Name, signIn.Token)

	if err != nil {
		req.SendError(err)

		return
	}

	req.SendMessage(data.MessageSignInResponse, data.SignInResponse{
		RunningGames: server.RunningGames(user),
	})
}
