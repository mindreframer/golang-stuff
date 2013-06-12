package main

import (
	"bytes"
	"fmt"
)

type field struct {
	squares              []int
	width, height        int
	numberOfSquaresMowed int
}

func NewField(fieldWidth, fieldHeight int) *field {
	field := field{
		squares: make([]int, fieldHeight*fieldWidth, fieldHeight*fieldWidth),
		width:   fieldWidth,
		height:  fieldHeight,
	}
	return &field
}

func (f *field) cut(x, y int) {
	index := y*f.width + x

	if f.squares[index] != 0 {
		return
	}

	f.numberOfSquaresMowed++
	f.squares[index] = f.numberOfSquaresMowed
}

func (f *field) allMowed() bool {
	return f.numberOfSquaresMowed == f.width*f.height
}

func (f *field) toString(startX, startY int) string {
	text := bytes.NewBuffer(make([]byte, 0, 5*len(f.squares)))

	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			step := f.squares[y*f.width+x]
			text.WriteString(fmt.Sprint(step))
			if y == startY && x == startX {
				text.WriteString("*")
			} else {
				text.WriteString(" ")
			}
			if step < 10 {
				text.WriteString(" ")
			}
			if step < 100 {
				text.WriteString(" ")
			}
		}
		text.WriteString("\n\n")
	}
	text.WriteString("\n* - starting location")
	return text.String()
}
