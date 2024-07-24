package main

import (
	"bufio"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type session struct {
	index    int
	name     string
	attached bool
}

type ByIndex []session

func (m ByIndex) Len() int           { return len(m) }
func (m ByIndex) Less(i, j int) bool { return m[i].index < m[j].index }
func (m ByIndex) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func initUI() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}
	return screen, nil
}

func closeUI(screen tcell.Screen) {
	screen.Fini()
}

func drawWindows(screen tcell.Screen, windows []string, selected int) {
	screen.Clear()

	cols, _ := screen.Size()
	for i, window := range windows {
		fg := tcell.ColorLightGray
		bg := tcell.ColorBlack
		if i == selected {
			fg = tcell.ColorBlack
			bg = tcell.ColorLightGray
		}
		indexStr := "(" + strconv.Itoa(i) + ") " + window

		for j, ch := range indexStr {
			if j < cols {
				screen.SetContent(j, i, ch, nil, tcell.StyleDefault.Foreground(fg).Background(bg))
			}
		}

		for j := len(indexStr); j < cols; j++ {
			screen.SetContent(j, i, ' ', nil, tcell.StyleDefault.Foreground(fg).Background(bg))
		}
	}
	screen.Show()
}

func getTmuxSessions() ([]session, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_id}~#{session_name}~#{session_attached}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sessionInfos []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		sessionInfos = append(sessionInfos, scanner.Text())
	}
	var retDate []session
	for _, sesData := range sessionInfos {
		splitted := strings.Split(sesData, "~")
		index, _ := strconv.Atoi(strings.Replace(splitted[0], "$", "", 1))
		attached := false
		if splitted[2] == "1" {
			attached = true
		}

		retDate = append(retDate, session{index: index, name: splitted[1], attached: attached})
	}

	sort.Sort(ByIndex(retDate))
	return retDate, scanner.Err()
}

func getTmuxWindows() ([]string, error) {
	sessions, err := getTmuxSessions()
	if err != nil {
		return nil, err
	}

	var allWindows []string
	for _, session := range sessions {
		cmd := exec.Command("tmux", "list-windows", "-t", session.name, "-F", "#{window_name}")
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
			commands = append(commands, session.name)
			name := "- " + session.name + " -> " + windows[0]
			if session.attached {
				name = name + " (attached)"
			}
			allWindows = append(allWindows, name)
		} else {
			name := "-" + session.name + ": " + strconv.Itoa(len(windows)) + " windows"
			if session.attached {
				name = name + " (attached)"
			}
			allWindows = append(allWindows, name)
			commands = append(commands, session.name)
			for windowIndex, win := range windows {
				commands = append(commands, session.name+":"+strconv.Itoa(windowIndex+1))
				allWindows = append(allWindows, "--> "+win)
			}
		}
	}

	return allWindows, nil
}

var commands []string

func main() {
	screen, err := initUI()
	if err != nil {
		panic(err)
	}
	defer closeUI(screen)

	var rows []string
	windows, err := getTmuxWindows()
	if err != nil {
		panic(err)
	}
	rows = append(rows, windows...)

	selected := 0
	for {
		drawWindows(screen, rows, selected)

		event := screen.PollEvent()
		switch event.(type) {
		case *tcell.EventKey:

			key := event.(*tcell.EventKey)
			switch key.Rune() {
			case 'q':
				return
			case 'n':
				//TODO: Open tmuxinator or zoxide session creation popup
			case 'j':
				selected++
				if selected >= len(rows) {
					selected = 0
				}
			case 'k':
				selected--
				if selected < 0 {
					selected = len(rows) - 1
				}
			case 'g':
				selected = 0
			case 'G':
				selected = len(rows) - 1
			}
			switch key.Key() {
			case tcell.KeyEsc, tcell.KeyCtrlC:
				return
			case tcell.KeyDown:
				selected++
				if selected >= len(rows) {
					selected = 0
				}
			case tcell.KeyUp:
				selected--
				if selected < 0 {
					selected = len(rows) - 1
				}
			case tcell.KeyCtrlK:
				selected--
				if selected < 0 {
					selected = len(rows) - 1
				}
			case tcell.KeyHome:
				selected = 0
			case tcell.KeyEnd:
				selected = len(rows) - 1
			case tcell.KeyEnter:
				cmd := exec.Command("tmux", "switch", "-t", commands[selected])
				_, err := cmd.Output()
				if err != nil {
					panic(err)
				}
				return
			}
		}
	}
}
