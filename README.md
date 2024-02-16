# commitgpt

CommitGPT is a command-line tool that generates a commit message based on the changes in the git diff, following the conventional commits standard.

## Installation

To install Commit GPT, you need to have Go installed on your system. Then, you can run the following command:

```bash
go install github.com/loadfms/commitgpt@v1.1.1
```

To configure your access and preferences, run:
```bash
commitgpt auth
```

> If you dont have a key already, visit [Api Keys](https://platform.openai.com/account/api-keys)

## Usage

To generate a commit message, navigate to the root directory of your git repository and run the following command:

```bash
$ git commit -m "$(commitgpt)"
```

Commit GPT will analyze the changes in the git diff and generate a commit message based on the conventional commits standard.

> PRO TIP: create alias on your .zshrc with command
```bash
alias cgpt='git commit -m "$(commitgpt)"'
```

### Sample of Usage
```bash
$ git add .

$ cgpt

$ git push
```

## Uninstall
To uninstall just remove the bin file from your $GOPATH/bin
```bash
rm $GOPATH/bin/commitgpt
```

## Conventional Commits

The conventional commits standard is a lightweight convention on top of commit messages. It provides an easy way to communicate the nature of changes to other developers and tools that work with the repository.

A conventional commit message consists of a type, a scope, and a subject, followed by a body and a footer (optional). Here's an example of a conventional commit message:

```
feat(parser): add support for JSON input

Add support for parsing JSON input in the parser module.
```

In this example, the type is feat (for a new feature), the scope is parser, and the subject is add support for JSON input. The body provides more details about the changes, and the footer contains optional metadata, such as references to issues or breaking changes.
License

Commit GPT is licensed under the MIT license. See the LICENSE file for more information.
