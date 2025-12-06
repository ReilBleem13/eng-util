package main

import (
	"english-util/domain"
	"english-util/tasks"
	"english-util/words"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type model struct {
	list              list.Model
	textInput         textinput.Model
	choice            *domain.Item
	wrongAnswers      []domain.Item
	showInput         bool
	showResult        bool
	showActionMenu    bool
	showExample       bool
	showDescribe      bool
	quiiting          bool
	inputResult       string
	actionSelection   int // 0 - перевести, 1 - пример
	generatedSentence string
	describeResult    string
	spinner           spinner.Model
	isLoading         bool
	loadingMessage    string
	client            *http.Client
	tasks             *tasks.Generation
	words             *words.Words
}

type sentenceGeneratedMsg struct {
	output *domain.SentenceGenerationOutput
	err    error
}

type describeGeneratedMsg struct {
	output *domain.SentenceTranslationOutput
	err    error
}

func initialModel(l list.Model, wordsClient *words.Words, tasksClient *tasks.Generation) model {
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	ti := textinput.New()
	ti.Placeholder = "Введи перевод..."
	ti.CharLimit = 156
	ti.Width = 80

	return model{
		list:           l,
		textInput:      ti,
		wrongAnswers:   []domain.Item{},
		showInput:      false,
		showResult:     false,
		showActionMenu: false,
		showExample:    false,
		showDescribe:   false,
		quiiting:       false,
		client:         &http.Client{Timeout: 30 * time.Second},
		spinner:        sp,
		isLoading:      false,
		tasks:          tasksClient,
		words:          wordsClient,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found:", err)
		os.Exit(1)
	}

	wordsClient := words.NewWords(os.Getenv("PATH_TO_WORDS"))
	tasksClient := tasks.NewGeneration(os.Getenv("MODEL"), os.Getenv("MODEL_PORT"))

	selectedWords := wordsClient.TakeWords()
	items := []list.Item{}

	for _, word := range selectedWords {
		items = append(items, word)
	}

	m := initialModel(*domain.NewList(items), wordsClient, tasksClient)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
