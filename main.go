package main

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
)

func main() {
	app := tview.NewApplication()
	commandBuffer := "" // Buffer to store the current command line input
	outputBuffer := ""  // Buffer to store the output of executed commands
	cursorX, cursorY := 0, 0

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true).SetTitle("Go Terminal")

	// Update the text view with the current buffers and cursor position
	updateTextView := func() {
		textView.Clear()
		fmt.Fprintf(textView, "%s\n$ %s", outputBuffer, commandBuffer)
		textView.Highlight(fmt.Sprintf("cursor_%d_%d", cursorY, cursorX))
	}

	// Execute a command and capture the output
	executeCommand := func(command string) {
		cmd := exec.Command("cmd", "/c", command)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			outputBuffer += fmt.Sprintf("error: %v\n", err)
		}
		outputBuffer += out.String()
	}

	// Handle input for command line
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyEnter:
			executeCommand(commandBuffer)
			commandBuffer = ""
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if len(commandBuffer) > 0 {
				commandBuffer = commandBuffer[:len(commandBuffer)-1]
			}
		default:
			if event.Key() == tcell.KeyRune {
				commandBuffer += string(event.Rune())
			}
		}
		updateTextView()
		return event
	})

	// Handle mouse events for cursor movement
	textView.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick {
			_, cy := event.Position()
			cursorY = cy
			updateTextView()
		}
		return action, event
	})

	updateTextView()

	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}
