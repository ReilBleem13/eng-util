package main

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	if m.quiiting {

		if len(m.wrongAnswers) != 0 {
			var res strings.Builder
			res.WriteString("\nСписок слов в которых были допущенны ошибки.\n\n")

			i := 1
			for _, ans := range m.wrongAnswers {
				res.WriteString(fmt.Sprintf("%d. %s - %s\n", i, ans.Word, strings.Join(ans.Translation, ", ")))
				i++
			}
			res.WriteString("\nПока! Жду снова!\n\n")
			return res.String()
		}

		return "\nПока! Жду снова!\n\n"
	}

	hint := "(esc — вернуться, enter — подтвердить)"

	if m.isLoading {
		return fmt.Sprintf(
			"\n%s %s\n\n(ожидание...)\n",
			m.spinner.View(),
			m.loadingMessage,
		)
	}

	if m.showActionMenu && m.choice != nil {
		options := []string{"Перевести слово", "Сгенерировать пример"}
		menu := ""

		for i, opt := range options {
			prefix := "  "
			if i == m.actionSelection {
				prefix = "> "
			}
			menu += fmt.Sprintf("%s%s\n", prefix, opt)
		}

		return fmt.Sprintf(
			"\nВыбрано слово: %s\n\n%s\n%s\n",
			m.choice.Word,
			menu,
			hint,
		)
	}

	if m.showExample {
		m.textInput.Width = 60
		m.textInput.Prompt = "Перевод: "

		inputLine := fmt.Sprintf("%*s%s", 4, "", m.textInput.View())

		return fmt.Sprintf(
			"\nСгенерировано предложение:\n%s\n\n%s\n\n%s\n",
			m.generatedSentence,
			inputLine,
			hint,
		)
	}

	if m.showDescribe {
		return fmt.Sprintf(
			"\nАнализ перевода:\n\n%s\n\n%s\n",
			m.describeResult,
			hint,
		)
	}

	if m.showInput && m.choice != nil {
		m.textInput.Width = 60
		m.textInput.Prompt = "› "

		inputLine := fmt.Sprintf("%*s%s", 4, "", m.textInput.View())

		return fmt.Sprintf(
			"\nСлово: %s\n\n%s\n\n%s\n",
			m.choice.Word,
			inputLine,
			hint,
		)
	}

	if m.showResult && m.choice != nil && m.inputResult != "" {
		res, correct := m.words.IsTranslationCorrect(
			m.inputResult,
			m.choice.Translation,
		)

		var msg string
		if correct {
			msg = fmt.Sprintf("Правильно!\n\nВарианты ответа:\n%s", res)
		} else {
			msg = fmt.Sprintf(
				"Неправильно!\nТы ввёл: %q\n\nВарианты ответа:\n%s",
				m.inputResult,
				res,
			)
		}

		return fmt.Sprintf("\n%s\n\n%s\n", msg, hint)
	}
	return "\n" + m.list.View()
}
