package tools

import (
	"fmt"
	"io"
	"strings"

	"github.com/arfrie22/arf-toolkit/src/lib/types"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "golang.org/x/sys/windows/registry"
)

var (
	name string = "Outlook Preview Tool"
	desc string = "Allows you to set the previewer used by Outlook for any file type."
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type previewer struct {
	name string
	id   string
}

func (p previewer) FilterValue() string {
	return p.name
}

type previewerDelegate struct{}

func (d previewerDelegate) Height() int                             { return 1 }
func (d previewerDelegate) Spacing() int                            { return 0 }
func (d previewerDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d previewerDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(previewer)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	stage         int
	fileExtInput  textinput.Model
	previewerList list.Model
}

func new() types.Tool {
	ti := textinput.New()
	ti.Placeholder = "File Extension"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 32

	previewers := []list.Item{}
	// k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\PreviewHandlers`, registry.QUERY_VALUE)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer k.Close()

	// previewerIds, err := k.ReadValueNames(-1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, id := range previewerIds {
	// 	p, _, err := k.GetStringValue(id)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	previewers = append(previewers, previewer{name: p, id: id})
	// }
	previewers = append(previewers, previewer{name: "test", id: "test"})
	previewers = append(previewers, previewer{name: "test2", id: "test2"})

	l := list.New(previewers, previewerDelegate{}, 20, 1)
	l.Title = "Previewers"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{
		stage:         0,
		fileExtInput:  ti,
		previewerList: l,
	}
}

func OutlookPreview() types.ToolItem {
	return types.ToolItem{
		Tool:        new,
		Name:        name,
		Description: desc,
	}
}

func (m model) Update(msg tea.Msg) (types.Tool, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			switch m.stage {
			case 0:
				v := m.fileExtInput.Value()
				if v != "" {
					m.stage = 1
				}
			case 1:
				i, ok := m.previewerList.SelectedItem().(previewer)
				if ok {
					fmt.Println(i)
					return m, tea.Quit
				}
				// if ok {
				// 	k, err := registry.OpenKey(registry.CLASSES_ROOT, `.`+m.fileExtInput.Value(), registry.CREATE_SUB_KEY)
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				// 	defer k.Close()

				// 	k, _, err = registry.CreateKey(k, `shellex`, registry.CREATE_SUB_KEY)
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				// 	defer k.Close()

				// 	k, _, err = registry.CreateKey(k, `{8895b1c6-b41f-4c1c-a562-0d564250836f}`, registry.CREATE_SUB_KEY)
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}
				// 	defer k.Close()

				// 	err = k.SetStringValue("", i.id)
				// 	if err != nil {
				// 		log.Fatal(err)
				// 	}

				// 	return m, tea.Quit
				// }
			}
		}

	case tea.WindowSizeMsg:
		m.fileExtInput.Width = msg.Width - 4
		m.previewerList.SetSize(msg.Width-4, msg.Height-4)
	}

	switch m.stage {
	case 0:
		var cmd tea.Cmd
		m.fileExtInput, cmd = m.fileExtInput.Update(msg)
		return m, cmd
	case 1:
		var cmd tea.Cmd
		m.previewerList, cmd = m.previewerList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.stage {
	case 0:
		return m.fileExtInput.View()
	case 1:
		return m.previewerList.View()
	}

	return ""
}
