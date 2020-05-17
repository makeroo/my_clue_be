package clue

const (
	MessageSignInRequest  = "sign_in"
	MessageSignInResponse = "sign_in_response"

	MessageCreateGameRequest  = "create_game"
	MessageCreateGameResponse = "create_game_resp"

	MessageJoinGameRequest  = "join_game"
	MessageJoinGameResponse = "join_game_resp"

	MessageSelectCharRequest = "select_char"

	MessageVoteStartRequest = "vote_start"

	MessageRollDicesRequest = "roll_dices"

	MessageMoveRequest = "move"

	MessageQuerySolutionRequest = "query_solution"

	MessageRevealRequest = "reveal"

	MessageDeclareSolutionRequest = "declare_solution"

	MessagePassRequest = "pass"

	MessageNotifyUserState = "notify_user_state"

	MessageNotifyGameStarted = "notify_game_started"

	MessageNotifyGameState = "notify_game_state"

	MessageError = "error"
)

const (
	NotSignedIn           = "not_signed_in"
	AlreadySignedIn       = "already_signed_in"
	CannotJoinRunningGame = "cannot_join_running_game"
	TableIsFull           = "table_is_full"
	TokenMismatch         = "token_mismatch"
	UnknownToken          = "unknown_token"
	TooManyGames          = "too_many_games"
	UnknownGame           = "unknown_game"
	AlreadyPlaying        = "already_playing"
	AlreadySelected       = "already_selected"
	NotACharacter         = "not_a_character"
	NotPlaying            = "not_playing"
	//CannotChangeCharacter = "cannot_change_character"
	GameAlreadyStarted   = "game_already_started"
	CharacterNotSelected = "character_not_selected"
	NotYourTurn          = "not_your_turn"
	IllegalState         = "illegal_state"
	IllegalMove          = "illegal_move"
	NotYourCard          = "not_your_card"
	MustShowACard        = "must_show_a_card"
)

type SignInRequest struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type SignInResponse struct {
	Token        string         `json:"token,omitempty"`
	RunningGames []GameSynopsis `json:"running_games,omitempty"`
}

type GameSynopsis struct {
	GameID    string `json:"game_id"`
	Character int    `json:"character"`
	PlayerID  int    `json:"player_id"`
}

type CreateGameRequest struct {
	// empty
}

type CreateGameResponse struct {
	GameID string `json:"game_id"`
}

type JoinGameRequest struct {
	GameID string `json:"game_id"`
}

type JoinGameResponse struct {
	Players []NotifyUserState `json:"players"`
}

type SelectCharacterRequest struct {
	Character int `json:"character"`
}

type VoteStartRequest struct {
	Vote bool `json:"vote"`
}

type RollDicesRequest struct {
	// empty
}

type MoveRequest struct {
	EnterRoom Card `json:"enter_room"`
	MapX      int  `json:"map_x"`
	MapY      int  `json:"map_y"`
}

type QuerySolutionRequest struct {
	Character Card `json:"character"`
	Room      Card `json:"room"`
	Weapon    Card `json:"weapon"`
}

type RevealRequest struct {
	Card Card `json:"card,omitempty"`
}

type DeclareSolutionRequest struct {
	Character Card `json:"character"`
	Room      Card `json:"room"`
	Weapon    Card `json:"weapon"`
}

type PassRequest struct {
	// empty
}

// Message is a frame going from fe to be or vicersa.
type MessageFrame struct {
	Header MessageHeader
	Body   interface{}
}

type MessageHeader struct {
	Type string `json:"type"`
	//	Error string `json:"error,omitempty"`
}

type NotifyError struct {
	Error string `json:"error"`
}

/*
NotifyUserState communicate weather a user is reachable, online is true, or not,
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

type NotifyGameStarted struct {
	Deck         []Card `json:"deck"`
	PlayersOrder []int  `json:"players_order"`
}

type NotifyGameState struct {
	State         State `json:"state"`
	CurrentPlayer int   `json:"current_player"`

	Dice1 int `json:"dice1,omitempty"`
	Dice2 int `json:"dice2,omitempty"`

	Room Card `json:"room,omitempty"`
	MapX int  `json:"map_x,omitempty"`
	MapY int  `json:"map_y,omitempty"`

	AnsweringPlayer int  `json:"answering_player,omitempty"`
	Character       Card `json:"character,omitempty"`
	//Room Card
	Weapon  Card `json:"weapon,omitempty"`
	Matched bool `json:"matched,omitempty"`
}
