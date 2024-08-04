package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Chat struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
}

func InitChat() Chat {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."

	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false
	vp := viewport.New(30, 5)

	vp.SetContent(`Welcome to the chat room! Type to send a message.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Chat{
		textarea:    ta,
		messages:    make([]string, 0),
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
	}
}

func (c Chat) Init() tea.Cmd {
	return textarea.Blink
}

func (s Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	s.textarea, tiCmd = s.textarea.Update(msg)
	s.viewport, vpCmd = s.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(s.textarea.Value())
			return s, tea.Quit
		case tea.KeyEnter:
			s.messages = append(s.messages, s.senderStyle.Render("You: ")+s.textarea.Value())
			s.viewport.SetContent(strings.Join(s.messages, "\n"))
			s.textarea.Reset()
			s.viewport.GotoBottom()
		}
	}

	return s, tea.Batch(tiCmd, vpCmd)
}

func (s Chat) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		s.viewport.View(),
		s.textarea.View(),
	) + "\n\n"
}
