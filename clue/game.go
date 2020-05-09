package clue

import (
	"errors"
	"math/rand"
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

var initialPositions = map[Card][2]int{
	MissScarlett: {7, 24},
	RevGreen:     {14, 0},
	ColMustard:   {0, 17},
	ProfPlum:     {23, 19},
	MrsPeacock:   {23, 7},
	MrsWhite:     {9, 0},
}

// Player collects all playing user data.
// A user can play multiple games simultaneously.
type Player struct {
	Game      *Game
	PlayerID  int
	Character int

	VoteStart bool

	Deck []Card

	// A player can either be in a room or somewhere in the hallways.
	// When in a room EnteredRoom holds the room card the player is in.
	// Otherwise EnteredRoom is -1 and MapX/Y point to the map cell.
	EnteredRoom int
	MapX        int
	MapY        int

	UserIO *UserIO
}

// Game is a clue table. A user can join multiple tables.
type Game struct {
	GameID  string
	Players map[int]*Player
	rand    *rand.Rand

	solutionRoom      Card
	solutionWeapon    Card
	solutionCharacter Card

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
func NewGame(gameID string, rand *rand.Rand) *Game {
	game := Game{
		GameID:  gameID,
		Players: make(map[int]*Player),
		rand:    rand,

		state: GameStateStarting,
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
		Game:     game,
		UserIO:   userIO,
		PlayerID: len(game.Players) + 1,
	}

	game.Players[player.PlayerID] = player

	return player, nil
}

func (game *Game) Start() {
	game.shufflePlayers()

	game.state = GameStateNewTurn
	game.currentPlayer = 0

	// create secret

	game.solutionCharacter = game.randomCard(MissScarlett, MrsWhite)
	game.solutionRoom = game.randomCard(Kitchen, Study)
	game.solutionWeapon = game.randomCard(Candlestick, Wrenck)

	deck := game.makeDeck()

	game.rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	cardsPerPlayer := len(deck) / len(game.Players)
	playersWithAnExtraCard := len(deck) % len(game.Players)
	start := 0

	for i, player := range game.Players {
		cards := cardsPerPlayer

		if i < playersWithAnExtraCard {
			cards++
		}

		player.Deck = deck[start : start+cards]
		start += cards
	}

	for _, player := range game.Players {
		player.EnteredRoom = -1

		pos := initialPositions[Card(player.Character)]

		player.MapX, player.MapY = pos[0], pos[1]
	}
}

func (game *Game) shufflePlayers() {
	game.rand.Shuffle(len(game.Players), func(i, j int) {
		game.Players[i], game.Players[j] = game.Players[j], game.Players[i]
	})
}

func (game *Game) randomCard(min Card, max Card) Card {
	n := game.rand.Intn(int(max)-int(min)) + int(min)

	return Card(n)
}

func (game *Game) makeDeck() []Card {
	deck := make([]Card, int(Wrenck))

	for i := range deck {
		deck[i] = Card(i)
	}

	return deck
}

/*func (game *Game) randomPlayerID() int {
	for {
		x := utils.RandomInt()

		if game.Players[x] == nil {
			return x
		}
	}
}*/

func (game *Game) SelectCharacter(player *Player, character int) (bool, error) {
	if game.state != GameStateStarting {
		return false, errors.New(GameAlreadyStarted)
	}

	if !IsCharacter(Card(character)) {
		return false, errors.New(NotACharacter)
	}

	if player.Character == character {
		return false, nil
	}

	for _, cplayer := range game.Players {
		if cplayer.Character == character {
			return false, errors.New(AlreadySelected)
		}
	}

	player.Character = character

	return true, nil
}

func (game *Game) VoteStart(player *Player, vote bool) (bool, error) {
	if game.state != GameStateStarting {
		return false, errors.New(GameAlreadyStarted)
	}

	if !IsCharacter(Card(player.Character)) {
		return false, errors.New(CharacterNotSelected)
	}

	if player.VoteStart == vote {
		return false, nil
	}

	player.VoteStart = vote

	if vote && len(game.Players) > 2 { // TODO: support 2 player version
		for _, p := range game.Players {
			if !p.VoteStart {
				return false, nil
			}
		}
	}

	return true, nil
}
