package game

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	viewport viewport.Model
}

var (
	board_style lipgloss.Style = lipgloss.NewStyle().Background(lipgloss.Color("#c7f9cc"))
)

func NewBoard() Board {
	return Board{
		viewport: viewport.New(10, 10),
	}
}

func (b Board) Init() tea.Cmd {
	return b.viewport.Init()
}

func (b Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		dim := msg.(tea.WindowSizeMsg)
		b.viewport = viewport.New(50, dim.Height)
		b.viewport.Init()
		b.viewport.SetContent(board_style.Render("test"))
	}
	var cmd tea.Cmd
	b.viewport, cmd = b.viewport.Update(msg)

	return b, cmd
}

func (b Board) View() string {
	return b.viewport.View()
}
