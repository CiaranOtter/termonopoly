package player

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type PlayerList struct {
	list     list.Model
	choice   string
	quitting bool
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (pl PlayerList) Init() tea.Cmd {
	return nil
}

func (pl PlayerList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		pl.list.SetWidth(msg.Width)
		return pl, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			pl.quitting = true
			return pl, tea.Quit

		case "enter":
			i, ok := pl.list.SelectedItem().(item)
			if ok {
				pl.choice = string(i)
			}
			return pl, tea.Quit
		}
	}

	pl.list, cmd = pl.list.Update(msg)
	cmds = append(cmds, cmd)
	return pl, tea.Batch(cmds...)
}

func (pl PlayerList) View() string {
	if pl.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", pl.choice))
	}
	if pl.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + pl.list.View()
}

func NewPlayerlist() PlayerList {
	items := []list.Item{
		item("ciaran"),
		item("dave"),
		item("louise"),
		item("andrew"),
		item("that_other_guy"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)

	l.Title = "Players in the game"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := PlayerList{list: l}

	return m
}
