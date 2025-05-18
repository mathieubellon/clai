# CLAI

CLAI is a command-line tool that leverages OpenAI's GPT models to translate natural language instructions into executable shell commands on macOS. It provides a safe, interactive interface to review, explain, and execute commands generated from your input.

## Features

- **Natural Language to Shell Command**: Enter a sentence describing what you want to do; CLAI generates the best shell command for you.
- **Interactive Review**: See the generated command before execution.
- **Step-by-Step Explanation**: Get a detailed explanation of what the command does.
- **Clipboard Integration**: Copy the command to your clipboard with a single keypress.
- **Safe Execution**: Choose to execute, explain, copy, or quit—no surprises.

## Usage

```sh
clai "find all .txt files in my Documents folder"
```

Or run without arguments for interactive mode:

```sh
clai
```

You will be prompted to enter a sentence. After the command is generated, you’ll see:

```
Command: <generated command>
Options: [Enter]=Execute, [e]=Explain, [c]=Copy, [q]=Quit:
```

- **[Enter]**: Execute the command.
- **[e]**: Show a step-by-step explanation.
- **[c]**: Copy the command to your clipboard.
- **[q]**: Quit without executing.

## Setup

1. **Clone the repository**

   ```sh
   git clone git@github.com:mathieubellon/clai.git
   cd clai
   ```

2. **Build the CLI**

   ```sh
   go build -o clai main.go
   ```

3. **Set your OpenAI API key**

   Export your API key as an environment variable:

   ```sh
   export OPENAI_API_KEY=sk-...
   ```

4. **Install with Go**

  If you have Go installed, you can use `go install` to install CLAI directly:

  ```sh
  go install github.com/mathieubellon/clai@latest
  ```

  This will compile and install the CLI to your Go bin directory. Make sure your Go bin directory is in your PATH.

## Requirements

- Go 1.18+
- macOS (uses `pbcopy` for clipboard integration)
- OpenAI API key

## Security

- Always review generated commands before executing.
- The tool does not execute any command without your explicit confirmation.

## License

MIT License

---

*CLAI is not affiliated with OpenAI. Use at your own risk.*
