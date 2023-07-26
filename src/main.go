package main

import (
	"fmt"
	"os"

	"github.com/arfrie22/arf-toolkit/src/lib/choose"
	types "github.com/arfrie22/arf-toolkit/src/lib/types"
	"github.com/arfrie22/arf-toolkit/src/tools"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	activeTool    types.Tool
	chooser       choose.Model
	height, width int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		}
	case choose.ChooseMsg:
		m.activeTool = msg.Tool
		return m, func() tea.Msg {
			return tea.WindowSizeMsg{Height: m.height, Width: m.width}
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	var cmd tea.Cmd
	if m.activeTool == nil {
		m.chooser, cmd = m.chooser.Update(msg)
	} else {
		m.activeTool, cmd = m.activeTool.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.activeTool == nil {
		return docStyle.Render(m.chooser.View())
	} else {
		return docStyle.Render(m.activeTool.View())
	}
}

func main() {
	m := model{chooser: choose.Choose()}
	m.chooser.AddTool(tools.OutlookPreview())

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
