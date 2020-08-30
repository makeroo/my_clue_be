package data

import (
	"github.com/makeroo/my_clue_be/internal/platform/game"
)

// MessageType is an enum of all the message types defined by My Clue BE API.
type MessageType string

const (
	// MessageSignInRequest is a constant for sign in request.
	MessageSignInRequest MessageType = "sign_in"
	// MessageSignInResponse is a constant for sign in response.
	MessageSignInResponse = "sign_in_response"

	// MessageCreateGameRequest is a constant for create game request.
	MessageCreateGameRequest = "create_game"
	// MessageCreateGameResponse is a constant for create game response.
	MessageCreateGameResponse = "create_game_resp"

	// MessageJoinGameRequest is a constant for join game request.
	MessageJoinGameRequest = "join_game"
	// MessageJoinGameResponse is a constant for join game response.
	MessageJoinGameResponse = "join_game_resp"

	// MessageSelectCharRequest is a constant for select char request.
	MessageSelectCharRequest = "select_char"

	// MessageVoteStartRequest is a constant for vote start request.
	MessageVoteStartRequest = "vote_start"

	// MessageRollDicesRequest is a constant for roll dices request.
	MessageRollDicesRequest = "roll_dices"

	// MessageMoveRequest is a constant for move request.
	MessageMoveRequest = "move"

	// MessageQuerySolutionRequest is a constant for query solution request.
	MessageQuerySolutionRequest = "query_solution"

	// MessageRevealRequest is a constant for reveal request.
	MessageRevealRequest = "reveal"

	// MessageDeclareSolutionRequest is a constant for declare solution request.
	MessageDeclareSolutionRequest = "declare_solution"

	// MessagePassRequest is a constant for pass request.
	MessagePassRequest = "pass"

	// MessageNotifyUserState is a constant for user state notification.
	MessageNotifyUserState = "notify_user_state"

	// MessageNotifyGameStarted is a constant for game started notification.
	MessageNotifyGameStarted = "notify_game_started"

	// MessageNotifyMoveRecord is a constant for move record notification.
	MessageNotifyMoveRecord = "notify_move_record"

	// MessageError is a constant for error notification.
	MessageError = "error"
)

// SignInRequest describes a sign in request.
// If only Name is defined, ie. non empty, then this is a register request and a new token will be
// assigned and returned in SignInResponse.
// If Token is defined, ie. non empty, then this is a sign in request.
type SignInRequest struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// SignInResponse describes a sign in response.
type SignInResponse struct {
	Token        string         `json:"token,omitempty"`
	RunningGames []GameSynopsis `json:"running_games,omitempty"`
}

// GamePlayer is a synthetic description of a Clue game player.
type GamePlayer struct {
	Character game.Card     `json:"character,omitempty"`
	ID        game.PlayerID `json:"player_id"`
	Name      string        `json:"name"`
	Online    bool          `json:"online"`
}

// GameSynopsis is a preview of a joined game.
type GameSynopsis struct {
	Game game.StateUpdate `json:"game"`
	//	Character game.Card     `json:"character,omitempty"`
	MyID    game.PlayerID `json:"my_player_id"`
	Players []GamePlayer  `json:"players,omitempty"`
}

// CreateGameResponse describes a create game response.
type CreateGameResponse struct {
	GameID string        `json:"game_id"`
	MyID   game.PlayerID `json:"my_player_id"`
}

// JoinGameRequest describes a join game request.
type JoinGameRequest struct {
	GameID string `json:"game_id"`
}

// JoinGameResponse describes a join game response.
type JoinGameResponse struct {
	Players []NotifyUserState `json:"players"`
	MyID    game.PlayerID     `json:"my_player_id"`
}

// SelectCharacterRequest describes a select char request.
type SelectCharacterRequest struct {
	Character game.Card `json:"character"`
}

// VoteStartRequest describes a vote start request.
type VoteStartRequest struct {
	Vote bool `json:"vote"`
}

// MoveRequest describes a move request.
type MoveRequest struct {
	EnterRoom game.Card `json:"enter_room"`
	MapX      int       `json:"map_x"`
	MapY      int       `json:"map_y"`
}

// QuerySolutionRequest describes a query solution request.
type QuerySolutionRequest struct {
	Character game.Card `json:"character"`
	Weapon    game.Card `json:"weapon"`
	// room is the room the querying player is in
}

// RevealRequest describes a reveal request.
type RevealRequest struct {
	Card game.Card `json:"card,omitempty"`
}

// DeclareSolutionRequest describes a declare solution request.
type DeclareSolutionRequest struct {
	game.Declaration
}

// NotifyError is an error message.
type NotifyError struct {
	Error string `json:"error"`
}

/*
// PlayerDeclaration is
type PlayerDeclaration struct {
	Declaration
	ID game.PlayerID `json:"player_id"`
}*/

/*
NotifyUserState communicates weather a user is reachable or not (Online),
and if she/he changed name or selected a character.
The message is sent to all users playing a game with her/him.
A user is identified by PlayerID.
Name can't be used because is user choosen an not guaranteed to be unique.
Character is unique only when game is started, otherwise is not defined,
ie. it is 0 for all players that haven't selected a character yet.
PlayerID is unique only in a game.
*/
type NotifyUserState struct {
	ID        game.PlayerID `json:"player_id"`
	Name      string        `json:"name,omitempty"`
	Character game.Card     `json:"character,omitempty"`
	Online    bool          `json:"online"`
}

// NotifyGameStarted is sent to all players of a table to signal that
// the game has started.
type NotifyGameStarted struct {
	Deck         []game.Card     `json:"deck"`
	PlayersOrder []game.PlayerID `json:"players_order"`
}

// MessageFrame is a message going from fe to be or vicersa.
// Body can be nil (eg. create game or pass requests) or an instance of
// the types above.
type MessageFrame struct {
	Header MessageHeader
	Body   interface{}
}

// MessageHeader is the first value  sent throught the web socket.
type MessageHeader struct {
	Type  MessageType `json:"type"`
	ReqID int         `json:"req_id"`
}
