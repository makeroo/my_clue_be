package clue

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/utils"
)

// Server orchestrates and handles all FE requests.
type Server struct {
	upgrader websocket.Upgrader

	signedUsers map[string]*User
	users       []*User

	register   chan *websocket.Conn
	unregister chan *UserIO
	process    chan request

	maxMessageSize int64
	pongWait       time.Duration
	pingPeriod     time.Duration
	writeWait      time.Duration

	maxGamesPerPlayer int

	games map[string]*Game
}

type UserIO struct {
	ws   *websocket.Conn
	send chan Message

	user   *User
	player *Player
}

// User is actually just a websocket connection.
// Eventually it will hold a user name
type User struct {
	server *Server

	Name  string
	Token string

	io []*UserIO

	joinedGames []*Player
}

type request struct {
	userIO  *UserIO
	message *Message
}

// NewServer builds a Server instance.
func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		signedUsers:       make(map[string]*User),
		users:             nil,
		games:             make(map[string]*Game),
		register:          make(chan *websocket.Conn),
		unregister:        make(chan *UserIO),
		process:           make(chan request),
		maxMessageSize:    1024,
		pongWait:          60 * time.Second,
		pingPeriod:        55 * time.Second,
		writeWait:         10 * time.Second,
		maxGamesPerPlayer: 10,
	}
}

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
	user := &User{
		server: server,
	}

	server.users = append(server.users, user)

	user.addIO(conn)
}

func (server *Server) removeClient(userIO *UserIO) {
	user := userIO.user

	for i, elem := range user.io {
		if userIO == elem {
			user.io[i] = user.io[len(user.io)-1]
			user.io = user.io[:len(user.io)-1]

			if len(user.io) == 0 {
				fmt.Println("user unreachable: ", user.Token)

				server.broadcast(user, func(me *Player, target *Player) Message {
					return Message{
						NotifyUserState: &NotifyUserState{
							OldName:   user.Name,
							NewName:   user.Name,
							Character: me.Character,
							Online:    false,
						},
					}
				})
			}

			return
		}
	}

	log.Println("warning, user not found")
}

func (server *Server) handleRequest(req request) {
	if signIn := req.message.SignIn; signIn != nil {
		user := req.userIO.user
		var oldName string

		if signIn.Token == "" {
			// this is a new user: generate a new token and return it

			if user.Token != "" {
				// can't happen
				req.userIO.send <- Message{
					Error: AlreadySignedIn,
				}

				return
			}

			oldName = ""
			user.Token = server.randomUserToken()
			user.Name = signIn.Name

			server.signedUsers[user.Token] = user

			fmt.Println("new user: name=", user.Name, "token=", user.Token)

			req.userIO.send <- Message{
				SignInResponse: &SignInResponse{
					Token: user.Token,
				},
			}

		} else if user.Token != "" && user.Token != signIn.Token {
			req.userIO.send <- Message{
				Error: TokenMismatch,
			}

			return

		} else if len(user.io) != 1 {
			fmt.Println("error: join request from a user with more than a socket")

			return

		} else {
			user = server.updateUser(req.userIO, signIn.Token)

			if user == nil {
				req.userIO.send <- Message{
					Error: UnknownToken,
				}

				return
			}

			oldName = user.Name
			user.Name = signIn.Name

			var runningGames []GameSynopsis = nil

			for _, player := range user.joinedGames {
				runningGames = append(runningGames, GameSynopsis{
					GameID:    player.Game.GameID,
					Character: player.Character,
				})
			}

			req.userIO.send <- Message{
				SignInResponse: &SignInResponse{
					RunningGames: runningGames,
				},
			}

			fmt.Println("user back online: token=", signIn.Token)
		}

		user.Name = signIn.Name

		server.broadcast(user, func(me *Player, target *Player) Message {
			return Message{
				NotifyUserState: &NotifyUserState{
					OldName:   oldName,
					NewName:   user.Name,
					Character: me.Character,
					Online:    true,
				},
			}
		})

	} else if req.message.CreateGame != nil {
		user := req.userIO.user

		if user.Token == "" {
			req.userIO.send <- Message{
				Error: NotSignedIn,
			}

			return
		}

		if len(user.joinedGames) >= server.maxGamesPerPlayer {
			req.userIO.send <- Message{
				Error: TooManyGames,
			}

			return
		}

		g := NewGame(server.randomGameToken())

		server.games[g.GameID] = g

		player, err := g.AddPlayer(req.userIO)

		if err != nil {
			req.userIO.send <- Message{
				Error: err.Error(),
			}

			return
		}

		user.joinedGames = append(user.joinedGames, player)

		req.userIO.send <- Message{
			CreateGameResponse: &CreateGameResponse{GameID: g.GameID},
		}
	}
}

/*func (server *Server) findUser(token string) (*User, bool) {
	for _, value := range server.users {
		if value.Token == token {
			return value, true
		}
	}

	return nil, false
}*/

func (server *Server) updateUser(userIO *UserIO, token string) *User {
	oldUser, found := server.signedUsers[token]

	if !found {
		return nil
	}

	for i, u := range server.users {
		if u == userIO.user {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
		}
	}

	userIO.user = oldUser
	oldUser.io = append(oldUser.io, userIO)

	return oldUser
}

func (server *Server) broadcast(user *User, messageBuilder func(me *Player, target *Player) Message) {
	for _, player := range user.joinedGames {
		game := player.Game

		for _, target := range game.Players {
			//if target == player {
			//	continue
			//}

			for _, io := range target.UserIO.user.io {
				io.send <- messageBuilder(player, target)
			}
		}
	}
}

func (server *Server) randomGameToken() string {
	for {
		t := utils.String(4)

		if server.games[t] == nil {
			return t
		}
	}
}

func (server *Server) randomUserToken() string {
	for {
		t := utils.String(4)

		if server.signedUsers[t] == nil {
			return t
		}
	}
}

func (user *User) addIO(conn *websocket.Conn) {
	newIO := &UserIO{
		ws:   conn,
		user: user,
		send: make(chan Message),
	}

	user.io = append(user.io, newIO)

	go newIO.writePump()
	go newIO.readPump()
}

func (userIO *UserIO) readPump() {
	ws := userIO.ws

	defer func() {
		userIO.user.server.unregister <- userIO
		ws.Close()
	}()

	ws.SetReadLimit(userIO.user.server.maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(userIO.user.server.pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(userIO.user.server.pongWait)); return nil })

	for {
		message := Message{}
		err := userIO.ws.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			// TODO: handle json decoding erros
			break
		}

		userIO.user.server.process <- request{userIO, &message}
		log.Println("delivered", message)
	}
}

func (userIO *UserIO) writePump() {
	ticker := time.NewTicker(userIO.user.server.pingPeriod)

	ws := userIO.ws
	user := userIO.user
	server := user.server

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
				ws.WriteJSON(message) // TODO: handle error
				return
			}

			if err := ws.WriteJSON(message); err != nil {
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
