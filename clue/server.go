package clue

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/utils"
)

// Server orchestrates and handles all FE requests.
type Server struct {
	upgrader *websocket.Upgrader
	rand     *rand.Rand

	// Users that have succesfully signed in.
	signedUsers map[string]*User
	// Users that has connected but not yet signed in.
	connectedUsers []*UserIO

	register   chan *websocket.Conn
	unregister chan *UserIO
	process    chan request

	maxMessageSize int64
	pongWait       time.Duration
	pingPeriod     time.Duration
	writeWait      time.Duration

	maxGamesPerPlayer int

	// All the games, starting, running or completed, this server knows of.
	games map[string]*Game
}

/* UserIO handles ws requests.
 * There is an instance per websocket/browser tab.
 * A user can have multiple tab/windows running different games.
 * Each tab is binded to one game at most though.
 * A use is not allowed to have more than one tab/window opened on the same game.
 */
type UserIO struct {
	ws   *websocket.Conn
	send chan Message

	user   *User
	player *Player
}

// User collects all the info to recognize a user and to allow her/him to play Clue.
type User struct {
	//server *Server

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
func NewServer(upgrader *websocket.Upgrader, rand *rand.Rand) *Server {
	return &Server{
		upgrader:          upgrader,
		rand:              rand,
		signedUsers:       make(map[string]*User),
		connectedUsers:    nil,
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
		send: make(chan Message),
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

				server.broadcast(user, func(me *Player, target *Player) Message {
					return Message{
						NotifyUserState: &NotifyUserState{
							PlayerID: me.PlayerID,
							Online:   false,
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

		if signIn.Token == "" {
			// this is a new user: generate a new token and return it

			if user != nil {
				req.userIO.send <- Message{
					Error: AlreadySignedIn,
				}

				return
			}

			user = &User{
				Token: server.randomUserToken(),
				Name:  signIn.Name,
			}

			req.userIO.user = user
			user.io = append(user.io, req.userIO)

			server.signedUsers[user.Token] = user

			fmt.Println("new user: name=", user.Name, "token=", user.Token)

			req.userIO.send <- Message{
				SignInResponse: &SignInResponse{
					Token: user.Token,
				},
			}

			server.removeConnectedUser(req.userIO)

		} else if user != nil {
			// signin request from an already signed in user

			if user.Token != signIn.Token {
				req.userIO.send <- Message{
					Error: TokenMismatch,
				}

				return
			}

			// a new tab from a known user?

			for _, io := range user.io {
				if io == req.userIO {
					req.userIO.send <- Message{
						Error: AlreadySignedIn, // TODO: remove error to reuse signin to change name
					}

					return
				}
			}

			if req.message.SignIn.Name != "" {
				user.Name = req.message.SignIn.Name
			}

			req.userIO.user = user
			user.io = append(user.io, req.userIO)

			server.removeConnectedUser(req.userIO)

		} else {
			// signin request from a disconnected user?

			user = server.signedUsers[signIn.Token]

			if user == nil {
				req.userIO.send <- Message{
					Error: UnknownToken,
				}

				return
			}

			server.removeConnectedUser(req.userIO)

			req.userIO.user = user
			user.io = append(user.io, req.userIO)

			user.Name = signIn.Name

			var runningGames []GameSynopsis = nil

			for _, player := range user.joinedGames {
				runningGames = append(runningGames, GameSynopsis{
					GameID:    player.Game.GameID,
					Character: player.Character,
					PlayerID:  player.PlayerID,
				})
			}

			req.userIO.send <- Message{
				SignInResponse: &SignInResponse{
					RunningGames: runningGames,
				},
			}

			fmt.Println("user back online: token=", signIn.Token)
		}

		server.broadcast(user, func(me *Player, target *Player) Message {
			return Message{
				NotifyUserState: &NotifyUserState{
					PlayerID:  me.PlayerID,
					Name:      user.Name,
					Character: me.Character,
					Online:    true,
				},
			}
		})

	} else if req.message.CreateGame != nil {
		user := req.userIO.user

		if user == nil {
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

		g := NewGame(server.randomGameToken(), server.rand)

		server.games[g.GameID] = g

		player, err := g.AddPlayer(req.userIO)

		if err != nil {
			req.userIO.send <- Message{
				Error: err.Error(),
			}

			return
		}

		req.userIO.player = player
		user.joinedGames = append(user.joinedGames, player)

		req.userIO.send <- Message{
			CreateGameResponse: &CreateGameResponse{GameID: g.GameID},
		}

	} else if req.message.JoinGame != nil {
		user := req.userIO.user

		if user == nil {
			req.userIO.send <- Message{
				Error: NotSignedIn,
			}

			return
		}

		if req.userIO.player != nil {
			req.userIO.send <- Message{
				Error: AlreadyPlaying, // TODO: what about changing game inside a tab?
			}
		}

		game := server.games[req.message.JoinGame.GameID]

		if game == nil {
			req.userIO.send <- Message{
				Error: UnknownGame,
			}

			return
		}

		found := false

		for _, player := range user.joinedGames {
			if player.Game.GameID == game.GameID {
				if player.UserIO != nil {
					req.userIO.send <- Message{
						Error: AlreadyPlaying,
					}

					return
				}

				// recover an already running game
				// ie. user disconnected for some reason and know she/he has come back!

				player.UserIO = req.userIO
				req.userIO.player = player

				found = true

				break
			}
		}

		if !found {
			player, err := game.AddPlayer(req.userIO)

			if err != nil {
				req.userIO.send <- Message{
					Error: err.Error(),
				}

				return
			}

			req.userIO.player = player
			user.joinedGames = append(user.joinedGames, player)
		}

		message := Message{
			NotifyUserState: &NotifyUserState{
				PlayerID:  req.userIO.player.PlayerID,
				Name:      user.Name,
				Character: req.userIO.player.Character,
				Online:    true,
			},
		}

		server.notifyPlayers(game, nil, func(player *Player) Message {
			return message
		})

	} else if req.message.SelectCharacter != nil {
		game, err := server.checkStartedGame(req)

		notify, err := game.SelectCharacter(req.userIO.player, req.message.SelectCharacter.Character)
		if err != nil {
			req.userIO.send <- Message{
				Error: err.Error(),
			}

			return
		}

		if !notify {
			return
		}

		message := Message{
			NotifyUserState: &NotifyUserState{
				PlayerID:  req.userIO.player.PlayerID,
				Character: req.message.SelectCharacter.Character,
			},
		}

		server.notifyPlayers(game, nil, func(player *Player) Message {
			return message
		})

	} else if req.message.VoteStart != nil {
		game, err := server.checkStartedGame(req)

		if err != nil {
			req.userIO.send <- Message{
				Error: err.Error(),
			}

			return
		}

		start, err := game.VoteStart(req.userIO.player, req.message.VoteStart.Vote)

		if err != nil {
			req.userIO.send <- Message{
				Error: err.Error(),
			}

			return
		}

		if !start {
			return
		}

		game.Start()

		playersOrder := []int{}

		for _, player := range game.Players {
			playersOrder = append(playersOrder, player.PlayerID)
		}

		server.notifyPlayers(game, nil, func(player *Player) Message {
			return Message{
				NotifyGameStarted: &NotifyGameStarted{
					Deck:         player.Deck,
					PlayersOrder: playersOrder,
				},
			}
		})

		newTurn := Message{
			NotifyGameState: &NotifyGameState{
				State:         game.state,
				CurrentPlayer: game.currentPlayer,
			},
		}

		server.notifyPlayers(game, nil, func(player *Player) Message {
			return newTurn
		})

	} else if req.message.RollDices != nil {
	}
}

func (server *Server) checkStartedGame(req request) (*Game, error) {
	user := req.userIO.user

	if user == nil {
		return nil, errors.New(NotSignedIn)
	}

	if req.userIO.player == nil {
		return nil, errors.New(NotPlaying)
	}

	game := req.userIO.player.Game

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

func (server *Server) broadcast(user *User, messageBuilder func(me *Player, target *Player) Message) {
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
				io.send <- messageBuilder(player, target)
			}
		}
	}
}

func (server *Server) notifyPlayers(game *Game, skipPlayer *Player, messageBuilder func(player *Player) Message) {
	for _, player := range game.Players {
		if player == skipPlayer {
			continue
		}

		player.UserIO.send <- messageBuilder(player)
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
		message := Message{}
		err := userIO.ws.ReadJSON(&message)
		if err != nil {
			if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
				log.Println("debug: unmarshal error, ignoring request", jsonErr)
				continue
			}

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			// TODO: handle json decoding erros
			break
		}

		server.process <- request{userIO, &message}
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
