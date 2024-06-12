package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	text := []string{""} // Text buffer, each entry is a line
	cursorX, cursorY := 0, 0

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true).SetTitle("Go Text Editor")

	// Display cursor
	updateTextView := func() {
		textView.Clear()
		for i, line := range text {
			if i == cursorY {
				textView.Write([]byte(line[:cursorX] + "|" + line[cursorX:]))
			} else {
				textView.Write([]byte(line))
			}
			textView.Write([]byte("\n"))
		}
	}

	// Allow editing text
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyCtrlS:
			// Handle save logic here (e.g., write to file)
			return nil
		case tcell.KeyUp:
			if cursorY > 0 {
				cursorY--
				if cursorX > len(text[cursorY]) {
					cursorX = len(text[cursorY])
				}
			}
		case tcell.KeyDown:
			if cursorY < len(text)-1 {
				cursorY++
				if cursorX > len(text[cursorY]) {
					cursorX = len(text[cursorY])
				}
			}
		case tcell.KeyLeft:
			if cursorX > 0 {
				cursorX--
			} else if cursorY > 0 {
				cursorY--
				cursorX = len(text[cursorY])
			}
		case tcell.KeyRight:
			if cursorX < len(text[cursorY]) {
				cursorX++
			} else if cursorY < len(text)-1 {
				cursorY++
				cursorX = 0
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if cursorX > 0 {
				line := text[cursorY]
				text[cursorY] = line[:cursorX-1] + line[cursorX:]
				cursorX--
			} else if cursorY > 0 {
				prevLine := text[cursorY-1]
				text = append(text[:cursorY], text[cursorY+1:]...)
				cursorY--
				cursorX = len(prevLine)
				text[cursorY] = prevLine + text[cursorY]
			}
		case tcell.KeyEnter:
			line := text[cursorY]
			text[cursorY] = line[:cursorX]
			text = append(text[:cursorY+1], append([]string{line[cursorX:]}, text[cursorY+1:]...)...)
			cursorY++
			cursorX = 0
		default:
			if event.Key() == tcell.KeyRune {
				line := text[cursorY]
				text[cursorY] = line[:cursorX] + string(event.Rune()) + line[cursorX:]
				cursorX++
			}
		}
		updateTextView()
		return event
	})

	updateTextView()

	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}
