package types

import tea "github.com/charmbracelet/bubbletea"

type Tool interface {
	Update(msg tea.Msg) (Tool, tea.Cmd)
	View() string
}

type ToolItem struct {
	Tool        func() Tool
	Name        string
	Description string
}
