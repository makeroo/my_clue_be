package game

const (
	xx = -1
	oo = 0
	ki = int(Kitchen)
	ba = int(Ballroom)
	co = int(Conservatory)
	di = int(DiningRoom)
	bi = int(BilliardRoom)
	li = int(Library)
	lo = int(Lounge)
	ha = int(Hall)
	st = int(Study)
)

var clueBoard = [25][24]int{
	//0   1   2   3   4   5   6   7   8   9  10  11  12  13  14  15  16  17  18  19  20  21  22  23
	{xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx}, // 00
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 01
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx}, // 02
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx}, // 03
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx}, // 04
	{xx, xx, xx, xx, xx, xx, oo, ba, xx, xx, xx, xx, xx, xx, xx, xx, ba, oo, xx, xx, xx, xx, xx, xx}, // 05
	{xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, co, oo, oo, oo, oo, xx}, // 06
	{oo, oo, oo, oo, ki, oo, oo, oo, xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, oo, oo, xx}, // 07
	{xx, oo, oo, oo, oo, oo, oo, oo, oo, ba, oo, oo, oo, oo, ba, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 08
	{xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, bi, xx, xx, xx, xx, xx, xx}, // 09
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 10
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 11
	{xx, xx, xx, xx, xx, xx, xx, xx, di, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 12

	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, li, oo, bi, xx}, // 13
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 14
	{xx, xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 15
	{xx, oo, oo, oo, oo, oo, di, oo, oo, oo, xx, xx, xx, xx, xx, oo, li, xx, xx, xx, xx, xx, xx, xx}, // 16
	{xx, oo, oo, oo, oo, oo, oo, oo, oo, oo, oo, ha, ha, oo, oo, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 17
	{xx, oo, oo, oo, oo, oo, lo, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, oo, xx, xx, xx, xx, xx, xx}, // 18
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, oo, oo, oo, oo, oo, oo, xx}, // 19
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, st, oo, oo, oo, oo, oo, xx}, // 20
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 21
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 22
	{xx, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, oo, oo, xx, xx, xx, xx, xx, xx, xx}, // 23
	{xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, oo, xx, xx, xx, xx, xx, xx, xx}, // 24
}

var initialPositions = map[Card]PawnPosition{
	MissScarlett: PositionAt(7, 24),
	RevGreen:     PositionAt(14, 0),
	ColMustard:   PositionAt(0, 17),
	ProfPlum:     PositionAt(23, 19),
	MrsPeacock:   PositionAt(23, 7),
	MrsWhite:     PositionAt(9, 0),
}
