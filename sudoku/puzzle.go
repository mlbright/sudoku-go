package sudoku

type puzzle []Cell

var blank puzzle

func (p *puzzle) Duplicate() puzzle {
	duplicate := make([]Cell, 81)
	for i := 0; i < 81; i++ {
		duplicate[i] = *(*p)[i].Duplicate()
	}
	return duplicate
}

func BlankPuzzle() puzzle {
	if blank == nil {
		p := make([]Cell, 81)
		for i := range p {
			p[i] = *NewCell()
		}
		blank = puzzle(p)
	}
	return blank.Duplicate()
}

func (p *puzzle) IsSet(i, j int) bool {
	return (*p)[i].IsSet(j)
}

func (p *puzzle) Unset(i, j int) {
	(*p)[i].Unset(j)
}

func (p *puzzle) Length(i int) uint {
	return (*p)[i].Length()
}