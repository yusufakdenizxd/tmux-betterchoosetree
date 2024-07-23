package main

import (
	"bufio"

	"os/exec"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

func initUI() error {
	return termbox.Init()
}

func closeUI() {
	termbox.Close()
}

func drawWindows(windows []string, selected int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	cols, _ := termbox.Size()
	for i, window := range windows {
		fg := termbox.ColorLightGray
		bg := termbox.ColorDefault
		if i == selected {
			fg = termbox.ColorBlack
			bg = termbox.ColorLightGray
		}
		indexStr := "(" + strconv.Itoa(i) + ") " + window

		for j, ch := range indexStr {
			if j < cols {
				termbox.SetCell(j, i, ch, fg, bg)
			}
		}

		for j := len(indexStr); j < cols; j++ {
			termbox.SetCell(j, i, ' ', fg, bg)
		}
	}
	termbox.Flush()
}

func getTmuxSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sessions []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		sessions = append(sessions, scanner.Text())
	}
	return sessions, scanner.Err()
}

func getTmuxWindows() ([]string, error) {
	sessions, err := getTmuxSessions()
	if err != nil {
		return nil, err
	}

	var allWindows []string
	for _, session := range sessions {
		cmd := exec.Command("tmux", "list-windows", "-t", session, "-F", "#W")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		var windows []string
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			windows = append(windows, scanner.Text())
		}

		if len(windows) == 1 {
			allWindows = append(allWindows, "- "+session+" - "+windows[0])
		} else {
			allWindows = append(allWindows, "- "+session+": "+strconv.Itoa(len(windows))+" windows")
			for _, win := range windows {

				allWindows = append(allWindows, "--> "+win)
			}
		}
	}

	return allWindows, nil
}

func main() {

	err := initUI()
	if err != nil {
		panic(err)
	}
	defer closeUI()

	var rows []string
	windows, err := getTmuxWindows()
	if err != nil {
		panic(err)
	}
	rows = append(rows, windows...)

	selected := 0
	for {
		drawWindows(rows, selected)

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 'q' || ev.Key == termbox.KeyEsc {
				return
			}
			if ev.Ch == 'j' {
				selected++
				if selected >= len(rows) {
					selected = 0
				}
			} else if ev.Ch == 'k' {
				selected--
				if selected < 0 {
					selected = len(rows) - 1
				}
			} else if ev.Ch == 'g' {
				selected = 0
			} else if ev.Ch == 'G' {
				selected = len(rows) - 1
			} else if ev.Key == termbox.KeyEnter {

				cmd := exec.Command("tmux", "switch", "-t", "tmux-betterchoosetree"+":"+strconv.Itoa(selected))
				_, err := cmd.Output()

				if err != nil {
					panic(err)
				}
				return
				// exec.Command("tmux", "switch-client", "-t", "tmux-betterchoosetree:1")
				// return
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
