package sudoku

import (
	"github.com/dropbox/godropbox/container/bitvector"
)

type Cell struct {
	data   bitvector.BitVector
	length int
}

func (c *Cell) Unset(i int) {
	c.data.Set(0, i)
	c.length -= 1
}

func (c *Cell) IsSet(i int) bool {
	return c.data.Element(i) == 1
}

func (c *Cell) Length() int {
	return c.length
}

func (c *Cell) Duplicate() *Cell {
	backing := make([]byte, 2)
	copy(backing, c.data.Bytes())
	return &Cell{
		*bitvector.NewBitVector(backing, 9),
		c.length,
	}
}

func NewCell() *Cell {
	return &Cell{
		*bitvector.NewBitVector([]byte{byte(0b11111111), byte(0b11111111)}, 9),
		9,
	}
}
