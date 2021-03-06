package vterm

import (
	"log"

	"github.com/aaronjanse/3mux/ecma48"
)

func (v *VTerm) handleEraseInDisplay(directive int) {
	switch directive {
	case 0: // clear from Cursor to end of screen
		for i := v.Cursor.X; i < len(v.Screen[v.Cursor.Y]); i++ {
			v.Screen[v.Cursor.Y][i] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
		}
		if v.Cursor.Y+1 < len(v.Screen) {
			for j := v.Cursor.Y + 1; j < len(v.Screen); j++ {
				for i := 0; i < len(v.Screen[j]); i++ {
					v.Screen[j][i] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
				}
			}
		}
		v.RedrawWindow()
	case 1: // clear from Cursor to beginning of screen
		for j := 0; j < v.Cursor.Y; j++ {
			for i := 0; i < len(v.Screen[j]); i++ {
				v.Screen[j][i] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
			}
		}
		v.RedrawWindow()
	case 2: // clear entire screen (and move Cursor to top left?)
		for i := 0; i < v.h; i++ {
			if i >= len(v.Screen) {
				newLine := make([]ecma48.StyledChar, v.w)
				for x := range newLine {
					newLine[x].Style = v.Cursor.Style
				}
				v.Screen = append(v.Screen, newLine)
			}
			for j := range v.Screen[i] {
				v.Screen[i][j] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
			}
		}
		v.setCursorPos(0, 0)
		v.RedrawWindow()
	case 3: // clear entire screen and delete all lines saved in scrollback buffer
		v.Scrollback = [][]ecma48.StyledChar{}
		for i := range v.Screen {
			for j := range v.Screen[i] {
				v.Screen[i][j] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
			}
		}
		v.setCursorPos(0, 0)
		v.RedrawWindow()
	default:
		log.Printf("Unrecognized erase in display directive: %d", directive)
	}
}

func (v *VTerm) handleEraseInLine(directive int) {
	switch directive {
	case 0: // clear from Cursor to end of line
		for i := v.Cursor.X; i < len(v.Screen[v.Cursor.Y]); i++ {
			v.Screen[v.Cursor.Y][i] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
		}
	case 1: // clear from Cursor to beginning of line
		for i := 0; i < v.Cursor.X; i++ {
			v.Screen[v.Cursor.Y][i] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
		}
	case 2: // clear entire line; Cursor position remains the same
		for i := 0; i < len(v.Screen[v.Cursor.Y]); i++ {
			v.Screen[v.Cursor.Y][i] = ecma48.StyledChar{Rune: ' ', Style: v.Cursor.Style}
		}
	default:
		log.Printf("Unrecognized erase in line directive: %d", directive)
	}
	v.RedrawWindow()
}
