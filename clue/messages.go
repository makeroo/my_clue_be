package clue

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
)

// Message is a frame going from fe to be or vicersa.
type Message struct {
	SignIn *struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	} `json:"sign_in,omitempty"`

	SignInResponse *SignInResponse `json:"sign_in_response,omitempty"`

	CreateGame *struct {
		// empty
	} `json:"create_game,omitempty"`

	CreateGameResponse *CreateGameResponse `json:"create_game_resp,omitempty"`

	JoinGame *struct {
		GameID string `json:"game_id"`
	} `json:"join_game,omitempty"`

	SelectCharacter *struct {
		Character int `json:"character"`
	} `json:"select_char,omitempty"`

	VoteStart *struct {
		Vote bool `json:"vote"`
	} `json:"vote_start,omitempty"`

	RollDices *struct {
		// empty
	} `json:"roll_dices,omitempty"`

	RollDicesResponse *RollDicesResonse `json:"roll_dices_response,omitempty"`

	NotifyUserState *NotifyUserState `json:"notify_user_state,omitempty"`

	NotifyGameStarted *NotifyGameStarted `json:"notify_game_started,omitempty"`

	NotifyGameState *NotifyGameState `json:"notify_game_state,omitempty"`

	Error string `json:"error,omitempty"`
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

type CreateGameResponse struct {
	GameID string `json:"game_id"`
}

type RollDicesResonse struct {
	Dice1 int `json:"dice1"`
	Dice2 int `json:"dice2"`
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
	Name      string `json:"name"`
	Character int    `json:"character"`
	Online    bool   `json:"online"`
}

type NotifyGameStarted struct {
	Deck         []Card `json:"deck"`
	PlayersOrder []int  `json:"players_order"`
}

type NotifyGameState struct {
	State         State `json:"state"`
	CurrentPlayer int   `json:"current_player"`
}

/*
welcome page

    join a game
        text field: enter game id
        send req: join(gameId) -> game

    create a game
        send req: create() -> game

game:
    starting
        selectChar(char, name) -> charId / alreadyTaken

        -> notifyChar(char, name)

        voteStart()

        -> notifyStart(name)

    turn(char)

        -> startTurn

        rollDices()

        -> notifyDices(dices)

        (LATER card)

        selectAction(still, secretPassage, move)

        move(cell)

        -> notifyMove()

        (if in room)

        ask( (implicit room), char, weap)

        -> notifyAsk(room, char, weap, aswer:bool)

        answer( room or char or weap )

        -> notifyAnswer( room or char or weap )
        -> notifyBlindAnswer

*/
