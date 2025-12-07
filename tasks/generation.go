package tasks

import (
	"bytes"
	"encoding/json"
	"english-util/domain"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type Generation struct {
	model string
	port  string
}

func NewGeneration(model, port string) *Generation {
	return &Generation{
		model: model,
		port:  port,
	}
}

const (
	generateSentencePrompt string = `Generate a single English sentence (at least 5 words) using the word "%s". 
		Return only the sentence, no extra text.`

	generateDescribePrompt string = `
		Входные данные:
		Оригинал на английском: [%s]
		Мой перевод на русский: [%s]

		Выходные данные:
		Мой перевод на русском.
		2-3 примера правильного перевода (НА РУССКОМ).
		Ответ в формате raw, без заголовков и тд
	`
)

func (g *Generation) GenerateSentenceTask(client *http.Client, input *domain.SentenceGenerationInput) (*domain.SentenceGenerationOutput, error) {
	timeNow := time.Now()

	reqBody := ollamaRequest{
		Model:  g.model,
		Prompt: fmt.Sprintf(generateSentencePrompt, input.Word),
		Stream: false,
	}

	body, err := g.createRequest(client, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var oResp ollamaResponse
	if err := json.Unmarshal(body, &oResp); err != nil {
		return nil, err
	}

	return &domain.SentenceGenerationOutput{
		Sentence:       oResp.Response,
		TimeToGenerate: time.Since(timeNow),
	}, nil
}

func (g *Generation) GenerateDescribeSentenceTask(client *http.Client, input *domain.SentenceTranslationInput) (*domain.SentenceTranslationOutput, error) {
	timeNow := time.Now()

	reqBody := ollamaRequest{
		Model:  g.model,
		Prompt: fmt.Sprintf(generateDescribePrompt, input.ToTranslate, input.Translated),
		Stream: false,
	}

	body, err := g.createRequest(client, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var oResp ollamaResponse
	if err := json.Unmarshal(body, &oResp); err != nil {
		return nil, err
	}

	return &domain.SentenceTranslationOutput{
		Discription:    oResp.Response,
		TimeToGenerate: time.Since(timeNow),
	}, nil
}

func (g *Generation) createRequest(client *http.Client, reqBody *ollamaRequest) ([]byte, error) {
	j, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:%s/api/generate", g.port),
		bytes.NewReader(j),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama error: %s", string(body))
	}
	return body, nil
}
