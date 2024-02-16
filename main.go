package main

import (
	"fmt"
	"os"

	"github.com/loadfms/commitgpt/services"
)

func main() {
	err := handleArguments()
	if err != nil {
		if err.Error() != "done" {
			fmt.Println(err)
		}
		return
	}

	cfg := services.NewConfigService()

	prompt := services.NewPromptService(cfg)
	changes, err := prompt.GetChanges()
	if err != nil {
		fmt.Println(err)
		return
	}

	svcOpenAi := services.NewOpenAiService(cfg, services.AIMODEL_GPT35, 0.5)
	result, err := svcOpenAi.GetResponse(changes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func handleArguments() error {
	commandService := services.NewCommandsService()
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "help":
			return commandService.Help()
		case "auth":
			return commandService.Auth()
		default:
			return fmt.Errorf("Invalid argument '%s'", args[0])
		}
	}
	return nil
}
