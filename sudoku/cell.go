package sudoku

import (
	"github.com/dropbox/godropbox/container/bitvector"
)

type Cell struct {
	bv     bitvector.BitVector
	length uint
}

func (c *Cell) Unset(i int) {
	c.bv.Set(0, i)
	c.length -= 1
}

func (c *Cell) IsSet(i int) bool {
	return c.bv.Element(i) == 1
}

func (c *Cell) Length() uint {
	return c.length
}

func (c *Cell) Duplicate() *Cell {
	duplicate := NewCell()
	duplicate.length = c.length
	copy(duplicate.bv.Bytes(), c.bv.Bytes())
	return duplicate
}

func NewCell() *Cell {
	return &Cell{
		*bitvector.NewBitVector([]byte{byte(0b11111111), byte(0b11111111)}, 9),
		9,
	}
}
