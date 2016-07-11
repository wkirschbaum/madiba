package main

import "github.com/nsf/termbox-go"

var row = 0
var col = 0
var buffer = [][]rune{}

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
}

func typeBackspace() {
	if col != 0 {
		buffer[row] = buffer[row][:len(buffer[row])-1]
		col--
	} else if row > 0 {
		row--
		col = len(buffer[row])
	}
}

func typeLetter(letter rune) {
	buffer[row] = append(buffer[row], letter)
	col++
}

func typeEnter() {
	buffer = append(buffer, []rune{})
	row++
	col = 0
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
			default:
				typeLetter(ev.Ch)
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw()
	}
}
