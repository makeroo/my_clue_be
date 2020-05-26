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

// Card is a card in the solution deck that comprises of all the characters,
// rooms and weapons. I prefer this solution instead of having three enums,
// one for each kind of object and a Deck of "union".
type Card int

const (
	// Candlestick is the candlestick card
	Candlestick Card = iota + 1
	// Knife is the knife card
	Knife
	// LeadPipe is the lead pipe card
	LeadPipe
	// Revolver is the revolver card
	Revolver
	// Rope is the rope card
	Rope
	// Wrenck is the wrenck card
	Wrenck

	// Kitchen is the kitchen card
	Kitchen
	// Ballroom is the ballroom card
	Ballroom
	// Conservatory is the conservatory card, the greenhouse in the ITA version.
	Conservatory
	// DiningRoom is the dining room card
	DiningRoom
	// BilliardRoom is the billiard card
	BilliardRoom
	// Library is the library card
	Library
	// Lounge is the lounge card
	Lounge
	// Hall is an hall card
	Hall
	// Study is the study card
	Study

	// MissScarlett is the Miss Scarlett card
	MissScarlett
	// RevGreen is the Rev. Green card
	RevGreen
	// ColMustard is the Col. Mustard card
	ColMustard
	// ProfPlum is the Prof. Plum card
	ProfPlum
	// MrsPeacock is the Mrs. Peacock card
	MrsPeacock
	// MrsWhite is the Mrs. White card
	MrsWhite // ITA: Orchid
)

// Cards is the number of cards in the room+weapon+char deck.
const Cards = int(MrsWhite)

var initialPositions = map[Card][2]int{
	MissScarlett: {7, 24},
	RevGreen:     {14, 0},
	ColMustard:   {0, 17},
	ProfPlum:     {23, 19},
	MrsPeacock:   {23, 7},
	MrsWhite:     {9, 0},
}

const (
	xx = -1
	oo = 0
	ki = int(Kitchen)
	ba = int(Ballroom)
	co = int(Conservatory)
	di = int(DiningRoom)
	bi = int(BilliardRoom)
	li = int(Library)
	lo = int(Lounge)
	ha = int(Hall)
	st = int(Study)
)

var clueBoard = [25][24]int{
	//0   1   2   3   4   5   6   7   8   9  10  11  12  13  14  15  16  17  18  19  20  21  22  23
	{xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx}, // 00
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 01
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx}, // 02
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx}, // 03
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx}, // 04
	{xx, xx, xx, xx, xx, xx, oo, ba, xx, xx, xx, xx, xx, xx, xx, xx, ba, oo, xx, xx, xx, xx, xx, xx}, // 05
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, co, oo, oo, oo, oo, xx}, // 06
	{oo, oo, oo, oo, ki, oo, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, oo, oo, xx}, // 07
	{xx, oo, oo, oo, oo, oo, oo, oo, oo, ba, oo, oo, oo, oo, ba, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 08
	{xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, bi, xx, xx, xx, xx, xx, xx}, // 09
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 10
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 11
	{xx, xx, xx, xx, xx, xx, xx, xx, di, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 12

	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, li, oo, bi, xx}, // 13
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 14
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 15
	{xx, oo, oo, oo, oo, oo, di, oo, oo, oo, xx, xx, xx, xx, xx, oo, li, xx, xx, xx, xx, xx, xx, xx}, // 16
	{xx, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, ha, ha, oo, oo, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 17
	{xx, oo, oo, oo, oo, oo, lo, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 18
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, oo, oo, oo, xx}, // 19
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, st, oo, oo, oo, oo, oo, xx}, // 20
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 21
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 22
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 23
	{xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, oo, xx, xx, xx, xx, xx, xx, xx}, // 24
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
	// When in a room Room holds the room card the player is in.
	// Otherwise Room is 0 and MapX/Y point to the map cell.
	Room Card
	MapX int
	MapY int

	// UserIO is defined if the user is connected, nil otherwise.
	UserIO *UserIO

	// User is a link to user's infos when user is offline and UserIO is nil
	User *User

	// FailedSolution is false until the player declared a solution
	// that is wrong.
	FailedSolution bool
}

// Game is a clue table. A user can join multiple tables.
type Game struct {
	GameID  string
	Players []*Player
	rand    *rand.Rand

	solutionRoom      Card
	solutionWeapon    Card
	solutionCharacter Card

	state          State
	currentPlayer  int
	dice1          int
	dice2          int
	remainingSteps int
	queryRoom      Card
	queryWeapon    Card
	queryCharacter Card
	queryingPlayer int

	secretPassages [][2]Card
}

// IsRoom returns true if the given card is a room.
func IsRoom(card Card) bool {
	return !IsWeapon(card) && !IsCharacter(card)
}

// IsWeapon returns true if the given card is a weapon.
func IsWeapon(card Card) bool {
	return 0 < card && card < Kitchen
}

// IsCharacter returns true if the given card is a character.
func IsCharacter(card Card) bool {
	return card > Study
}

