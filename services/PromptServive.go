package services

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/loadfms/commitgpt/models"
)

type PromptService struct {
	cfg *ConfigService
}

func NewPromptService(cfg *ConfigService) *PromptService {
	return &PromptService{cfg: cfg}
}

func (s *PromptService) GetChanges() (string, error) {
	cmd := exec.Command("git", "status", "-v")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error executing git status")
	}

	if strings.Contains(string(out), "no changes added to commit") {
		return "", fmt.Errorf("No commits detected. HINT: Did you run 'git add .'?")
	}

	// This shit will cause bugs, because if we have some changes
	// that contains this string voialà, we will not be able to commit!
	// git status -v at first is dangerous, and might be avoided before doing this validations
	if strings.Contains(string(out), "nothing to commit, working tree clean") {
		return "", fmt.Errorf("No changes detected. Your working tree is clean.")
	}

	return fmt.Sprintf(`%s %s`, s.cfg.CustomPrompt, out), nil
}

// This function should get the changes from git status
// just like the function above, but now, it will be allowed
// to add everything to staging in case nothing is in it.
func (s *PromptService) GetAllChanges() (string, error) {
	status, err := exec.Command("git", "status").Output()
	if err != nil {
		return "", fmt.Errorf("Error executing git status")
	}

	if strings.Contains(string(status), "no changes added to commit") {
		// Add all changes to the git stage
		_, err := exec.Command("git", "add", ".").Output()
		if err != nil {
			return "", fmt.Errorf("Error executing git add")
		}
		// Get the changes again
		status, err = exec.Command("git", "status").Output()
		if err != nil {
			return "", fmt.Errorf("Error executing git status")
		}
	}

	// Ignore this validation for now
	// For example, if the user just want to execute something like this:
	// `commitgpt -i "dude, I just want to reset the last pushed commit and rename it to "fix: something"`
	// This will not work, because the user will not have any changes in the working tree

	// if strings.Contains(string(status), "nothing to commit, working tree clean") {
	// 	return "", fmt.Errorf("No changes detected. Your working tree is clean.")
	// }

	changes, err := exec.Command("git", "status", "-v").Output()
	if err != nil {
		return "", fmt.Errorf("Error executing git status -v")
	}

	return fmt.Sprintf(`%s %s`, s.cfg.CustomPrompt, string(changes)), nil
}

func (s *PromptService) InteractivePrompt(args []string) (string, error) {
	// Get all the arguments passed to the command
	arguments := args[1:]

	// Check if the prompt is empty
	if len(arguments) == 0 {
		return "", fmt.Errorf("No prompt provided. Please provide a prompt.")
	}

	// Basically converting the arguments to a string
	prompt := strings.Join(arguments, " ")

	// Get all the changes
	changes, err := s.GetAllChanges()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`%s %s Prompt: ChatGPT, %s`, models.INTERACTIVE_PROMPT, changes, prompt), nil
}
