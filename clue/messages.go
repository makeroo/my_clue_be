package clue

const (
	// MessageSignInRequest is a constant for sign in request.
	MessageSignInRequest = "sign_in"
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

	// MessageNotifyGameState is a constant for game state notification.
	MessageNotifyGameState = "notify_game_state"

	// MessageError is a constant for error notification.
	MessageError = "error"
)

const (
	// NotSignedIn error: issued for game related requests.
	NotSignedIn = "not_signed_in"
	// CannotJoinRunningGame error: illegal join game request.
	CannotJoinRunningGame = "cannot_join_running_game"
	// TableIsFull error: cannot join.
	TableIsFull = "table_is_full"
	// TokenMismatch error: the user provided a different token.
	TokenMismatch = "token_mismatch"
	// UnknownToken error: the provided token is unknown (maybe expired?).
	UnknownToken = "unknown_token"
	// TooManyGames error: per user game limit exceeded.
	TooManyGames = "too_many_games"
	// UnknownGame error: illegal join request.
	UnknownGame = "unknown_game"
	// AlreadyPlaying error: one tab per game limit.
	AlreadyPlaying = "already_playing"
	// AlreadySelected error: the choosen character has already been selected.
	AlreadySelected = "already_selected"
	// NotACharacter error: illegal select char request.
	NotACharacter = "not_a_character"
	// NotPlaying error: illegal game related request.
	NotPlaying = "not_playing"
	// GameAlreadyStarted error: illegal game setup request (eg. vote start / select char).
	GameAlreadyStarted = "game_already_started"
	// CharacterNotSelected error: cannot vote start before having selected a character.
	CharacterNotSelected = "character_not_selected"
	// NotYourTurn error: illegal game related request.
	NotYourTurn = "not_your_turn"
	// IllegalState error: illegal game related request.
	IllegalState = "illegal_state"
	// IllegalMove error: illegal move parameters.
	IllegalMove = "illegal_move"
	// NotYourCard error: cannot reveal a card not in your deck.
	NotYourCard = "not_your_card"
	// MustShowACard error: cannot pass if you have a card to show.
	MustShowACard = "must_show_a_card"
	// NotInARoom error: cannot query solution if you are not in a room.
	NotInARoom = "not_in_a_room"
)

// SignInRequest describes a sign in request.
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
	Character int    `json:"character,omitempty"`
	PlayerID  int    `json:"player_id"`
	Name      string `json:"name"`
	Online    bool   `json:"online"`
}

// GameSynopsis is a preview of a joined game.
type GameSynopsis struct {
	GameID    string       `json:"game_id"`
	Character int          `json:"character,omitempty"`
	PlayerID  int          `json:"player_id"`
	Others    []GamePlayer `json:"others,omitempty"`
}

// CreateGameResponse describes a create game response.
type CreateGameResponse struct {
	GameID   string `json:"game_id"`
	PlayerID int    `json:"player_id"`
}

// JoinGameRequest describes a join game request.
type JoinGameRequest struct {
	GameID string `json:"game_id"`
}

// JoinGameResponse describes a join game response.
type JoinGameResponse struct {
	Players  []NotifyUserState `json:"players"`
	PlayerID int               `json:"player_id"`
}

// SelectCharacterRequest describes a select char request.
type SelectCharacterRequest struct {
	Character int `json:"character"`
}

// VoteStartRequest describes a vote start request.
type VoteStartRequest struct {
	Vote bool `json:"vote"`
}

// MoveRequest describes a move request.
type MoveRequest struct {
	EnterRoom Card `json:"enter_room"`
	MapX      int  `json:"map_x"`
	MapY      int  `json:"map_y"`
}

// QuerySolutionRequest describes a query solution request.
type QuerySolutionRequest struct {
	Character Card `json:"character"`
	Weapon    Card `json:"weapon"`
}

// RevealRequest describes a reveal request.
type RevealRequest struct {
	Card Card `json:"card,omitempty"`
}

// DeclareSolutionRequest describes a declare solution request.
type DeclareSolutionRequest struct {
	Character Card `json:"character"`
	Room      Card `json:"room"`
	Weapon    Card `json:"weapon"`
}

// NotifyError is an error message.
type NotifyError struct {
	Error string `json:"error"`
}

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
	PlayerID  int    `json:"player_id"`
	Name      string `json:"name,omitempty"`
	Character int    `json:"character,omitempty"`
	Online    bool   `json:"online"`
}

// NotifyGameStarted is sent to all players of a table to signal that
// the game has started.
type NotifyGameStarted struct {
	Deck         []Card `json:"deck"`
	PlayersOrder []int  `json:"players_order"`
}

// PlayerPosition notifies a player position.
// A player can be either in a room or in a hallway.
// Room and MapX/Y are used respectively.
type PlayerPosition struct {
	PlayerID int  `json:"player_id"`
	Room     Card `json:"room,omitempty"`
	MapX     int  `json:"map_x,omitempty"`
	MapY     int  `json:"map_y,omitempty"`
}

// NotifyGameState is sent to all players of a table whenever something
// happens.
type NotifyGameState struct {
	State         State `json:"state"`
	CurrentPlayer int   `json:"current_player,omitempty"`

	Dice1          int `json:"dice1,omitempty"`
	Dice2          int `json:"dice2,omitempty"`
	RemainingSteps int `json:"remaining_steps,omitempty"`

	PlayerPositions []PlayerPosition `json:"player_positions,omitempty"`

	AnsweringPlayer int  `json:"answering_player,omitempty"`
	Character       Card `json:"character,omitempty"`
	Room            Card `json:"room,omitempty"`
	Weapon          Card `json:"weapon,omitempty"`
	Revealed        bool `json:"revealed,omitempty"`
	RevealedCard    Card `json:"revealed_card,omitempty"`
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
	Type  string `json:"type"`
	ReqID int    `json:"req_id"`
}
