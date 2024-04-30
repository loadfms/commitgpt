package main

import (
	"fmt"
	"os"

	"github.com/loadfms/commitgpt/services"
)

func main() {
	// Get configs from Config file
	cfg := services.NewConfigService()

	// Get prompt message for GPT
	promptSvc := services.NewPromptService(cfg)

	// Create a new OpenAiService
	openAiSvc := services.NewOpenAiService(cfg, services.AIMODEL_GPT35, 0.5)

	// Handle arguments
	result, err := handleArguments(promptSvc, openAiSvc)
	if err != nil || result == "done" {
		// result == "done" is a workaround to avoid printing the error message
		// this means some command was executed successfully and the program should exit
		if result != "done" {
			fmt.Println(err)
		}
		return
	}
	fmt.Println(result)
}

func handleArguments(promptSvc *services.PromptService, openAiSvc *services.OpenAiService) (string, error) {
	commandService := services.NewCommandsService(promptSvc, openAiSvc)
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "help", "--help", "-h":
			return commandService.Help()
		case "auth", "--auth", "-a":
			return commandService.Auth()
		case "interactive", "-i", "--interactive":
			return commandService.Interactive(args[1:]) // Remove the interactive flag
		case "version", "--version":
			return commandService.Version()
		default:
			// As a default behaviour on args, we will use interactive mode
			return commandService.Interactive(args) // Pass the args to the interactive mode
		}
	}
	// Default behaviour `git commit -m "$(commitgpt)"` is handled here:

	// If no arguments are passed, we get the changes from git status
	prompt, err := promptSvc.GetChanges()
	if err != nil {
		return "", err
	}

	// Get the response from OpenAI
	result, err := openAiSvc.GetResponse(prompt)
	if err != nil {
		return "", err
	}

	return result, nil
}
