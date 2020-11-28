package game

import (
	"math/rand"
	"time"
	//"github.com/makeroo/my_clue_be/internal/platform/web"
)

// State is the state of the Game FSM
type State int

const (
	// GameStateStarting is the state a Game begins in: players are joining
	// the table
	GameStateStarting State = iota
	// GameStateNewTurn is the state when a player starts her/his turn.
	// Dices have not rolled yet.
	GameStateNewTurn
	// GameStateCard is the state entered if dices contain a '1' (lens)
	// and card have to be drawed from the "hints" deck
	// TODO: actually this is a composite state depending on card instructions
	GameStateCard
	// GameStateMove is the state entered when dice have rolled, eventually
	// the hints card have been obeyed and the player have to choose what to do.
	GameStateMove
	// GameStateQuery is the state entered when the player completed her/his
	// move and, if in a room, query a solution.
	GameStateQuery
	// GameStateTrySolution is the state entered when the player can declare
	// the solution or pass.
	GameStateTrySolution
	// GameEnded is the state entered when someone find out what the solution is.
	GameEnded
)

// Game is a clue table. A user can join multiple tables.
type Game struct {
	gameID  string
	players []*Player
	rand    *rand.Rand

	solution Declaration

	state           State
	currentPlayer   int
	dice1           int
	dice2           int
	remainingSteps  int
	query           Declaration
	answeringPlayer int

	revealed     bool
	revealedCard Card

	secretPassages [][2]Card

	history []*MoveRecord
}

// New create a Game instance.
func New(gameID string, rand *rand.Rand) *Game {
	game := Game{
		gameID: gameID,
		rand:   rand,

		state: GameStateStarting,

		secretPassages: [][2]Card{
			{Kitchen, Study},
			{Study, Kitchen},
			{Lounge, Conservatory},
			{Conservatory, Lounge},
		},
	}

	return &game
}

// ID returns the game id.
func (game *Game) ID() string {
	return game.gameID
}

// Started return true if the game has started.
func (game *Game) Started() bool {
	return game.state != GameStateStarting
}

// AddPlayer adds a player to the table.
func (game *Game) AddPlayer( /*userIO *web.UserIO*/ ) (*Player, error) {
	if game.state != GameStateStarting {
		return nil, CannotJoinRunningGame
	}

	if len(game.players) == 6 {
		return nil, TableIsFull
	}

	player := &Player{
		game: game,
		//UserIO: userIO,
		//User:   userIO.user,
		id: PlayerID(len(game.players) + 1),
	}

	game.players = append(game.players, player)

	return player, nil
}

// Start starts a new game.
func (game *Game) Start() error {
	if game.state != GameStateStarting {
		return GameAlreadyStarted
	}

	game.shufflePlayers()

	game.state = GameStateNewTurn
	game.currentPlayer = 0

	// create secret

	game.solution = Declaration{
		Character: game.randomCard(MissScarlett, MrsWhite),
		Room:      game.randomCard(Kitchen, Study),
		Weapon:    game.randomCard(Candlestick, Wrenck),
	}

	deck := game.makeDeckWithoutSolution()

	game.rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	cardsPerPlayer := len(deck) / len(game.players)
	playersWithAnExtraCard := len(deck) % len(game.players)
	start := 0

	for i, player := range game.players {
		cards := cardsPerPlayer

		if i < playersWithAnExtraCard {
			cards++
		}

		player.deck = deck[start : start+cards]
		start += cards
	}

	for _, player := range game.players {
		player.position = initialPositions[player.character]
	}

	return nil
}

func (game *Game) shufflePlayers() {
	game.rand.Shuffle(len(game.players), func(i, j int) {
		game.players[i], game.players[j] = game.players[j], game.players[i]
	})
}

func (game *Game) randomCard(min Card, max Card) Card {
	n := game.rand.Intn(int(max)-int(min)+1) + int(min)

	return Card(n)
}

func (game *Game) makeDeckWithoutSolution() []Card {
	deck := make([]Card, Cards-3)

	p := 0

	for c := Candlestick; c <= MrsWhite; c++ {
		if c == game.solution.Character || c == game.solution.Room || c == game.solution.Weapon {
			continue
		}

		deck[p] = c
		p++
	}

	return deck
}

// SelectCharacter assigns a character to a player.
func (game *Game) SelectCharacter(player *Player, character Card) (bool, error) {
	if game.state != GameStateStarting {
		return false, GameAlreadyStarted
	}

	if !IsCharacter(character) {
		return false, NotACharacter
	}

	if player.character == character {
		return false, nil
	}

	for _, cplayer := range game.players {
		if cplayer.character == character {
			return false, AlreadySelected
		}
	}

	player.character = character

	return true, nil
}

