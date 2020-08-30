package game

// Error is an alias for string to be able to declare errors as constants.
// See https://dave.cheney.net/2016/04/07/constant-errors
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// NotSignedIn error: issued for game related requests.
	NotSignedIn = Error("not_signed_in")
	// CannotJoinRunningGame error: illegal join game request.
	CannotJoinRunningGame = Error("cannot_join_running_game")
	// TableIsFull error: cannot join.
	TableIsFull = Error("table_is_full")
	// TokenMismatch error: the user provided a different token.
	TokenMismatch = Error("token_mismatch")
	// UnknownToken error: the provided token is unknown (maybe expired?).
	UnknownToken = Error("unknown_token")
	// TooManyGames error: per user game limit exceeded.
	TooManyGames = Error("too_many_games")
	// UnknownGame error: illegal join request.
	UnknownGame = Error("unknown_game")
	// AlreadyPlaying error: one tab per game limit.
	AlreadyPlaying = Error("already_playing")
	// AlreadySelected error: the choosen character has already been selected.
	AlreadySelected = Error("already_selected")
	// NotACharacter error: illegal select char request.
	NotACharacter = Error("not_a_character")
	// NotPlaying error: illegal game related request.
	NotPlaying = Error("not_playing")
	// GameAlreadyStarted error: illegal game setup request (eg. vote start / select char).
	GameAlreadyStarted = Error("game_already_started")
	// CharacterNotSelected error: cannot vote start before having selected a character.
	CharacterNotSelected = Error("character_not_selected")
	// NotYourTurn error: illegal game related request.
	NotYourTurn = Error("not_your_turn")
	// IllegalState error: illegal game related request.
	IllegalState = Error("illegal_state")
	// IllegalMove error: illegal move parameters.
	IllegalMove = Error("illegal_move")
	// NotYourCard error: cannot reveal a card not in your deck.
	NotYourCard = Error("not_your_card")
	// MustShowACard error: cannot pass if you have a card to show.
	MustShowACard = Error("must_show_a_card")
	// NotInARoom error: cannot query solution if you are not in a room.
	NotInARoom = Error("not_in_a_room")
)
