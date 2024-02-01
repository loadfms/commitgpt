package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/cli/browser"
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
	err := handleArguments()
	if err != nil {
		if err.Error() != "done" {
			fmt.Println(err)
		}
		return
	}

	endpoint := "https://api.openai.com/v1/chat/completions"
	apiKey, err := getApiKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	changes, err := getChanges()
	if err != nil {
		fmt.Println(err)
		return
	}

	request := OpenAIRequest{
		Model:       "gpt-3.5-turbo-1106",
		Temperature: 0.5,
		Messages:    make([]OpenAIRequestMessages, 0),
	}

	message := OpenAIRequestMessages{
		Role:    "user",
		Content: changes,
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

func handleArguments() error {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "help":
			return help()
		case "auth":
			return auth()
		case "config":
		default:
			return fmt.Errorf("Invalid argument '%s'", args[0])
		}
	}
	return nil
}

func help() error {
	// Colors
	Reset := "\033[0m"
	White := "\033[97m"

	fmt.Println(White + "CommitGPT is a command-line tool that generates a commit message based on the changes in the git diff, following the conventional commits standard." + Reset)
	fmt.Println("")
	fmt.Println("Available commands for CommitGPT:")
	fmt.Println("")
	fmt.Println(White + "   auth:" + Reset)
	fmt.Println("     Configure your OpenAI credentials.")
	fmt.Println("     Redirects you to OpenAI Website, gets the API Key and automatically stores it.")
	fmt.Println("")
	return fmt.Errorf("done")
}

func auth() error {
	url := "https://platform.openai.com/api-keys"
	err := browser.OpenURL(url)
	if err != nil {
		return err
	}
	fmt.Println("Your browser has been opened to visit: ")
	fmt.Printf("  %s\n\n", url)

	fmt.Print("Paste your API Key here: ")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return err
	}

	path := "/.config/openai/"
	filename := "config"

	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("Error getting current user")
	}

	// Create directory and config file in case it doesn't exists
	// or open the config file in case it already exists and update its value
	os.Mkdir(currentUser.HomeDir+path, os.ModePerm)
	file, err := os.Create(currentUser.HomeDir + path + filename)
	if err != nil {
		file, err = os.Open(currentUser.HomeDir + path + filename)
		if err != nil {
			return fmt.Errorf("Error opening file in $HOME/.config/openai/config")
		}
	}
	defer file.Close()
	_, err = file.Write([]byte(input))
	if err != nil {
		return fmt.Errorf("Error writing OpenAI's config file")
	}

	return fmt.Errorf("done")
}

func getChanges() (string, error) {
	cmd := exec.Command("git", "status", "-v")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error executing git status")
	}

	if strings.Contains(string(out), "no changes added to commit") {
		return "", fmt.Errorf("No commits detected. HINT: Did you run 'git add .'?")
	}

	if strings.Contains(string(out), "nothing to commit, working tree clean") {
		return "", fmt.Errorf("No changes detected. Your working tree is clean.")
	}

	return fmt.Sprintf(`Write a commit message following the Conventional Commits standard and use Markdown formatting if needed. Please do not include the character count in the message, any author information or code snippet. The commit message should describe the changes made by this commit. these are changes:  %s`, out), nil
}

func getApiKey() (string, error) {
	filename := "/.config/openai/config"

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Error getting current user")
	}

	file, err := os.Open(currentUser.HomeDir + filename)
	if err != nil {
		return "", fmt.Errorf("Error opening file in $HOME/.config/openai/config")
	}
	defer file.Close()

	// new comment testing commit
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Error reading file")
	}

	return strings.TrimSuffix(string(data), "\n"), nil
}
