package clue

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/utils"
)

// MessageReader is a callback that reads expected payload from the ws.
// I didn't succede in creating a "generic" readJSON function so everytime I
// have to read from a ws I have to create a MessageReader func.
type MessageReader func(*websocket.Conn) (interface{}, error)

// RequestHandlerDescriptor is the descriptor of a request handler.
type RequestHandlerDescriptor struct {
	BodyReader MessageReader
	Handler    RequestHandler
}

// Server orchestrates and handles all FE requests.
type Server struct {
	upgrader *websocket.Upgrader
	rand     *rand.Rand

	handlerDescriptors map[string]RequestHandlerDescriptor

	// Users that have succesfully signed in.
	signedUsers map[string]*User
	// Users that has connected but not yet signed in.
	connectedUsers []*UserIO

	register   chan *websocket.Conn
	unregister chan *UserIO
	process    chan Request

	maxMessageSize int64
	pongWait       time.Duration
	pingPeriod     time.Duration
	writeWait      time.Duration

	maxGamesPerPlayer int

	// All the games, starting, running or completed, this server knows of.
	games map[string]*Game
}

/* UserIO collects data to handle ws I/O.
 * There is an instance per websocket/browser tab.
 * A user can have multiple tab/windows running different games.
 * Each tab is binded to one game at most though.
 * A use is not allowed to have more than one tab/window opened on the same game.
 */
type UserIO struct {
	ws   *websocket.Conn
	send chan MessageFrame

	// user is defined after a sign in request
	user *User
	// player  is defined after a create or join game request
	player *Player
}

// User collects all the info to recognize a user and to allow her/him to play Clue.
type User struct {
	// Name is visible to all users.
	Name string
	// Token is the secret that used to recognize a user.
	Token string

	// io is a collection of all opened websockets of a user.
	io []*UserIO

	joinedGames []*Player
}

// RequestHandler implements the logic of a specific request.
type RequestHandler func(*Server, Request)

// Request is an incoming request to be served.
type Request struct {
	// UserIO is the user who issued the request.
	UserIO  *UserIO
	Body    interface{}
	Handler RequestHandler
}

// NewServer builds a Server instance.
func NewServer(upgrader *websocket.Upgrader, rand *rand.Rand) *Server {
	return &Server{
		upgrader:          upgrader,
		rand:              rand,
		signedUsers:       make(map[string]*User),
		connectedUsers:    nil,
		games:             make(map[string]*Game),
		register:          make(chan *websocket.Conn),
		unregister:        make(chan *UserIO),
		process:           make(chan Request),
		maxMessageSize:    1024,
		pongWait:          60 * time.Second,
		pingPeriod:        55 * time.Second,
		writeWait:         10 * time.Second,
		maxGamesPerPlayer: 10,

		handlerDescriptors: map[string]RequestHandlerDescriptor{
			MessageSignInRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := SignInRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleSignInRequest,
			},
			MessageCreateGameRequest: {
				BodyReader: nil,
				Handler:    HandleSignInRequest,
			},
			MessageJoinGameRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := JoinGameRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleJoinGameRequest,
			},
			MessageSelectCharRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := SelectCharacterRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleSelectCharacterRequest,
			},
			MessageVoteStartRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := VoteStartRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleVoteStartRequest,
			},
			MessageRollDicesRequest: {
				BodyReader: nil,
				Handler:    HandleRollDicestRequest,
			},
			MessageMoveRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := MoveRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleMoveRequest,
			},
			MessageQuerySolutionRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := QuerySolutionRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleQuerySolutionRequest,
			},
			MessageRevealRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := RevealRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleRevealRequest,
			},
			MessagePassRequest: {
				BodyReader: nil,
				Handler:    HandlePassRequest,
			},
			MessageDeclareSolutionRequest: {
				BodyReader: func(ws *websocket.Conn) (interface{}, error) {
					body := DeclareSolutionRequest{}
					err := ws.ReadJSON(&body)
					return &body, err
				},
				Handler: HandleDeclareSolutionRequest,
			},
		},
	}
}

// Run starts a server.
func (server *Server) Run() {
	go func() {
		for {
			select {
			case ws := <-server.register:
				server.addClient(ws)
			case userIO := <-server.unregister:
				server.removeClient(userIO)
			case req := <-server.process:
				server.handleRequest(req)
			}
		}
	}()
}

// Handle receives an HTTP request and upgrade to websocket protocol.
// It is the ws entry point.
func (server *Server) Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := server.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("ws upgrade failed: request=", r, "error=", err)

		// TODO: return error

		return
	}

	server.register <- ws
}

func (server *Server) addClient(conn *websocket.Conn) {
	userIO := &UserIO{
		ws:   conn,
		send: make(chan MessageFrame),
	}

	go userIO.writePump(server)
	go userIO.readPump(server)

	server.connectedUsers = append(server.connectedUsers, userIO)
}

