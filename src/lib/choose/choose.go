package choose

import (
	"github.com/arfrie22/arf-toolkit/src/lib/types"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	tool types.ToolItem
}

func (i item) Title() string       { return i.tool.Name }
func (i item) Description() string { return i.tool.Description }
func (i item) FilterValue() string { return i.tool.Name }

type ChooseMsg struct {
	Run func()
}

type Model struct {
	list list.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			var cmd tea.Cmd
			cmd = nil

			i, ok := m.list.SelectedItem().(item)
			if ok {
				cmd = func() tea.Msg {
					return ChooseMsg{Run: i.tool.Run}
				}
			}

			return m, cmd
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}

type chooseKeyMap struct {
	selectItem key.Binding
}

func newChooseKeyMap() *chooseKeyMap {
	return &chooseKeyMap{
		selectItem: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select tool"),
		),
	}
}

func (m *Model) AddTool(tool types.ToolItem) {
	m.list.InsertItem(len(m.list.Items()), item{tool: tool})
}

func Choose() Model {
	var keyMap = newChooseKeyMap()
	items := []list.Item{}

	m := Model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "ARF Toolkit"
	m.list.SetStatusBarItemName("tool", "tools")
	m.list.InfiniteScrolling = true
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.selectItem,
		}
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keyMap.selectItem,
		}
	}

	return m
}
