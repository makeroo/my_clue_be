package web

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/makeroo/my_clue_be/internal/platform/data"
	"github.com/makeroo/my_clue_be/internal/platform/game"
	"github.com/makeroo/my_clue_be/randomstring"
)

// RequestHandler reads the request payload from the ws and then executes it.
type RequestHandler interface {
	RequestType() data.MessageType

	// BodyReader is a callback that reads expected payload from the ws.
	// I didn't succede in creating a "generic" readJSON function so everytime I
	// have to read from a ws I have to know the expected type and parse it.
	BodyReader(*websocket.Conn) (interface{}, error)

	// Handle implements the logic of a specific request.
	Handle(*Server, *Request)
}

type gameUser struct {
	user *User
	io   *UserIO
	// player is io.player but io is defined only if the user is reachable
	player *game.Player
}

type serverGame struct {
	game    *game.Game
	players []*gameUser
}

// Server orchestrates and handles all FE requests.
type Server struct {
	upgrader *websocket.Upgrader
	rand     *rand.Rand

	handlerDescriptors map[data.MessageType]RequestHandler

	// Users that have succesfully signed in.
	signedUsers map[string]*User
	// Users that has connected but not yet signed in.
	connectedUsers []*UserIO

	register   chan *websocket.Conn
	unregister chan *UserIO
	process    chan *Request

	maxMessageSize int64
	pongWait       time.Duration
	pingPeriod     time.Duration
	writeWait      time.Duration

	maxGamesPerPlayer int

	// All the games, starting, running or completed, this server knows of.
	games map[string]*serverGame
}

// New builds a Server instance.
func New(upgrader *websocket.Upgrader, rand *rand.Rand) *Server {
	return &Server{
		upgrader:          upgrader,
		rand:              rand,
		signedUsers:       make(map[string]*User),
		connectedUsers:    nil,
		games:             make(map[string]*serverGame),
		register:          make(chan *websocket.Conn),
		unregister:        make(chan *UserIO),
		process:           make(chan *Request),
		maxMessageSize:    1024,
		pongWait:          60 * time.Second,
		pingPeriod:        55 * time.Second,
		writeWait:         10 * time.Second,
		maxGamesPerPlayer: 10,

		handlerDescriptors: map[data.MessageType]RequestHandler{ /*
				data.MessageVoteStartRequest: {
					BodyReader: func(ws *websocket.Conn) (interface{}, error) {
						body := data.VoteStartRequest{}
						err := ws.ReadJSON(&body)
						return &body, err
					},
					Handler: handlers.HandleVoteStartRequest,
				},
				data.MessageRollDicesRequest: {
					BodyReader: nil,
					Handler:    handlers.HandleRollDicestRequest,
				},
				data.MessageMoveRequest: {
					BodyReader: func(ws *websocket.Conn) (interface{}, error) {
						body := data.MoveRequest{}
						err := ws.ReadJSON(&body)
						return &body, err
					},
					Handler: handlers.HandleMoveRequest,
				},
				data.MessageQuerySolutionRequest: {
					BodyReader: func(ws *websocket.Conn) (interface{}, error) {
						body := data.QuerySolutionRequest{}
						err := ws.ReadJSON(&body)
						return &body, err
					},
					Handler: handlers.HandleQuerySolutionRequest,
				},
				data.MessageRevealRequest: {
					BodyReader: func(ws *websocket.Conn) (interface{}, error) {
						body := data.RevealRequest{}
						err :¨handlers.HandleRevealRequest,
				},
				data.MessagePassRequest: {
					BodyReader: nil,
					Handler:    handlers.HandlePassRequest,
				},
				data.MessageDeclareSolutionRequest: {
					BodyReader: func(ws *websocket.Conn) (interface{}, error) {
						body := data.DeclareSolutionRequest{}
						err := ws.ReadJSON(&body)
						return &body, err
					},
					Handler: handlers.HandleDeclareSolutionRequest,
				},*/
		},
	}
}

func (server *Server) RegisterHandler(handler RequestHandler) {
	server.handlerDescriptors[handler.RequestType()] = handler
}

