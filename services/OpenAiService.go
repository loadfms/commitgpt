package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/loadfms/commitgpt/models"
)

type AiModel string

const (
	AIMODEL_GPT35 AiModel = "gpt-3.5-turbo-1106"
)

const (
	ROUTE_COMPLETIONS = "https://api.openai.com/v1/chat/completions"
)

type OpenAiService struct {
	cfg           *ConfigService
	OpenAiRequest models.OpenAIRequest
}

func NewOpenAiService(cfg *ConfigService, aiModel AiModel, temperature float64) *OpenAiService {
	return &OpenAiService{
		cfg: cfg,
		OpenAiRequest: models.OpenAIRequest{
			Model:       string(aiModel),
			Temperature: temperature,
			Messages:    make([]models.OpenAIRequestMessages, 0),
		},
	}
}

func (s *OpenAiService) GetResponse(prompt string) (string, error) {
	message := models.OpenAIRequestMessages{
		Role:    "user",
		Content: prompt,
	}

	s.OpenAiRequest.Messages = append(s.OpenAiRequest.Messages, message)

	jsonBytes, err := json.Marshal(s.OpenAiRequest)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", ROUTE_COMPLETIONS, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res models.OpenAIReponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Could not interact with OpenAI: %s", resp.Status)
	}

	if res.Choices == nil || len(res.Choices) == 0 {
		return "", fmt.Errorf("No response from OpenAI.\nPlease, check your API KEY and try again.")
	}

	return res.Choices[0].Message.Content, nil
}
