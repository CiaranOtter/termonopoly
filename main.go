package main

import (
	"fmt"
	"os"
	"termonopoly/chat"
	"termonopoly/game"
	"termonopoly/player"
	"termonopoly/setup"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Screen struct {
	// tea.Model
	ready      bool
	board      tea.Model
	spinner    spinner.Model
	chat       tea.Model
	playerlist tea.Model
}

// func (s Screen) SideBar(c string) string {
// 	return lipgloss.NewStyle().
// 		Height(s.Viewport.Height).
// 		Width(s.Viewport.Width).
// 		Border(lipgloss.DoubleBorder(), true, false).Render(c)
// }

func NewModel() Screen {
	screen := Screen{
		ready: false,
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	screen.spinner = s

	screen.chat = chat.InitChat()
	screen.playerlist = player.NewPlayerlist()

	screen.board = game.NewBoard()

	return screen
}

func (s Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	s.playerlist, cmd = s.playerlist.Update(msg)
	cmds = append(cmds, cmd)

	s.chat, cmd = s.chat.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	// // Is it a key press?
	case tea.KeyMsg:

		// 	// Cool, what was the actual key pressed?
		switch msg.String() {

		// 	// These keys should exit the program.
		case "ctrl+c", "q":
			return s, tea.Quit

			// 	// The "up" and "k" keys move the cursor up
			// 	case "up", "k":
			// 		if m.cursor > 0 {
			// 			m.cursor--
			// 		}

			// 	// The "down" and "j" keys move the cursor down
			// 	case "down", "j":
			// 		if m.cursor < len(m.choices)-1 {
			// 			m.cursor++
			// 		}

			// 	// The "enter" key and the spacebar (a literal space) toggle
			// 	// the selected state for the item that the cursor is pointing at.
			// 	case "enter", " ":
			// 		_, ok := m.selected[m.cursor]
			// 		if ok {
			// 			delete(m.selected, m.cursor)
			// 		} else {
			// 			m.selected[m.cursor] = struct{}{}
			// 		}
		}

	case tea.WindowSizeMsg:
		if !s.ready {
			// s.viewport = viewport.New(msg.Width, msg.Height)
			// s.viewport.YPosition = headerHeight
			// s.viewport.HighPerformanceRendering = true
			// s.viewport.SetContent(s.body)
			s.ready = true
		}

	default:
		s.spinner, cmd = s.spinner.Update(msg)
		cmds = append(cmds, cmd)
		// return s, cmd
	}

	s.board, cmd = s.board.Update(msg)
	cmds = append(cmds, cmd)
	// Get the actual, physical dimensions of the text block.
	// width := lipgloss.Width(block)
	// height := lipgloss.Height(block)

	// // Here's a shorthand function.
	// w, h := lipgloss.Size(block)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return s, tea.Batch(cmds...)
}

func (s Screen) Init() tea.Cmd {
	return tea.Batch(s.chat.Init(), s.playerlist.Init(), s.board.Init(), tea.ClearScreen, s.spinner.Tick)
}

func (s Screen) View() string {

	if !s.ready {
		return fmt.Sprintf("%s Loading game", s.spinner.View())
	}
	// The header
	// viewport.New(s.Viewport.Height, s.Viewport.Width).

	// Send the UI for rendering
	return lipgloss.JoinHorizontal(0, s.playerlist.View(), s.board.View(), s.chat.View())

	// return style.Render(s.SideBar("test"))
}

func main() {
	// p := tea.NewProgram(chat.InitChat(), tea.WithAltScreen())
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	game := game.NewGame()

	start := setup.ReadCsv("./data/Economics.csv", game)

	pl := player.NewPlayer("Ciaran", start, 1500)

	pl.Space.Print()

	// for {
	// 	roll := rand.Intn(12)
	// 	fmt.Printf("Moved forward %d spaces\n", roll)
	// 	pl.Move(roll, player.FORWARD)
	// 	time.Sleep(1 * time.Second)
	// }

}
