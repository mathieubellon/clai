# clai

**clai** is a command-line tool that uses OpenAI's GPT models to turn natural language instructions into the best possible shell command for macOS, and then lets you execute it with a single keystroke.

## Features

- Converts your sentence into a shell command using GPT-3.5/4.
- Shows you the generated command before running it.
- Executes the command in your shell (supports pipes, redirection, etc.).
- Designed for macOS.

## Installation

1. **Clone the repository:**
   ```
   git clone https://github.com/yourusername/clai.git
   cd clai
   ```

2. **Install Go dependencies:**
   ```
   go mod tidy
   ```

3. **Build the CLI:**
   ```
   go build -o clai
   ```

## Usage

Set your OpenAI API key in your environment:
```
export OPENAI_API_KEY=sk-...
```

Run the CLI with your instruction:
```
./clai show me all files modified today
```

You’ll see the generated command. Press Enter to execute it, or Ctrl+C to cancel.

## Notes

- Requires Go 1.18+.
- Only the best command is generated (no options).
- The tool is tailored for macOS shell commands.
- Your OpenAI API key is required and billed by OpenAI.

## License

MIT