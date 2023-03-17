package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type OpenAIRequest struct {
	Model       string                  `json:"model"`
	Messages    []OpenAIRequestMessages `json:"messages"`
	Temperature float64                 `json:"temperature"`
}
type OpenAIRequestMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIReponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

func main() {
	endpoint := "https://api.openai.com/v1/chat/completions"
	apiKey := getApiKey()

	request := OpenAIRequest{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
		Messages:    make([]OpenAIRequestMessages, 0),
	}

	message := OpenAIRequestMessages{
		Role:    "user",
		Content: getChanges(),
	}

	request.Messages = append(request.Messages, message)

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var res OpenAIReponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Choices[0].Message.Content)
}

func getChanges() string {
	cmd := exec.Command("git", "status", "-v")
	out, err := cmd.Output()
	if err != nil {
		panic("Error executing git status")
	}

	if strings.Contains(string(out), "no changes added to commit") {
		panic("No commits detected")
	}

	return fmt.Sprintf(`Write a commit message following the Conventional Commits standard and use Markdown formatting if needed. Please do not include the character count in the message, any author information or code snippet. The commit message should describe the changes made by this commit. these are changes:  %s`, out)
}

func getApiKey() string {
	filename := "/.config/openai/config"

	currentUser, err := user.Current()
	if err != nil {
		panic("Error getting current user")
	}

	file, err := os.Open(currentUser.HomeDir + filename)
	if err != nil {
		panic("Error opening file")
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Error reading file")
	}

	return strings.TrimSuffix(string(data), "\n")
}
