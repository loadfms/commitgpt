package models

const (
	CONFIG_FOLDER      = "/.config/openai/"
	FILENAME           = "config.toml"
	DEFAULT_PROMPT     = "Write a commit message following the Conventional Commits standard and use Markdown formatting if needed. Please do not include the character count in the message, any author information or code snippet. The commit message should describe the changes made by this commit. these are changes: "
	INTERACTIVE_PROMPT = "You are now in control of a bash. You need to control a git repository using all the necessary git commands - which need to be extremely precise to avoid errors. You need to send the git commands to do an operation. Send the command in a pipeline format. You are allowed to write a commit message if needed, following the Conventional Commits standard and use Markdown formatting if necessary. Here's the operation you need to do, based on this git changes: "
)
