package services

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"

	"github.com/cli/browser"

	"github.com/loadfms/commitgpt/models"
)

// We might want to pass args to the commands service
// if things get more complex
type CommandsService struct {
	prompt    *PromptService
	openAiSvc *OpenAiService
}

func NewCommandsService(prompt *PromptService, openAiSvc *OpenAiService) *CommandsService {
	return &CommandsService{
		prompt:    prompt,
		openAiSvc: openAiSvc,
	}
}

func (c *CommandsService) Help() (string, error) {
	// Colors
	Reset := "\033[0m"
	White := "\033[97m"

	fmt.Println(White + "CommitGPT is a command-line tool that generates a commit message based on the changes in the git diff, following the conventional commits standard." + Reset)
	fmt.Println("")
	fmt.Println("Available commands for CommitGPT:")
	fmt.Println("")
	fmt.Println(White + "   auth, --auth, -a:" + Reset)
	fmt.Println("     Configure your OpenAI credentials.")
	fmt.Println("     Redirects you to OpenAI Website, gets the API Key and automatically stores it.")
	fmt.Println("")
	fmt.Println(White + "   interactive, --interactive, -i:" + Reset)
	fmt.Println("     Generates a commit message based on the changes in the git diff.")
	fmt.Println("     The user can interact with the generated message and decide whether to apply it.")
	fmt.Println("     After typing the command, the user will be prompted to either accept the command [y], reject it [n] or to retry [r].")
	fmt.Println("")
	return "done", nil
}

func (c *CommandsService) Auth() (string, error) {
	url := "https://platform.openai.com/api-keys"
	err := browser.OpenURL(url)
	if err != nil {
		return "", err
	}
	fmt.Println("Your browser has been opened to visit: ")
	fmt.Printf("  %s\n\n", url)

	fmt.Print("Paste your API Key here: ")
	reader := bufio.NewReader(os.Stdin)
	inputApiKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return "", err
	}

	inputApiKey = strings.TrimSpace(inputApiKey)

	fmt.Print("Do you like to add a custom prompt? (Leave blank to default) ")
	reader = bufio.NewReader(os.Stdin)
	inputPrompt, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return "", err
	}

	inputPrompt = strings.TrimSpace(inputPrompt)

	if len(inputPrompt) == 0 {
		inputPrompt = models.DEFAULT_PROMPT
	} else {
		// Add the default prompt to the custom prompt
		// This is to ensure that the commit message is generated based on the changes in the git diff
		inputPrompt = fmt.Sprintf("%s, %s", inputPrompt, models.DEFAULT_PROMPT)
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Error getting current user")
	}

	// Create directory if it doesn't exist
	os.Mkdir(currentUser.HomeDir+models.CONFIG_FOLDER, os.ModePerm)

	// Construct the file path
	filePath := currentUser.HomeDir + models.CONFIG_FOLDER + models.FILENAME

	// Create a Config struct
	cfgContent := FileConfig{}

	// Assign content to the API key field
	cfgContent.ApiKey.Key = inputApiKey
	cfgContent.Prompt.Custom = inputPrompt

	SaveConfigFile(filePath, cfgContent)

	return "done", nil
}

func (c *CommandsService) Interactive(args []string) (string, error) {
	prompt, err := c.prompt.InteractivePrompt(args)
	if err != nil {
		return "", err
	}

	// Get the response from OpenAI
	result, err := c.openAiSvc.GetResponse(prompt)
	if err != nil {
		return "", err
	}

	// Replace the ' && ' to '\n' to print to the user
	Reset := "\033[0m"
	Green := "\033[32m"
	commands := Green + strings.ReplaceAll(result, " && ", "\n") + Reset

	var confirm string
	fmt.Printf("Here is the commands to execute: \n\n%s\n\nDo you want to apply it? [y/n/r]: ", commands)
	_, err = fmt.Scan(&confirm)
	if err != nil {
		return "", err
	}
	fmt.Println("")

	// If the user confirms, execute the command
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "r" {
		return "", fmt.Errorf("Command execution aborted. '%s'", confirm)
	}

	// If the user wants to retry, call the interactive function again
	if strings.ToLower(confirm) == "r" {
		return c.Interactive(args)
	}

	output, err := executeCommand(result)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func executeCommand(cmd string) (string, error) {
	// Split the result into commands
	commands := strings.Split(cmd, " && ")

	// Execute the commands
	var err error
	var output []byte
	for _, command := range commands {
		fmt.Println(command)
		output, err = exec.Command("bash", "-c", command).Output()
		if err != nil {
			return "", err
		}
	}
	return string(output), nil
}

// function to use regex to get string inside double quotes
func getCommitMessage(s string) string {
	re := regexp.MustCompile(`"([^"]+)"`)
	return re.FindString(s)
}
