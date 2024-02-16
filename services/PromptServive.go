package services

import (
	"fmt"
	"os/exec"
	"strings"
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

	if strings.Contains(string(out), "nothing to commit, working tree clean") {
		return "", fmt.Errorf("No changes detected. Your working tree is clean.")
	}

	return fmt.Sprintf(`%s %s`, s.cfg.CustomPrompt, out), nil
}
