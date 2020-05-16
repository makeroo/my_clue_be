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
	NotYourTurn          = "not_your_turn"
	IllegalState         = "illegal_state"
	IllegalMove          = "illegal_move"
	NotYourCard          = "not_your_card"
	MustShowACard        = "must_show_a_card"
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

	JoinGameResponse *JoinGameResponse `json:"join_game_resp,omitempty"`

	SelectCharacter *struct {
		Character int `json:"character"`
	} `json:"select_char,omitempty"`

	VoteStart *struct {
		Vote bool `json:"vote"`
	} `json:"vote_start,omitempty"`

	RollDices *struct {
		// empty
	} `json:"roll_dices,omitempty"`

	Move *struct {
		EnterRoom Card `json:"enter_room"`
		MapX      int  `json:"map_x"`
		MapY      int  `json:"map_y"`
	} `json:"move,omitempty"`

	QuerySolution *struct {
		Character Card `json:"character"`
		Room      Card `json:"room"`
		Weapon    Card `json:"weapon"`
	} `json:"query_solution,omitempty"`

	Reveal *struct {
		card Card `json:"card,omitempty"`
	} `json:"reveal,omitempty"`

	DeclareSolution *struct {
		Character Card `json:"character"`
		Room      Card `json:"room"`
		Weapon    Card `json:"weapon"`
	} `json:"declare_solution,omitempty"`

	Pass *struct {
		// empty
	} `json:"pass,omitempty"`

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

type JoinGameResponse struct {
	Players []NotifyUserState `json:"players"`
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