// VoteStart records the start vote of a player.
func (game *Game) VoteStart(player *Player, vote bool) (bool, error) {
	if game.state != GameStateStarting {
		return false, GameAlreadyStarted
	}

	if !IsCharacter(Card(player.character)) {
		return false, CharacterNotSelected
	}

	if player.votedStart == vote {
		return false, nil
	}

	player.votedStart = vote

	if !vote || len(game.players) < 2 {
		return false, nil
	}

	for _, p := range game.players {
		if !p.votedStart {
			return false, nil
		}
	}

	return true, nil
}

// CurrentPlayer returns current player if the game has started.
func (game *Game) CurrentPlayer() *Player {
	if game.state == GameStateStarting {
		return nil
	}

	return game.players[game.currentPlayer]
}

// AnsweringPlayer returns the anwering player if the game state is .
func (game *Game) AnsweringPlayer() *Player {
	if game.state != GameStateQuery {
		return nil
	}

	if game.answeringPlayer == -1 {
		return nil
	}

	return game.players[game.answeringPlayer]
}

// RollDices rolls dices for current player.
func (game *Game) RollDices() (*MoveRecord, error) {
	if game.state != GameStateNewTurn {
		return nil, IllegalState
	}

	game.dice1 = game.rand.Intn(6) + 1
	game.dice2 = game.rand.Intn(6) + 1

	game.remainingSteps = game.dice1 + game.dice2

	// TODO: cards
	//if game.dice1 == 1 || game.dice2 == 1 {
	//	game.state = GameStateCard
	//} else {
	game.state = GameStateMove
	//}

	record := &MoveRecord{
		PlayerID:  game.players[game.currentPlayer].id,
		Timestamp: time.Now(),
		Move: &RollDicesMove{
			Dice1: game.dice1,
			Dice2: game.dice2,
		},
		StateDelta: StateUpdate{
			State:          game.state,
			Dice1:          game.dice1,
			Dice2:          game.dice2,
			RemainingSteps: game.remainingSteps,
		},
	}

	game.history = append(game.history, record)

	return record, nil
}

