package main

import (
	"english-util/domain"
	"english-util/priority"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) refreshListFromQueue() {
	items := make([]*priority.Item, len(m.queue.PQ))
	copy(items, m.queue.PQ)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = *it.Data
	}

	m.list.SetItems(listItems)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		if m.showInput || m.showExample {
			m.textInput.Width = msg.Width - 10
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		if key == "q" || key == "ctrl+c" {
			m.quiiting = true
			return m, tea.Quit
		}

		if !m.showActionMenu && !m.showExample && !m.showDescribe {
			switch key {
			case "h", "left", "pgup", "page up":
				m.list.PrevPage()
				return m, nil
			case "l", "right", "pgdn", "page down":
				m.list.NextPage()
				return m, nil
			}
		}

		if m.showActionMenu {
			switch key {
			case "up", "k":
				m.actionSelection = 0
				return m, nil
			case "down", "j":
				m.actionSelection = 1
				return m, nil
			}
		}

		if key == "enter" {
			if m.isLoading {
				return m, nil
			}

			if m.showActionMenu {
				if m.actionSelection == 0 {
					m.showActionMenu = false
					m.showInput = true
					return m, m.textInput.Focus()
				}

				if m.actionSelection == 1 {
					m.showActionMenu = false
					m.showExample = true
					m.isLoading = true
					m.loadingMessage = "Генерирую пример..."

					input := &domain.SentenceGenerationInput{
						Word: m.choice.Word,
					}

					return m, tea.Batch(
						func() tea.Msg {
							out, err := m.tasks.GenerateSentenceTask(m.client, input)
							return sentenceGeneratedMsg{out, err}
						},
						m.spinner.Tick,
					)
				}
			}

			if m.showExample {
				m.inputResult = m.textInput.Value()
				m.showExample = false
				m.isLoading = true
				m.loadingMessage = "Анализирую перевод..."

				input := &domain.SentenceTranslationInput{
					ToTranslate: m.generatedSentence,
					Translated:  m.inputResult,
				}

				return m, tea.Batch(
					func() tea.Msg {
						out, err := m.tasks.GenerateDescribeSentenceTask(m.client, input)
						return describeGeneratedMsg{out, err}
					},
					m.spinner.Tick,
				)
			}

			if m.showInput {
				m.inputResult = m.textInput.Value()
				m.textInput.SetValue("")
				m.showInput = false
				m.showResult = true

				_, correct := m.words.IsTranslationCorrect(m.inputResult, m.choice.Translation)
				if !correct {
					m.queue.Increase(m.choice.Word)
					m.wrongAnswers = append(m.wrongAnswers, domain.Item{
						Word:        m.choice.Word,
						Translation: m.choice.Translation,
					})
				} else {
					m.queue.Decrease(m.choice.Word)
				}

				m.refreshListFromQueue()

				return m, nil
			}

			if !m.showResult && !m.showDescribe {
				if sel, ok := m.list.SelectedItem().(domain.Item); ok {
					m.choice = &domain.Item{
						Word:        sel.Word,
						Translation: sel.Translation,
					}
					m.showActionMenu = true
					m.actionSelection = 0
				}
			}

			return m, nil
		}

		if key == "esc" {
			if m.isLoading {
				return m, nil
			}

			switch {
			case m.showInput:
				m.showInput = false
				m.textInput.SetValue("")
			case m.showResult:
				m.showResult = false
				m.inputResult = ""
			case m.showActionMenu:
				m.showActionMenu = false
				m.choice = nil
			case m.showExample:
				m.showExample = false
				m.generatedSentence = ""
				m.inputResult = ""
			case m.showDescribe:
				m.showDescribe = false
				m.describeResult = ""
				m.inputResult = ""
				m.generatedSentence = ""
			}
			return m, nil
		}

		if m.showInput || m.showExample {
			m.textInput, cmd = m.textInput.Update(msg)
		} else if !m.showResult && !m.showActionMenu && !m.showDescribe {
			m.list, cmd = m.list.Update(msg)
		}
		return m, cmd

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case sentenceGeneratedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.describeResult = "Ошибка генерации: " + msg.err.Error()
			m.showExample = false
			m.showDescribe = true
			return m, nil
		}

		m.generatedSentence = msg.output.Sentence
		m.textInput.SetValue("")
		return m, m.textInput.Focus()

	case describeGeneratedMsg:
		m.isLoading = false
		if msg.err != nil {
			m.describeResult = "Ошибка анализа: " + msg.err.Error()
		} else {
			m.describeResult = msg.output.Discription
		}
		m.showDescribe = true
		return m, nil
	}
	return m, nil
}
