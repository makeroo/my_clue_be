package game

// Card is a card in the solution deck that comprises of all the characters,
// rooms and weapons. I prefer this solution instead of having three enums,
// one for each kind of object and a Deck of "union".
type Card int

const (
	// NoCard is used to signal the absence of a card.
	NoCard Card = iota
	// Candlestick is the candlestick card
	Candlestick
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

// IsRoom returns true if the given card is a room.
func IsRoom(card Card) bool {
	return Kitchen <= card && card <= Study
}

// IsWeapon returns true if the given card is a weapon.
func IsWeapon(card Card) bool {
	return Candlestick <= card && card <= Wrenck
}

// IsCharacter returns true if the given card is a character.
func IsCharacter(card Card) bool {
	return MissScarlett <= card && card <= MrsWhite
}

// IsCard returns true if the card is valid.
// Used when casting from int (json).
func IsCard(card Card) bool {
	return card >= Candlestick && card <= MrsWhite
}
