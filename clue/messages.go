package clue

const (
	NotSignedIn           = "not_signed_in"
	AlreadySignedIn       = "already_signed_in"
	CannotJoinRunningGame = "cannot_join_running_game"
	TableIsFull           = "table_is_full"
	TokenMismatch         = "token_mismatch"
	UnknownToken          = "unknown_token"
	TooManyGames          = "too_many_games"
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

	NotifyUserState *NotifyUserState `json:"notify_user_online,omitempty"`

	Error string `json:"error,omitempty"`
}

type SignInResponse struct {
	Token        string         `json:"token,omitempty"`
	RunningGames []GameSynopsis `json:"running_games,omitempty"`
}

type GameSynopsis struct {
	GameID    string `json:"game_id"`
	Character int    `json:"character"`
}

type CreateGameResponse struct {
	GameID string `json:"game_id"`
}

type NotifyUserState struct {
	OldName   string `json:"old_name"`
	NewName   string `json:"new_name"`
	Character int    `json:"character"`
	Online    bool   `json:"online"`
}

/* TODO type NotifyNameChange struct {

}*/

/*
type JoinGameRequest struct {
	GameID string `json:"game_id"`
	Name   string `json:"name"`
}

type JoinGameResponse struct {
	PlayerToken string `json:"token"`
}

// type CreateGameRequest empty


type SelectCharacterRequest struct {
	Character int    `json:"character"`
	Name      string `json:"name"`
}

type SelectCharacterResponse struct {
	PlayerToken  string `json:"player_token,omitempty"`
	AlreadyTaken bool   `json:"already_taken"`
}

type NotifyNewPlayer struct {
	Character int    `json:"character"`
	Name      string `json:"name"`
}
*/
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