// Move moves current player.
func (game *Game) Move(room Card, mapX int, mapY int) (*MoveRecord, error) {
	if game.state != GameStateMove {
		return nil, IllegalState
	}

	player := game.players[game.currentPlayer]
	var playerPosition PlayerPosition
	var move2 Move

	if IsRoom(room) {
		if room == player.position.Room {
			// the player choose to remain in the same room she/he was in
			game.state = GameStateQuery
			game.answeringPlayer = -1

			record := &MoveRecord{
				PlayerID:  player.id,
				Timestamp: time.Now(),
				Move: &EnterRoomMove{
					Room: player.position.Room,
				},
				StateDelta: StateUpdate{
					State: game.state,
					//AnsweringPlayer: PlayerID(game.answeringPlayer),
				},
			}

			game.history = append(game.history, record)

			return record, nil
		}

		if player.position.InRoom() {
			// changing room is possible only through secret passages

			if !game.IsSecretPassage(player.position.Room, room) {
				return nil, IllegalMove
			}

			player.position.EnterRoom(room)

			game.state = GameStateQuery
			game.answeringPlayer = -1

			record := &MoveRecord{
				PlayerID:  player.id,
				Timestamp: time.Now(),
				Move: &EnterRoomMove{
					Room: player.position.Room,
				},
				StateDelta: StateUpdate{
					State: game.state,
					//AnsweringPlayer: PlayerID(game.answeringPlayer),
					Positions: []PlayerPosition{
						{
							PlayerID:     player.id,
							PawnPosition: player.position,
						},
					},
				},
			}

			game.history = append(game.history, record)

			return record, nil
		}

		// the player is in the hallway, check if she/he is in front of a door of
		// the room she/he wants to enter in

		cellType := Card(clueBoard[player.position.MapY][player.position.MapX])

		if cellType != room {
			return nil, IllegalMove
		}

		player.position.EnterRoom(room)

		playerPosition = PlayerPosition{
			PlayerID:     player.id,
			PawnPosition: player.position,
		}

		move2 = &EnterRoomMove{
			Room: room,
		}

		game.state = GameStateQuery
		game.answeringPlayer = -1

	} else if player.position.InRoom() {
		if !game.IsValidPosition(mapX, mapY) {
			return nil, IllegalMove
		}

		if p := game.IsOccupied(mapX, mapY); p != nil && p != player {
			return nil, IllegalMove
		}

		// the player just exited a room
		// check the hallway pos she/he selected is one in front of a door of
		// the room she/he was in
		cellType := Card(clueBoard[mapY][mapX])

		if cellType != player.position.Room {
			return nil, IllegalMove
		}

		player.position.MoveTo(mapX, mapY)

		playerPosition = PlayerPosition{
			PlayerID:     player.id,
			PawnPosition: player.position,
		}

		move2 = &MovingInTheHallwayMove{
			MapX: mapX,
			MapY: mapY,
		}

		game.remainingSteps--

		// if the player has just exited a room then this is the first step
		// no need to check remainingSteps because the min is 2 so it is at least 1

	} else {
		if !game.IsValidPosition(mapX, mapY) {
			return nil, IllegalMove
		}

		if clueBoard[mapY][mapX] < 0 {
			return nil, IllegalMove
		}

		if game.IsOccupied(mapX, mapY) != nil {
			return nil, IllegalMove
		}

		if !player.position.IsAdjacent(mapX, mapY) {
			return nil, IllegalMove
		}

		player.position.MoveTo(mapX, mapY)

		playerPosition = PlayerPosition{
			PlayerID:     player.id,
			PawnPosition: player.position,
		}

		move2 = &MovingInTheHallwayMove{
			MapX: mapX,
			MapY: mapY,
		}

		game.remainingSteps--

		if game.remainingSteps == 0 {
			game.state = GameStateTrySolution
			game.answeringPlayer = -1
		}
	}

	var answeringPlayer PlayerID
	if game.answeringPlayer == -1 {
		answeringPlayer = 0
	} else {
		answeringPlayer = game.players[game.answeringPlayer].id
	}

	record := &MoveRecord{
		PlayerID:  player.id,
		Timestamp: time.Now(),
		Move:      move2,
		StateDelta: StateUpdate{
			State:           game.state,
			AnsweringPlayer: answeringPlayer,
			RemainingSteps:  game.remainingSteps,
			Positions: []PlayerPosition{
				playerPosition,
			},
		},
	}

	game.history = append(game.history, record)

	return record, nil
}

// IsValidPosition checks coordinate ranges.
func (game *Game) IsValidPosition(mapX, mapY int) bool {
	return mapX >= 0 && mapX <= 23 && mapY >= 0 && mapY <= 24
}

// IsOccupied checks if position is occupied by a player.
func (game *Game) IsOccupied(mapX, mapY int) *Player {
	for _, p := range game.players {
		if p.position.MapX == mapX && p.position.MapY == mapY {
			return p
		}
	}

	return nil
}

// IsSecretPassage checks if there is a secret passage.
func (game *Game) IsSecretPassage(from, to Card) bool {
	for _, secretPassage := range game.secretPassages {
		if from == secretPassage[0] && to == secretPassage[1] {
			return true
		}
	}

	return false
}

// QuerySolution starts a query solution process.
func (game *Game) QuerySolution(character, weapon Card) (*MoveRecord, error) {
	if game.state != GameStateQuery || game.answeringPlayer != -1 {
		return nil, IllegalState
	}

	currentPlayer := game.players[game.currentPlayer]

	room := currentPlayer.position.Room
	if !IsRoom(room) {
		return nil, NotInARoom
	}

	// I need a copy to be referred by MoveRecord
	// taking &game.query is not an option, history would be modified

	query := Declaration{
		Character: character,
		Weapon:    weapon,
		Room:      room,
	}

	game.query = query
	game.answeringPlayer = game.NextAnsweringPlayer(game.currentPlayer)

	var moves []PlayerPosition

	for _, player := range game.players {
		if Card(player.character) == character {
			if player.position.Room == room {
				break
			}

			player.position.EnterRoom(room)

			moves = append(moves, PlayerPosition{
				PlayerID:     player.id,
				PawnPosition: player.position,
			})

			break
		}
	}

	record := &MoveRecord{
		PlayerID:  currentPlayer.id,
		Timestamp: time.Now(),
		Move: &QuerySolutionMove{
			Character: character,
			Weapon:    weapon,
		},
		StateDelta: StateUpdate{
			State:     game.state,
			Positions: moves,

			Query: &query,

			AnsweringPlayer: game.players[game.answeringPlayer].id,
		},
	}

	game.history = append(game.history, record)

	return record, nil
}

