package sudoku

import (
	"github.com/dropbox/godropbox/container/bitvector"
)

type puzzle bitvector.BitVector

var allOnes []byte

func (p *puzzle) Duplicate() puzzle {
	bv := bitvector.BitVector(*p)
	backing := make([]byte, len(bv.Bytes()))
	copy(backing, bv.Bytes())
	duplicate := bitvector.NewBitVector(backing, bv.Length())
	return puzzle(*duplicate)
}

func BlankPuzzle() puzzle {
	if allOnes == nil {
		allOnes = make([]byte, 92)
		for i := range allOnes {
			allOnes[i] = byte(0b11111111)
		}
	}
	backing := make([]byte, 92)
	copy(backing, allOnes)
	p := bitvector.NewBitVector(backing, 81*9)
	return puzzle(*p)
}

func (p *puzzle) IsSet(i, j int) bool {
	position := (i * 9) + j
	bv := bitvector.BitVector(*p)
	return bv.Element(position) == 1
}

func (p *puzzle) Unset(i, j int) {
	position := (i * 9) + j
	bv := bitvector.BitVector(*p)
	bv.Set(0, position)
}

func (p *puzzle) Length(i int) uint {
	start := i * 9
	bv := bitvector.BitVector(*p)
	var l uint = 0
	for i := 0; i <= 8; i++ {
		if bv.Element(start+i) == 1 {
			l++
		}
	}
	return l
}
