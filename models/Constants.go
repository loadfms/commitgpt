package models

const (
	NAME               = "CommitGPT"
	VERSION            = "1.3.0"
	CONFIG_FOLDER      = "/.config/openai/"
	FILENAME           = "config.toml"
	DEFAULT_PROMPT     = "Write a commit message following the Conventional Commits standard and use Markdown formatting if needed. Please do not include the character count in the message, any author information or code snippet. The commit message should describe the changes made by this commit. these are changes: "
	INTERACTIVE_PROMPT = "You are now in control of a bash. You need to control a git repository using all the necessary git commands - which need to be extremely precise to avoid errors. You need to send the git commands to do an operation. Send the command in a pipeline format. You are allowed to write a commit message if needed, following the Conventional Commits standard and use Markdown formatting if necessary. As for your response, no formattation on commands, no markdown, (like bash) just the commands in plain text. Also, do not break lines, but instead, use && (pipeline). Do as the prompt says, if it says 'push to origin' include a 'git push -u origin $(branch)' command and so on. Use your git skills to find solutions using git commands. Here's the operation you need to do, based on this git changes: "
)
