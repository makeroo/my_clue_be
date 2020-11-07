package game

import (
	"encoding/json"
	"time"
	//"github.com/my_clue_be/internal/platform/web"
)

// MoveType is an identifier for the actions the players can enact in the game.
type MoveType int

const (
	// Start action.
	Start MoveType = iota + 1
	// RollDices action.
	RollDices
	// MovingInTheHallway action: either just exited a room or continuing in a corridor.
	MovingInTheHallway
	// EnterRoom action.
	EnterRoom
	// QuerySolution action: a player investigation.
	QuerySolution
	// NoCardToReveal action: answering player has no card to reveal to querying player.
	NoCardToReveal
	// RevealCard action.
	RevealCard
	// DeclareSolution action: the player will end her/his game wether as winner or not.
	DeclareSolution
	// Pass action. Used to skip investigation and solution declaration.
	// Note: to remain in a room after having rolled the dices, a player use EnterRoom action specifying the same room she/he is in.
	Pass
)

// Move is a marker.
type Move interface {
	MoveType() MoveType
}

// StartMove is a marker for the start of game record.
type StartMove struct{}

// MoveType returns Start action.
func (start *StartMove) MoveType() MoveType {
	return Start
}

// RollDicesMove describes dice rolling result.
type RollDicesMove struct {
	Dice1 int
	Dice2 int
}

// MoveType returns RollDices action.
func (move *RollDicesMove) MoveType() MoveType {
	return RollDices
}

// MovingInTheHallwayMove describes a pawn move in a corridor.
type MovingInTheHallwayMove struct {
	MapX int
	MapY int
}

// MoveType returns MovingInTheHallway action.
func (move *MovingInTheHallwayMove) MoveType() MoveType {
	return MovingInTheHallway
}

// EnterRoomMove describes a player entering a room or remaining in the same room she/he was in.
type EnterRoomMove struct {
	Room Card
}

// MoveType returns EnterRoom.
func (move *EnterRoomMove) MoveType() MoveType {
	return EnterRoom
}

// QuerySolutionMove describes a player investigation.
// Because a player can investigate only if she/he is in a room and can investigate that room only
// room is implicit and not specified in the action.
type QuerySolutionMove struct {
	Character Card
	Weapon    Card
}

// MoveType eturns QuerySolution action.
func (move *QuerySolutionMove) MoveType() MoveType {
	return QuerySolution
}

// PassMove describes a pass action.
type PassMove struct {
}

// MoveType returns Pass action.
func (move *PassMove) MoveType() MoveType {
	return Pass
}

// RevealCardMove describes the card the answering player shows to the querying one.
type RevealCardMove struct {
	Card Card
}

// MoveType returns RevealCard action.
func (move *RevealCardMove) MoveType() MoveType {
	return RevealCard
}

// NoCardToRevealMove signals that the answering player has none of the queried cards.
type NoCardToRevealMove struct {
}

// MoveType returns NoCardToReveal action.
func (move *NoCardToRevealMove) MoveType() MoveType {
	return NoCardToReveal
}

// DeclareSolutionMove describes a solution declaration.
type DeclareSolutionMove struct {
	Declaration
}

// MoveType returns DeclareSolution action.
func (move *DeclareSolutionMove) MoveType() MoveType {
	return DeclareSolution
}

// MoveRecord comprises of the player executing the action, the time she/he did it, which action executed, and its results.
type MoveRecord struct {
	PlayerID   PlayerID    `json:"player_id"`
	Timestamp  time.Time   `json:"timestamp"`
	Move       Move        `json:"move,omitempty"`
	StateDelta StateUpdate `json:"state_delta"`
}

// Declaration is a triple of cards. They must be a character card, a room card and a weapon card.
type Declaration struct {
	Room      Card `json:"room"`
	Weapon    Card `json:"weapon"`
	Character Card `json:"character"`
}

// EmptyDeclaration is a constant used to "reset" a declaration.
var EmptyDeclaration = Declaration{
	Room:      NoCard,
	Character: NoCard,
	Weapon:    NoCard,
}

// PlayerPosition describes a player pawn position.
type PlayerPosition struct {
	PlayerID PlayerID `json:"player_id"`
	PawnPosition
}

// PlayerDeclaration describes a player declaration.
type PlayerDeclaration struct {
	PlayerID PlayerID `json:"player_id"`
	Declaration
}

// StateUpdate describes the effects of an action or the current game state.
// Only the relevant fields in each cases are defined.
type StateUpdate struct {
	State         State    `json:"state"`
	CurrentPlayer PlayerID `json:"current_player,omitempty"`

	Dice1          int `json:"dice1,omitempty"`
	Dice2          int `json:"dice2,omitempty"`
	RemainingSteps int `json:"remaining_steps,omitempty"`

	Positions    []PlayerPosition    `json:"positions,omitempty"`
	Declarations []PlayerDeclaration `json:"declarations,omitempty"`

	Query           *Declaration `json:"query,omitempty"`
	AnsweringPlayer PlayerID     `json:"answering_player,omitempty"`
	Revealed        bool         `json:"revealed,omitempty"`
	RevealedCard    Card         `json:"revealed_card,omitempty"`

	Solution *Declaration `json:"solution,omitempty"`
}

// AsMessageFor return a record containing only the informations visible by the specified player.
func (record MoveRecord) AsMessageFor(player *Player) MoveRecord {
	if record.PlayerID == player.id {
		return record
	}

	if record.StateDelta.State == GameEnded {
		return record
	}

	if !IsCard(record.StateDelta.RevealedCard) {
		return record
	}

	// otherwise revealed card is not visible to the other players
	r := record

	r.StateDelta.RevealedCard = NoCard

	return r
}

// JSONMoveRecord is an alias used to break cycles in MarshalJSON invocation.
// The trick is required to add a field to the struct implementing Move interface
// containing MoveType() result.
type JSONMoveRecord MoveRecord

// MarshalJSON produces a json comprising of MoveRecord declared fields plus one declaring move type.
func (record MoveRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		JSONMoveRecord
		Type MoveType `json:"type"`
	}{
		JSONMoveRecord: JSONMoveRecord(record),
		Type:           record.Move.MoveType(),
	})
}