// Run starts a server.
func (server *Server) Run() {
	go func() {
		for {
			select {
			case ws := <-server.register:
				//log.Println("client to be added")
				server.addClient(ws)
				//log.Println("client added")
			case userIO := <-server.unregister:
				//log.Println("client to be removed", userIO)
				server.removeClient(userIO)
				//log.Println("client removed", userIO)
			case req := <-server.process:
				//log.Println("request to be handled", req, req.UserIO)
				server.handleRequest(req)
				//log.Println("request handled", req)
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
		send: make(chan data.MessageFrame),
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

	log.Println("user disconnected: ", user.token)

	if userIO.player != nil {
		sg := server.games[userIO.player.Game().ID()]

		// find gameUser and reset

		for _, p := range sg.players {
			if p.io == userIO {
				p.io = nil
				break
			}
		}

		userState := data.NotifyUserState{
			ID:        userIO.player.ID(),
			Name:      user.name,
			Character: userIO.player.Character(),
			Online:    false,
		}

		sg.notifyPlayers(userIO.player, data.MessageNotifyUserState, func(target *game.Player) interface{} {
			return userState
		})
	}

	for i, elem := range user.io {
		if userIO == elem {
			user.io[i] = user.io[len(user.io)-1]
			user.io = user.io[:len(user.io)-1]

			return
		}
	}

	log.Println("warning, user not found")
}

func (server *Server) handlerForHeader(msgType data.MessageType) RequestHandler {
	return server.handlerDescriptors[msgType]
}

func (server *Server) handleRequest(req *Request) {
	req.handler.Handle(server, req)
}

// CheckStartedGame performs a few check on incoming request.
func (server *Server) CheckStartedGame(userIO *UserIO) (*game.Game, error) {
	user := userIO.user

	if user == nil {
		return nil, game.NotSignedIn
	}

	if userIO.player == nil {
		return nil, game.NotPlaying
	}

	game := userIO.player.Game()

	return game, nil
}

// CheckCurrentPlayer verifies that request is issued by current player and returns the game.
func (server *Server) CheckCurrentPlayer(req *Request) (*game.Game, error) {
	g, err := server.CheckStartedGame(req.UserIO)

	if err != nil {
		return nil, err
	}

	if g.CurrentPlayer() != req.UserIO.player {
		return nil, game.NotYourTurn
	}

	return g, nil
}

// CheckAnsweringPlayer verifies that request is issued by answering player and returns the game.
func (server *Server) CheckAnsweringPlayer(req *Request) (*game.Game, error) {
	g, err := server.CheckStartedGame(req.UserIO)

	if err != nil {
		return nil, err
	}

	if g.AnsweringPlayer() != req.UserIO.player {
		return nil, game.NotYourTurn
	}

	return g, nil
}

// NotifyPlayers broadcast a message to all the players of a given game.
func (server *Server) NotifyPlayers(g *game.Game, skipPlayer *game.Player, messageType data.MessageType, builder func(*game.Player) interface{}) {
	sg := server.games[g.ID()]

	sg.notifyPlayers(skipPlayer, messageType, builder)
}

func (server *Server) removeConnectedUser(userIO *UserIO) {
	for i, u := range server.connectedUsers {
		if u == userIO {
			server.connectedUsers[i] = server.connectedUsers[len(server.connectedUsers)-1]
			server.connectedUsers = server.connectedUsers[:len(server.connectedUsers)-1]
		}
	}
}

func (g *serverGame) notifyPlayers(skipPlayer *game.Player, message data.MessageType, messageBuilder func(player *game.Player) interface{}) {
	for _, gu := range g.players {
		if gu.player == skipPlayer {
			continue
		}

		if gu.io == nil {
			continue
		}

		//log.Println("notify msg gu", gu.io, message)

		gu.io.send <- data.MessageFrame{
			Header: data.MessageHeader{
				Type: message,
			},
			Body: messageBuilder(gu.player),
		}
	}
}

func (server *Server) randomGameToken() string {
	for {
		t := randomstring.String(server.rand, 4)

		if _, ok := server.games[t]; !ok {
			return t
		}
	}
}

func (server *Server) randomUserToken() string {
	for {
		t := randomstring.String(server.rand, 4)

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
		message := data.MessageHeader{}
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("something went wrong, better to shutdown ws", err)
			break
		}

		requestHandler := server.handlerForHeader(message.Type)

		if requestHandler == nil {
			log.Println("error: unknown request", message.Type)
			continue
		}

		var body interface{}

		body, err = requestHandler.BodyReader(ws)

		if err != nil {
			log.Println("something went wrong, better to shutdown ws", err)
			break
		}

		server.process <- &Request{
			UserIO:  userIO,
			ReqID:   message.ReqID,
			Body:    body,
			handler: requestHandler,
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

			if err := ws.WriteJSON(message.Header); err != nil {
				log.Println("user send failed: user=", user, "error=", err)
				return
			}

			if message.Body != nil {
				if err := ws.WriteJSON(message.Body); err != nil {
					log.Println("user send failed: user=", user, "error=", err)
					return
				}
			}

			if !ok {
				// The hub closed the channel.
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

// NewGame creates a new table.
func (server *Server) NewGame(userIO *UserIO) (*game.Game, *game.Player, error) {
	user := userIO.user

	if user == nil {
		return nil, nil, game.NotSignedIn
	}

	if len(user.joinedGames) >= server.maxGamesPerPlayer {
		return nil, nil, game.TooManyGames
	}

	g := game.New(server.randomGameToken(), server.rand)
	player, err := g.AddPlayer()

	if err != nil {
		return nil, nil, err
	}

	sg := &serverGame{
		game: g,
	}

	gu := &gameUser{
		user:   user,
		io:     userIO,
		player: player,
	}

	sg.players = append(sg.players, gu)

	server.games[g.ID()] = sg

	userIO.player = player
	user.joinedGames = append(user.joinedGames, gu)

	return g, player, nil
}

// RunningGames return a an array of synopses of the non completed games joined by the user.
func (server *Server) RunningGames(user *User) []data.GameSynopsis {
	var runningGames []data.GameSynopsis = nil

	for _, gu := range user.joinedGames {
		runningGames = append(runningGames, server.Synopsis(gu))
	}

	return runningGames
}

// Synopsis returns a game synopsis to fill sign in response.
func (server *Server) Synopsis(targetPlayer *gameUser) data.GameSynopsis {
	var players []data.GamePlayer = nil

	g := targetPlayer.player.Game()
	sg := server.games[g.ID()]

	for _, gu := range sg.players {
		players = append(players, data.GamePlayer{
			Character: gu.player.Character(),
			ID:        gu.player.ID(),
			Name:      gu.user.name,
			Online:    gu.io != nil,
		})
	}

	synopsis := data.GameSynopsis{
		ID:      g.ID(),
		Game:    g.FullState(targetPlayer.player.ID()),
		MyID:    targetPlayer.player.ID(),
		Players: players,
	}

	return synopsis
}

// State return the player state to be notified to the f/e.
func (user *gameUser) State() data.NotifyUserState {
	return data.NotifyUserState{
		ID:        user.player.ID(),
		Name:      user.user.name,
		Character: user.player.Character(),
		Online:    user.io != nil,
	}
}

// JoinGame assign add a player to the given game.
func (server *Server) JoinGame(gameID string, userIO *UserIO) (*data.JoinGameResponse, error) {
	user := userIO.user

	if user == nil {
		return nil, game.NotSignedIn
	}

	if userIO.player != nil {
		// TODO: what about changing game inside a tab?
		// workaround: close and repone ws
		return nil, game.AlreadyPlaying
	}

	sg, ok := server.games[strings.ToUpper(gameID)]

	if !ok {
		return nil, game.UnknownGame
	}

	var rPlayer *game.Player = nil

	for _, gu := range user.joinedGames {
		if gu.player.Game().ID() == gameID {
			if gu.io != nil {
				// no more than 1 tab(ws) per game
				return nil, game.AlreadyPlaying
			}

			// recover an already running game
			// ie. user disconnected for some reason and know she/he has come back!

			gu.io = userIO
			userIO.player = gu.player
			userIO.game = sg

			rPlayer = gu.player

			break
		}
	}

	if rPlayer == nil {
		var err error

		rPlayer, err = sg.game.AddPlayer()

		if err != nil {
			return nil, err
		}

		gu := &gameUser{
			user:   user,
			io:     userIO,
			player: rPlayer,
		}

		sg.players = append(sg.players, gu)

		userIO.player = rPlayer
		userIO.game = sg
		user.joinedGames = append(user.joinedGames, gu)
	}

	players := make([]data.NotifyUserState, len(sg.players))

	for i, gu := range sg.players {
		players[i] = gu.State()
	}

	return &data.JoinGameResponse{
		Players: players,
		MyID:    rPlayer.ID(),
	}, nil
}

// CompleteJoin sends game state to newly joined player and notifies the others.
// FIXME: JoinRequest handler is ugly maybe I should split join in two requests
func (server *Server) CompleteJoin(userIO *UserIO) {
	sg := userIO.game

	if sg.game.Started() {
		userIO.send <- data.MessageFrame{
			Header: data.MessageHeader{
				Type: data.MessageNotifyGameStarted,
			},
			Body: data.NotifyGameStarted{
				PlayersOrder: sg.game.PlayerTurnSequence(),
				Deck:         userIO.player.Deck(),
			},
		}

		sg.game.History(func(record game.MoveRecord) {
			userIO.send <- data.MessageFrame{
				Header: data.MessageHeader{
					Type: data.MessageNotifyMoveRecord,
				},
				Body: record.AsMessageFor(userIO.player),
			}
		})
	}

	message := data.NotifyUserState{
		ID:        userIO.player.ID(),
		Character: userIO.player.Character(),
		Name:      userIO.user.name,
		Online:    true,
	}

	sg.notifyPlayers(userIO.player, data.MessageNotifyUserState, func(player *game.Player) interface{} {
		return message
	})
}

// SelectCharacter assign given character to player provided any other player has not selected it yet.
func (server *Server) SelectCharacter(userIO *UserIO, character game.Card) (*data.NotifyUserState, error) {
	if userIO.player == nil {
		return nil, game.NotPlaying
	}

	if userIO.user == nil {
		return nil, game.NotSignedIn
	}

	game := userIO.player.Game()

	notify, err := game.SelectCharacter(userIO.player, character)

	if err != nil {
		return nil, err
	}

	if !notify {
		return nil, nil
	}

	return &data.NotifyUserState{
		ID:        userIO.player.ID(),
		Character: userIO.player.Character(),
		Name:      userIO.user.name,
		Online:    true,
	}, nil
}

// VoteStart acknowledges player vote and start the game if every player is ready.
func (server *Server) VoteStart(userIO *UserIO, vote bool) (*game.Game, error) {
	g, err := server.CheckStartedGame(userIO)

	if err != nil {
		return nil, err
	}

	started, err := g.VoteStart(userIO.player, vote)

	if err != nil {
		return nil, err
	}

	if !started {
		return nil, nil
	}

	g.Start()

	return g, nil
}
