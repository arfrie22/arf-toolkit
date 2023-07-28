package tools

import (
	"fmt"

	"github.com/arfrie22/arf-toolkit/lib/types"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/muesli/termenv"
	"golang.org/x/sys/windows/registry"
)

var (
	name string = "Outlook Preview Tool"
	desc string = "Allows you to set the previewer used by Outlook for any file type."
)

type previewer struct {
	name string
	id   string
}

func (p previewer) FilterValue() string {
	return p.name
}

func run() {
	ti := textinput.New("What file extension would you like to set a previewer for?")
	ti.Placeholder = "e.g. txt"

	extension, err := ti.RunPrompt()
	if err != nil {
		fmt.Println(err)
		return
	}

	previewers := []previewer{}
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\PreviewHandlers`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer k.Close()

	previewerIds, err := k.ReadValueNames(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, id := range previewerIds {
		p, _, err := k.GetStringValue(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		previewers = append(previewers, previewer{name: p, id: id})
	}

	blue := termenv.ANSI256Color(32)
	keywordStyle := termenv.String("keyword").Foreground(termenv.ANSI256Color(33)).Bold()
	cancelStyle := termenv.String("cancel").Foreground(termenv.ANSI256Color(196)).Bold()
	succeedStyle := termenv.String("succeed").Foreground(termenv.ANSI256Color(46)).Bold()

	func() {
		k, err := registry.OpenKey(registry.CLASSES_ROOT, `.`+extension+`\shellex\{8895b1c6-b41f-4c1c-a562-0d564250836f}`, registry.QUERY_VALUE)
		if err != nil {
			return
		}
		defer k.Close()

		id, _, err := k.GetStringValue("")
		if err != nil {
			return
		}

		for _, p := range previewers {
			if p.id == id {
				fmt.Println("The previewer for " + keywordStyle.Styled(extension) + " is currently set to " + keywordStyle.Styled(p.name) + ".")
				return
			}
		}

		fmt.Println("The previewer for " + keywordStyle.Styled(extension) + " is set but unknown (" + keywordStyle.Styled(id) + ").")
	}()

	previewers = append([]previewer{previewer{name: "None", id: ""}}, previewers...)

	sp := selection.New("Select a previewer", previewers)

	sp.SelectedChoiceStyle = func(c *selection.Choice[previewer]) string {
		return termenv.String(c.String).Foreground(blue).Bold().Styled(c.Value.name)
	}
	sp.UnselectedChoiceStyle = func(c *selection.Choice[previewer]) string {
		return c.Value.name
	}
	sp.FinalChoiceStyle = func(c *selection.Choice[previewer]) string {
		return termenv.String(c.String).Foreground(blue).Bold().Styled(c.Value.name)
	}

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Using " + keywordStyle.Styled(choice.name) + " for " + keywordStyle.Styled(extension) + ".")

	co := confirmation.New("Continue", confirmation.Undecided)
	co.Template = confirmation.TemplateYN

	confirmed, err := co.RunPrompt()
	if err != nil {
		fmt.Println(err)
		return
	}

	if confirmed {
		if choice.id == "" {
			k, err := registry.OpenKey(registry.CLASSES_ROOT, `.`+extension+`\shellex`, registry.SET_VALUE)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer k.Close()

			err = registry.DeleteKey(k, `{8895b1c6-b41f-4c1c-a562-0d564250836f}`)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			k, err := registry.OpenKey(registry.CLASSES_ROOT, `.`+extension, registry.CREATE_SUB_KEY)
			if err != nil {

			}
			defer k.Close()

			k, _, err = registry.CreateKey(k, `shellex`, registry.CREATE_SUB_KEY)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer k.Close()

			k, _, err = registry.CreateKey(k, `{8895b1c6-b41f-4c1c-a562-0d564250836f}`, registry.SET_VALUE)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer k.Close()

			err = k.SetStringValue("", choice.id)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println(succeedStyle.Styled("Done."))
	} else {
		fmt.Println(cancelStyle.Styled("Cancelled."))
	}
}

func OutlookPreview() types.ToolItem {
	return types.ToolItem{
		Run:         run,
		Name:        name,
		Description: desc,
	}
}
