package domain

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 15

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	defaultWidth      = 20
)

func ResultHintStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Margin(1, 0, 0, 2). // верх, право, низ, лево
		Foreground(lipgloss.Color("244")).
		Italic(true)
}

func HeloStyle() lipgloss.Style {
	return helpStyle
}

type Item struct {
	Word        string
	Translation []string
}

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Word)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewList(items []list.Item) *list.Model {
	l := list.New(items, ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Выбери слово для перевода"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return &l
}

// Структуры для генерации заданий
// Генерация предложения.
type SentenceGenerationInput struct {
	Word string `json:"word"`
}

type SentenceGenerationOutput struct {
	Sentence       string        `json:"sentence"`
	TimeToGenerate time.Duration `json:"time_to_generate"`
}

// Генерация анализа перевода.
type SentenceTranslationInput struct {
	ToTranslate string `json:"to_translate"`
	Translated  string `json:"translated"`
}

type SentenceTranslationOutput struct {
	Discription    string        `json:"discription"`
	TimeToGenerate time.Duration `json:"time_to_generate"`
}
