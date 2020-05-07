package clue

import (
	"errors"
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
	// GameStateAction is the state entered when dice have rolled, eventually
	// the hints card have been obeyed and the player have to choose what to do.
	GameStateAction
	// GameStateQuery is the state entered when the player completed her/his
	// move and, if in a room, query a solution.
	GameStateQuery
	// TODO: assert solution
)

// Card is a card in the solution deck that comprises of all the characters,
// rooms and weapons. I prefer this solution instead of having three enums,
// one for each kind of object and a Deck of "union".
type Card int

const (
	// Kitchen is the kitchen room
	Kitchen Card = iota
	// Ballroom is the ballroom
	Ballroom
	// Conservatory is the conservatory room, the greenhouse in the ITA version.
	Conservatory
	// DiningRoom is the dining room.
	DiningRoom
	BilliardRoom
	Library
	Lounge
	Hall
	Study

	MissScarlett
	RevGreen
	ColMustard
	ProfPlum
	MrsPeacock
	MrsWhite // ITA: Orchid

	Candlestick
	Knife
	LeadPipe
	Revolver
	Rope
	Wrenck
)

// Player collects all playing user data.
// A user can play multiple games simultaneously.
type Player struct {
	Game      *Game
	Character int
	Deck      []Card

	UserIO *UserIO
}

// Game is a clue table. A user can join multiple tables.
type Game struct {
	GameID   string
	Players  []*Player
	solution []Card

	state          State
	currentPlayer  int
	dice1          int
	dice2          int
	remainingSteps int
}

// IsRoom return true if the given card is a room.
func IsRoom(card Card) bool {
	return card < MissScarlett
}

// IsWeapon return true if the given card is a weapon.
func IsWeapon(card Card) bool {
	return card > MrsWhite
}

// IsCharacter return true if the given card is a character.
func IsCharacter(card Card) bool {
	return !IsRoom(card) && !IsWeapon(card)
}

// NewGame create a Game instance.
func NewGame(gameID string) *Game {
	game := Game{
		GameID:         gameID,
		Players:        nil,
		solution:       nil,
		state:          GameStateStarting,
		currentPlayer:  0,
		dice1:          0,
		dice2:          0,
		remainingSteps: 0,
	}

	return &game
}

func (game *Game) AddPlayer(userIO *UserIO) (*Player, error) {
	if game.state != GameStateStarting {
		return nil, errors.New(CannotJoinRunningGame)
	}

	if len(game.Players) == 6 {
		return nil, errors.New(TableIsFull)
	}

	player := &Player{
		Game:   game,
		UserIO: userIO,
	}

	game.Players = append(game.Players, player)

	return player, nil
}

//func (game *Game) SelectCharacter(playerToken string, character int) {
//
//}