// nextTurnPlayer returns the next current player.
func (game *Game) nextTurnPlayer() (int, bool) {
	c := game.currentPlayer
	next := c

	for {
		next = (next + 1) % len(game.players)

		if next == c {
			return 0, true
		}

		if !game.players[next].FailedSolution() {
			return next, false
		}
	}
}

// NextAnsweringPlayer returns the player due to try to reveal a card after from player.
func (game *Game) NextAnsweringPlayer(from int) int {
	return (from + 1) % len(game.players)
}

// Reveal processes query solution answer.
func (game *Game) Reveal(card Card) (*MoveRecord, error) {
	answeringPlayer := game.players[game.answeringPlayer]

	if IsCard(card) {
		if !answeringPlayer.HasCard(card) {
			return nil, NotYourCard
		}

		game.state = GameStateTrySolution
		game.revealed = true
		game.revealedCard = card

		record := &MoveRecord{
			PlayerID:  answeringPlayer.id,
			Timestamp: time.Now(),
			Move: &RevealCardMove{
				Card: card,
			},
			StateDelta: StateUpdate{
				// current player did not change but I need it to decide
				// if send revealed card or not (only answering and current player see the revealed card)
				CurrentPlayer: game.CurrentPlayer().id,
				State:         game.state,
				Revealed:      game.revealed,
				RevealedCard:  game.revealedCard,
			},
		}

		game.history = append(game.history, record)

		return record, nil
	}

	//currentPlayer := game.Players[game.currentPlayer]

	if answeringPlayer.HasCard(game.query.Character) || answeringPlayer.HasCard(game.query.Room) || answeringPlayer.HasCard(game.query.Weapon) {
		return nil, MustShowACard
	}

	game.answeringPlayer = game.NextAnsweringPlayer(game.answeringPlayer)

	var nextAnsweringPlayerID PlayerID

	if game.answeringPlayer == game.currentPlayer {
		game.revealed = false
		game.revealedCard = NoCard
		game.state = GameStateTrySolution

		nextAnsweringPlayerID = 0

	} else {
		nextAnsweringPlayerID = game.AnsweringPlayer().id
	}

	record := &MoveRecord{
		PlayerID:  answeringPlayer.id,
		Timestamp: time.Now(),
		Move:      &NoCardToRevealMove{},
		StateDelta: StateUpdate{
			State:           game.state,
			AnsweringPlayer: nextAnsweringPlayerID,
			Revealed:        game.revealed,
			RevealedCard:    game.revealedCard,
		},
	}

	game.history = append(game.history, record)

	return record, nil
}

// Pass skips to the next turn.
func (game *Game) Pass() (*MoveRecord, error) {
	switch game.state {
	case GameStateQuery:
		if game.answeringPlayer != -1 {
			return nil, NotYourTurn
		}

		player := game.players[game.currentPlayer]

		game.state = GameStateTrySolution

		record := &MoveRecord{
			PlayerID:  player.id,
			Timestamp: time.Now(),
			Move:      &PassMove{},
			StateDelta: StateUpdate{
				State: game.state,
			},
		}

		game.history = append(game.history, record)

		return record, nil

	case GameStateTrySolution:
		game.state = GameStateNewTurn

		player := game.players[game.currentPlayer]

		nextPlayer, _ := game.nextTurnPlayer()

		game.currentPlayer = nextPlayer

		game.query = EmptyDeclaration

		game.revealed = false
		game.revealedCard = NoCard

		record := &MoveRecord{
			PlayerID:  player.id,
			Timestamp: time.Now(),
			Move:      &PassMove{},
			StateDelta: StateUpdate{
				State:         game.state,
				CurrentPlayer: game.players[game.currentPlayer].id,
			},
		}

		game.history = append(game.history, record)

		return record, nil

	default:
		return nil, IllegalState
	}
}

