//go:generate go-winres make --product-version=git-tag --file-version=git-tag
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/arfrie22/arf-toolkit/lib/choose"
	"github.com/arfrie22/arf-toolkit/tools"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	nextTool func()
	chooser  choose.Model
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		}
	case choose.ChooseMsg:
		m.nextTool = msg.Run
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.chooser, cmd = m.chooser.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return docStyle.Render(m.chooser.View())
}

func main() {
	m := model{chooser: choose.Choose()}
	m.chooser.AddTool(tools.OutlookPreview())

	p := tea.NewProgram(&m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	} else if m.nextTool != nil {
		m.nextTool()
	}

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