// IsCard returns true if the card is valid.
// Used when casting from int (json).
func IsCard(card Card) bool {
	return card >= Candlestick && card <= MrsWhite
}

// HasCard checks if the player has the card in her/his deck.
func (player *Player) HasCard(card Card) bool {
	for _, c := range player.Deck {
		if c == card {
			return true
		}
	}

	return false
}

// NewGame create a Game instance.
func NewGame(gameID string, rand *rand.Rand) *Game {
	game := Game{
		GameID: gameID,
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

// AddPlayer adds a player to the table.
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
		User:     userIO.user,
		PlayerID: len(game.Players) + 1,
	}

	game.Players = append(game.Players, player)

	return player, nil
}

// Start starts a new game.
func (game *Game) Start() {
	game.shufflePlayers()

	game.state = GameStateNewTurn
	game.currentPlayer = 0

	// create secret

	game.solutionCharacter = game.randomCard(MissScarlett, MrsWhite)
	game.solutionRoom = game.randomCard(Kitchen, Study)
	game.solutionWeapon = game.randomCard(Candlestick, Wrenck)

	deck := game.makeDeckWithoutSolution()

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
		player.Room = Candlestick // any non room card is fine

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

func (game *Game) makeDeckWithoutSolution() []Card {
	deck := make([]Card, Cards-3)

	p := 0

	for i := 0; i < Cards; i++ {
		c := Card(i)

		if c == game.solutionCharacter || c == game.solutionRoom || c == game.solutionWeapon {
			continue
		}

		deck[p] = Card(i)
		p++
	}

	return deck
}

// SelectCharacter assigns a character to a player.
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

// VoteStart records the start vote of a player.
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

	if !vote || len(game.Players) < 3 { // TODO: support 2 player version
		return false, nil
	}

	for _, p := range game.Players {
		if !p.VoteStart {
			return false, nil
		}
	}

	return true, nil
}

// IsCurrentPlayer checks if it is the turn of the player.
func (game *Game) IsCurrentPlayer(player *Player) bool {
	return game.Players[game.currentPlayer].PlayerID == player.PlayerID
}