func (server *Server) removeClient(userIO *UserIO) {
	user := userIO.user

	if user == nil {
		// disconnected before signin in
		server.removeConnectedUser(userIO)

		return
	}

	if userIO.player != nil {
		userIO.player.UserIO = nil
	}

	for i, elem := range user.io {
		if userIO == elem {
			user.io[i] = user.io[len(user.io)-1]
			user.io = user.io[:len(user.io)-1]

			if len(user.io) == 0 {
				fmt.Println("user unreachable: ", user.Token)

				server.broadcast(user, MessageNotifyUserState, func(me *Player, target *Player) interface{} {
					return NotifyUserState{
						PlayerID: me.PlayerID,
						Online:   false,
					}
				})
			}

			return
		}
	}

	log.Println("warning, user not found")
}

func (server *Server) handlerForHeader(msgType string) (RequestHandlerDescriptor, bool) {
	requestHandler, ok := server.handlerDescriptors[msgType]
	return requestHandler, ok
}

func (server *Server) handleRequest(req Request) {
	req.Handler(server, req)
}

func (server *Server) sendError(userIO *UserIO, err string) {
	userIO.send <- MessageFrame{
		Header: MessageHeader{
			Type: MessageError,
		},
		Body: NotifyError{
			Error: err,
		},
	}
}

func (server *Server) checkStartedGame(req Request) (*Game, error) {
	user := req.UserIO.user

	if user == nil {
		return nil, errors.New(NotSignedIn)
	}

	if req.UserIO.player == nil {
		return nil, errors.New(NotPlaying)
	}

	game := req.UserIO.player.Game

	return game, nil
}

func (server *Server) checkCurrentPlayer(req Request) (*Game, error) {
	game, err := server.checkStartedGame(req)

	if err != nil {
		return nil, err
	}

	if !game.IsCurrentPlayer(req.UserIO.player) {
		return nil, errors.New(NotYourTurn)
	}

	return game, nil
}

func (server *Server) removeConnectedUser(userIO *UserIO) {
	for i, u := range server.connectedUsers {
		if u == userIO {
			server.connectedUsers[i] = server.connectedUsers[len(server.connectedUsers)-1]
			server.connectedUsers = server.connectedUsers[:len(server.connectedUsers)-1]
		}
	}
}

func (server *Server) broadcast(user *User, message string, messageBuilder func(me *Player, target *Player) interface{}) {
	if user.joinedGames == nil {
		return
	}

	for _, player := range user.joinedGames {
		game := player.Game

		for _, target := range game.Players {
			//if target == player {
			//	continue
			//}

			if target.UserIO == nil {
				continue
			}

			for _, io := range target.UserIO.user.io {
				io.send <- MessageFrame{
					Header: MessageHeader{
						Type: message,
					},
					Body: messageBuilder(player, target),
				}
			}
		}
	}
}

func (server *Server) notifyPlayers(game *Game, skipPlayer *Player, message string, messageBuilder func(player *Player) interface{}) {
	for _, player := range game.Players {
		if player == skipPlayer {
			continue
		}

		player.UserIO.send <- MessageFrame{
			Header: MessageHeader{
				Type: message,
			},
			Body: messageBuilder(player),
		}
	}
}

func (server *Server) randomGameToken() string {
	for {
		t := utils.String(server.rand, 4)

		if server.games[t] == nil {
			return t
		}
	}
}

func (server *Server) randomUserToken() string {
	for {
		t := utils.String(server.rand, 4)

		if server.signedUsers[t] == nil {
			return t
		}
	}
}

func (userIO *UserIO) readPump(server *Server) {
	ws := userIO.ws

	defer func() {
		server.unregister <- userIO
		ws.Close()
	}()

	ws.SetReadLimit(server.maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(server.pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(server.pongWait)); return nil })

	for {
		message := MessageHeader{}
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("something went wrong, better to shutdown ws", err)
			break
		}

		requestHandler, ok := server.handlerForHeader(message.Type)

		if !ok {
			log.Println("error: unknown request", message.Type)
			continue
		}

		var body interface{}
		if requestHandler.BodyReader != nil {
			body, err = requestHandler.BodyReader(ws)

			if err != nil {
				log.Println("something went wrong, better to shutdown ws", err)
				break
			}

		} else {
			body = nil
		}

		server.process <- Request{
			UserIO:  userIO,
			Body:    body,
			Handler: requestHandler.Handler,
		}

		log.Println("delivered", message)
	}
}

func (userIO *UserIO) writePump(server *Server) {
	ticker := time.NewTicker(server.pingPeriod)

	ws := userIO.ws
	user := userIO.user

	defer func() {
		ticker.Stop()
		ws.Close()
	}()

	for {
		select {
		case message, ok := <-userIO.send:
			ws.SetWriteDeadline(time.Now().Add(server.writeWait))
			if !ok {
				// The hub closed the channel.
				ws.WriteJSON(message.Header) // TODO: handle error
				ws.WriteJSON(message.Body)   // TODO: handle error
				return
			}

			if err := ws.WriteJSON(message.Header); err != nil {
				log.Println("user send failed: user=", user, "error=", err)
				return
			}

			if err := ws.WriteJSON(message.Body); err != nil {
				log.Println("user send failed: user=", user, "error=", err)
				return
			}

		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(server.writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
