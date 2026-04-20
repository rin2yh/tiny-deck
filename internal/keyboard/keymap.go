package keyboard

// Physical key layout (4 cols × 3 rows, column-major index):
//
//   col=0 col=1 col=2 col=3
//     0     3     6     9    row=0 (top)
//     1     4     7    10    row=1
//     2     5     8    11    row=2 (bottom)

const (
	KeyCols  = 4
	KeyRows  = 3
	KeyCount = KeyCols * KeyRows
)

type Key uint8

const (
	KeyC0R0 Key = 0
	KeyC0R1 Key = 1
	KeyC0R2 Key = 2
	KeyC1R0 Key = 3
	KeyC1R1 Key = 4
	KeyC1R2 Key = 5
	KeyC2R0 Key = 6
	KeyC2R1 Key = 7
	KeyC2R2 Key = 8
	KeyC3R0 Key = 9
	KeyC3R1 Key = 10
	KeyC3R2 Key = 11
)
