package game

//import "github.com/makeroo/my_clue_be/internal/platform/web"

// PawnPosition describes a pawn position on Cluedo board game.
// A player can be either in a room or in a hallway.
// When in a room, MapX/Y are set to 0, viceversa when in a hallway Room is set to 0/NoCard.
// Zeroed fields are not omitted so that JS client can use spread operator to and MapX/Y are used respectively.
// Note: 0 is neither a valid room neither a valid position.
type PawnPosition struct {
	Room Card `json:"room"`
	MapX int  `json:"map_x"`
	MapY int  `json:"map_y"`
}

// PositionAt builds a PawnPosition pointing at a cell in a hallway.
// Note: x/y are not validated.
func PositionAt(x, y int) PawnPosition {
	return PawnPosition{
		Room: NoCard,
		MapX: x,
		MapY: y,
	}
}

// InRoom returns true if the position is inside a room.
func (position PawnPosition) InRoom() bool {
	return IsRoom(position.Room)
}

// EnterRoom sets the pawn position inside the given room.
// Note: room is not validated.
func (position *PawnPosition) EnterRoom(room Card) {
	position.Room = room
	position.MapX = 0
	position.MapY = 0
}

// MoveTo sets the pawn position to the given coords.
// Note: mapX/Y are not validated.
func (position *PawnPosition) MoveTo(mapX, mapY int) {
	position.Room = 0
	position.MapX = mapX
	position.MapY = mapY
}

// IsAdjacent checks if this position is next to the given one.
func (position PawnPosition) IsAdjacent(x1, y1 int) bool {
	dx := x1 - position.MapX
	dy := y1 - position.MapY

	return (dx == 0 && (dy == -1 || dy == +1)) || (dy == 0 && (dx == -1 || dx == 1))
}

// PlayerID is the unique player idenfier assigned by Clue API server.
// A player can by identified by this id or by the ordinal in the players array that
// defines the turn cycle.
// A custom type helps to avoid ambiguity and signal mismatch when using the wrong id.
type PlayerID int

// Player collects all playing user data.
// A user can play multiple games simultaneously.
type Player struct {
	game *Game
	id   PlayerID
	// character is a valid character Card only when the user select an available character.
	character Card

	votedStart bool

	deck []Card

	position PawnPosition

	// UserIO is defined if the user is connected, nil otherwise.
	// Because UserIO is a websocket, this one-to-one binding limits to one tab per game.
	//UserIO *web.UserIO

	// User is a link to user's infos when user is offline and UserIO is nil
	//User *web.User

	declaration *Declaration
}

// ID returns the player id.
func (player *Player) ID() PlayerID {
	return player.id
}

// Game returns the game the player is in.
func (player *Player) Game() *Game {
	return player.game
}

// Character returns the character the player is.
func (player *Player) Character() Card {
	return player.character
}

// Deck returns the player deck.
// FIXME: this function breaks "private" fields design: deck is a pointer and can be modified
func (player *Player) Deck() []Card {
	return player.deck
}

// FailedSolution return wether the player is in play or not.
func (player *Player) FailedSolution() bool {
	return player.declaration != nil && *player.declaration != player.game.solution
}

// HasCard checks if the player has the card in her/his deck.
func (player *Player) HasCard(card Card) bool {
	for _, c := range player.deck {
		if c == card {
			return true
		}
	}

	return false
}