// RollDices rolls dices for current player.
func (game *Game) RollDices() error {
	if game.state != GameStateNewTurn {
		return errors.New(IllegalState)
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

	return nil
}

// Move moves current player.
func (game *Game) Move(room Card, mapX int, mapY int) error {
	if game.state != GameStateMove {
		return errors.New(IllegalState)
	}

	player := game.Players[game.currentPlayer]

	if IsRoom(room) {
		if room == player.Room {
			// the player choose to remain in the same room she/he was in
			game.state = GameStateQuery
			game.queryingPlayer = -1

			return nil
		}

		if IsRoom(player.Room) {
			// changing room is possible only through secret passages

			if !game.IsSecretPassage(player.Room, room) {
				return errors.New(IllegalMove)
			}

			player.Room = room
			player.MapX = 0
			player.MapY = 0

			game.state = GameStateQuery
			game.queryingPlayer = -1

			return nil
		}

		// the player is in the hallway, check if she/he is in front of a door of
		// the room she/he wants to enter in

		cellType := Card(clueBoard[player.MapY][player.MapX])

		if cellType != room {
			return errors.New(IllegalMove)
		}

		player.Room = room
		player.MapX = 0
		player.MapY = 0

		game.state = GameStateQuery
		game.queryingPlayer = -1

	} else if IsRoom(player.Room) {
		if !game.IsValidPosition(mapX, mapY) {
			return errors.New(IllegalMove)
		}

		if p := game.IsOccupied(mapX, mapY); p != nil && p != player {
			return errors.New(IllegalMove)
		}

		// the player just exited a room
		// check the hallway pos she/he selected is one in front of a door of
		// the room she/he was in
		cellType := Card(clueBoard[mapY][mapX])

		if cellType != player.Room {
			return errors.New(IllegalMove)
		}

		player.Room = 0
		player.MapX = mapX
		player.MapY = mapY

		game.remainingSteps--

		// if the player has just exited a room then this is the first step
		// no need to check remainingSteps because the min is 2 so it is at least 1

	} else {
		if !game.IsValidPosition(mapX, mapY) {
			return errors.New(IllegalMove)
		}

		if clueBoard[mapY][mapX] < 0 {
			return errors.New(IllegalMove)
		}

		if p := game.IsOccupied(mapX, mapY); p != nil && p != player {
			return errors.New(IllegalMove)
		}

		if !IsPositionAdjacent(player.MapX, player.MapY, mapX, mapY) {
			return errors.New(IllegalMove)
		}

		player.MapX = mapX
		player.MapY = mapY

		game.remainingSteps--

		if game.remainingSteps == 0 {
			game.state = GameStateQuery
			game.queryingPlayer = -1
		}
	}

	return nil
}

// IsValidPosition checks coordinate ranges.
func (game *Game) IsValidPosition(mapX, mapY int) bool {
	return mapX < 0 || mapX > 23 || mapY < 0 || mapY > 24
}

// IsPositionAdjacent checks if x/y2 is next to x/y1.
func IsPositionAdjacent(x1, y1, x2, y2 int) bool {
	dx := x1 - x2
	dy := y1 - y2

	return (dx == 0 && (dy == -1 || dy == +1)) || (dy == 0 && (dx == -1 || dx == 1))
}

// IsOccupied checks if position is occupied by a player.
func (game *Game) IsOccupied(mapX, mapY int) *Player {
	for _, p := range game.Players {
		if p.MapX == mapX || p.MapY == mapY {
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
func (game *Game) QuerySolution(character, weapon Card) error {
	if game.state != GameStateQuery || game.queryingPlayer != -1 {
		return errors.New(IllegalState)
	}

	room := game.Players[game.currentPlayer].Room
	if !IsRoom(room) {
		return errors.New(NotInARoom)
	}

	game.queryCharacter = character
	game.queryRoom = room
	game.queryWeapon = weapon
	// TODO: not len(players) but playingPlayers
	// because we have to exclude those who tried and failed and are not playing anymore
	game.queryingPlayer = game.NextPlayer()

	return nil
}

// NextPlayer returns the next current player.
func (game *Game) NextPlayer() int {
	for {
		next := (game.currentPlayer + 1) % len(game.Players)

		if !game.Players[next].FailedSolution {
			return next
		}
	}
}

// Reveal processes query solution answer.
func (game *Game) Reveal(card Card) (bool, error) {
	answeringPlayer := game.Players[game.queryingPlayer]

	if IsCard(card) {
		if !answeringPlayer.HasCard(card) {
			return false, errors.New(NotYourCard)
		}

		game.state = GameStateTrySolution

		return true, nil

	}

	currentPlayer := game.Players[game.currentPlayer]

	if currentPlayer.HasCard(game.queryCharacter) || currentPlayer.HasCard(game.queryRoom) || currentPlayer.HasCard(game.queryWeapon) {
		return false, errors.New(MustShowACard)
	}

	game.queryingPlayer = game.NextPlayer()

	if game.queryingPlayer == game.currentPlayer {
		game.state = GameStateTrySolution
	}

	return false, nil
}

// Pass skips to the next turn.
func (game *Game) Pass() error {
	switch game.state {
	case GameStateQuery:
		if game.queryingPlayer != -1 {
			return errors.New(NotYourTurn)
		}

		game.state = GameStateTrySolution
		return nil

	case GameStateTrySolution:
		game.state = GameStateNewTurn
		game.currentPlayer = game.NextPlayer()
		return nil

	default:
		return errors.New(IllegalState)
	}
}

// CheckSolution verifies the solution.
func (game *Game) CheckSolution(character, room, weapon Card) error {
	if game.state != GameStateTrySolution {
		return errors.New(IllegalState)
	}

	if game.solutionCharacter == character && game.solutionRoom == room && game.solutionWeapon == weapon {
		game.state = GameEnded

	} else {
		player := game.Players[game.currentPlayer]

		player.FailedSolution = true

		game.state = GameStateNewTurn
		game.currentPlayer = game.NextPlayer()
	}

	return nil
}

// GameStartedMessage return a NotifyGameStarted message so that web client can reinitialize from scratch.
func (game *Game) GameStartedMessage(player *Player) NotifyGameStarted {
	var playersOrder []int

	for _, player := range game.Players {
		playersOrder = append(playersOrder, player.PlayerID)
	}

	return NotifyGameStarted{
		Deck:         player.Deck,
		PlayersOrder: playersOrder,
	}
}

// FullState return a fully compiled NotifyGameState so that web client can reinitialize from stratch.
// Usually NotifyGameState contains only incremental changes.
func (game *Game) FullState() NotifyGameState {
	r := NotifyGameState{
		State:         game.state,
		CurrentPlayer: game.currentPlayer,
	}

	switch game.state {
	case GameStateStarting:
		// nop
		break
	case GameStateNewTurn:
		r.PlayerPositions = game.PlayerPositions()
		break
	case GameStateCard:
		// nop
		break
	case GameStateMove:
		r.PlayerPositions = game.PlayerPositions()
		r.Dice1 = game.dice1
		r.Dice2 = game.dice2
		r.RemainingSteps = game.remainingSteps
		break
	case GameStateQuery:
		r.PlayerPositions = game.PlayerPositions()
		r.Room = game.queryRoom
		r.Weapon = game.queryWeapon
		r.Character = game.queryCharacter
		r.AnsweringPlayer = game.queryingPlayer
		break
	case GameStateTrySolution:
		r.PlayerPositions = game.PlayerPositions()
		break
	case GameEnded:
		// nop
		break
	}

	return r
}

// PlayerPositions return an array containing all players' position.
func (game *Game) PlayerPositions() []PlayerPosition {
	var r []PlayerPosition

	for _, player := range game.Players {
		r = append(r, PlayerPosition{
			PlayerID: player.PlayerID,
			Room:     player.Room,
			MapX:     player.MapX,
			MapY:     player.MapY,
		})
	}

	return r
}