// CheckSolution verifies the solution.
func (game *Game) CheckSolution(character, room, weapon Card) ([]*MoveRecord, error) {
	if game.state != GameStateTrySolution {
		return nil, IllegalState
	}

	player := game.players[game.currentPlayer]
	player.declaration = &Declaration{
		Character: character,
		Room:      room,
		Weapon:    weapon,
	}

	historyLen := len(game.history)

	if game.solution == *player.declaration {
		game.state = GameEnded

		game.history = append(game.history, &MoveRecord{
			PlayerID:  player.id,
			Timestamp: time.Now(),
			Move: &DeclareSolutionMove{
				Declaration: *player.declaration,
			},
			StateDelta: StateUpdate{
				State: game.state,
			},
		})

	} else {
		nextPlayer, _ := game.nextTurnPlayer()

		game.currentPlayer = nextPlayer

		nextPlayerID := game.players[game.currentPlayer].id

		game.history = append(game.history, &MoveRecord{
			PlayerID:  player.id,
			Timestamp: time.Now(),
			Move: &DeclareSolutionMove{
				Declaration: *player.declaration,
			},
			StateDelta: StateUpdate{
				State:         GameStateNewTurn,
				CurrentPlayer: nextPlayerID,
			},
		})

		// check if nextPlayer is the last one who has not failed to find the solution

		if _, gameEnded := game.nextTurnPlayer(); gameEnded {
			game.state = GameEnded

			game.history = append(game.history, &MoveRecord{
				PlayerID:  nextPlayerID,
				Timestamp: time.Now(),
				Move: &DeclareSolutionMove{
					Declaration: game.solution,
				},
				StateDelta: StateUpdate{
					State: game.state,
				},
			})

		} else {
			game.state = GameStateNewTurn
		}
	}

	return game.history[historyLen:], nil
}

/*// HasWinner returns true if the game ended and it is not a draw. In this case the current player is the winner.
func (game *Game) HasWinner() bool {
	if game.state != GameEnded {
		return false
	}

	player := game.players[game.currentPlayer]

	return player.declaration != nil && *player.declaration == game.solution
}*/

// PlayerTurnSequence returns the order in which players play each turn.
func (game *Game) PlayerTurnSequence() []PlayerID {
	var playersOrder []PlayerID

	for _, player := range game.players {
		playersOrder = append(playersOrder, player.id)
	}

	return playersOrder
}

// StartState returns the initial Cluedo game state.
func (game *Game) StartState() StateUpdate {
	var positions []PlayerPosition

	for _, p := range game.players {
		positions = append(positions, PlayerPosition{
			PlayerID:     p.id,
			PawnPosition: initialPositions[p.character],
		})
	}

	return StateUpdate{
		State:         GameStateNewTurn,
		CurrentPlayer: game.players[0].id,
		Positions:     positions,
	}
}

// FullState return a fully compiled NotifyGameState so that web client can reinitialize from stratch.
// Usually NotifyGameState contains only incremental changes.
func (game *Game) FullState(askingPlayer PlayerID) StateUpdate {
	r := StateUpdate{
		State:         game.state,
		CurrentPlayer: game.players[game.currentPlayer].id,
	}

	switch game.state {
	case GameStateStarting:
		// nop
		break
	case GameStateNewTurn:
		r.Positions = game.PlayerPositions()
		break
	case GameStateCard:
		// TODO
		break
	case GameStateMove:
		r.Positions = game.PlayerPositions()
		r.Dice1 = game.dice1
		r.Dice2 = game.dice2
		r.RemainingSteps = game.remainingSteps
		break
	case GameStateQuery:
		r.Positions = game.PlayerPositions()
		r.Query = &game.query
		if game.answeringPlayer >= 0 {
			r.AnsweringPlayer = game.players[game.answeringPlayer].id
		}
		break
	case GameStateTrySolution:
		r.Positions = game.PlayerPositions()

		r.Query = &game.query
		r.Revealed = game.revealed
		if askingPlayer == game.players[game.currentPlayer].id {
			r.RevealedCard = game.revealedCard
		}
		break
	case GameEnded:
		r.Solution = &game.solution
		break
	}

	return r
}

// PlayerPositions return an array containing all players' position.
func (game *Game) PlayerPositions() []PlayerPosition {
	var r []PlayerPosition

	for _, player := range game.players {
		r = append(r, PlayerPosition{
			PlayerID:     player.id,
			PawnPosition: player.position,
		})
	}

	return r
}

// Players enumerate game players invoking mapPlayer on each of them.
func (game *Game) Players(mapPlayer func(player *Player)) {
	for _, player := range game.players {
		mapPlayer(player)
	}
}

// History enumerate game move records invoking recordHandler on each of them.
func (game *Game) History(recordHandler func(record MoveRecord)) {
	firstPlayer := game.players[0]

	recordHandler(MoveRecord{
		PlayerID:   firstPlayer.ID(),
		Timestamp:  time.Now(),
		Move:       &StartMove{},
		StateDelta: game.StartState(),
	})

	for _, record := range game.history {
		recordHandler(*record)
	}
}
