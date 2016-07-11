package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var row = 0
var col = 0
var buffer = [][]rune{}

func drawFooter(text string) {
	for i, c := range text {
		termbox.SetCell(i, 20, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func redraw() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	defer termbox.Flush()

	for crow, ccols := range buffer {
		for cpos, char := range ccols {
			termbox.SetCell(cpos, crow, char, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	termbox.SetCursor(col, row)

	drawFooter(fmt.Sprintf("pos: %d:%d", row, col))
}

func typeBackspace() {
	if col != 0 {
		rest := buffer[row][col:]
		buffer[row] = buffer[row][:col-1]
		buffer[row] = append(buffer[row], rest...)
		col--
	} else if row > 0 {
		row--
		col = len(buffer[row])
	}
}

func typeLetter(letter rune) {
	var rest = make([]rune, len(buffer[row][col:]))
	copy(rest, buffer[row][col:])
	buffer[row] = append(buffer[row][:col], letter)
	buffer[row] = append(buffer[row], rest...)
	col++
}

func moveLeft() {
	if col != 0 {
		col--
	} else if row > 0 {
		row--
		col = len(buffer[row])
	}
}

func moveRight() {
	if col < len(buffer[row]) {
		col++
	} else if row < len(buffer)-1 {
		col = 0
		row++
	}
}

func typeEnter() {
	var rest = make([]rune, len(buffer[row][col:]))
	copy(rest, buffer[row][col:])
	buffer[row] = buffer[row][:col]

	if row < len(buffer) {
		var rest = make([][]rune, len(buffer[row+1:]))
		copy(rest, buffer[row+1:])
		buffer = append(buffer[:row+1], []rune{})
		buffer = append(buffer, rest...)
	} else {
		buffer = append(buffer, []rune{})
	}

	row++
	col = 0
	buffer[row] = rest
}

func main() {
	buffer = append(buffer, []rune{})
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt)
	redraw()
	exit := false

	for !exit {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC:
				exit = true
			case termbox.KeyEnter:
				typeEnter()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				typeBackspace()
			case termbox.KeyArrowLeft:
				moveLeft()
			case termbox.KeyArrowRight:
				moveRight()
			default:
				typeLetter(ev.Ch)
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw()
	}
}
