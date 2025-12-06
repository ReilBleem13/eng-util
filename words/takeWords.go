package words

import (
	"bufio"
	"english-util/domain"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Words struct {
	path string
}

func NewWords(path string) *Words {
	return &Words{
		path: path,
	}
}

func (w *Words) TakeWords() []domain.Item {
	file, err := os.Open(w.path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	res, err := parseFile(file)
	if err != nil {
		fmt.Println("error occured: ", err)
	}
	return res
}

func (w *Words) IsTranslationCorrect(input string, correct []string) (string, bool) {
	var res strings.Builder

	flag := false
	c := 1

	input = strings.TrimSpace(strings.ToLower(input))
	for _, translation := range correct {
		res.WriteString(fmt.Sprintf("%d. %s\n", c, translation))
		c++
		if translation == input {
			flag = true
		}
	}
	return res.String(), flag
}

func parseFile(file *os.File) ([]domain.Item, error) {
	var res []domain.Item
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data := strings.TrimSpace(scanner.Text())

		if data == "" {
			continue
		}

		if strings.HasSuffix(data, " -") {
			newItem := domain.Item{}
			newItem.Word = strings.TrimSuffix(data, "-")

			if !scanner.Scan() {
				return nil, fmt.Errorf("word '%s' does not contain translation (unexpected EOF)", newItem.Word)
			}

			translationLine := strings.TrimSpace(scanner.Text())
			if translationLine == "" {
				return nil, fmt.Errorf("word '%s' has empty translation", newItem.Word)
			}

			newItem.Translation = translation(translationLine)
			res = append(res, newItem)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return res, nil
}

func translation(rawData string) []string {
	data := strings.Split(rawData, ",")

	res := make([]string, 0, len(data))
	for _, d := range data {
		res = append(res, strings.ToLower(removeNonLetterrsFromEnd(d)))
	}
	return res
}

func removeNonLetterrsFromEnd(s string) string {
	s = strings.TrimSpace(s)
	return strings.TrimRightFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}
