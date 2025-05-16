package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/eiannone/keyboard"
)

func main() {
	var sentence string
	if len(os.Args) > 1 {
		sentence = strings.Join(os.Args[1:], " ")
	} else {
		fmt.Print("Please enter a sentence: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		sentence = strings.TrimSpace(input)
	}

	options, err := getOpenAIOptions(sentence)
	if err != nil {
		fmt.Println("Error getting options:", err)
		return
	}

	selected, ok := selectOption(options)
	if !ok {
		fmt.Println("Selection cancelled")
		return
	}

	// Extract the actual content after the "Option X: " prefix
	parts := strings.SplitN(selected, ": ", 2)
	content := selected
	if len(parts) == 2 {
		content = parts[1]
	}

	fmt.Println(content)
	fmt.Println("Press Enter to execute this as a command...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	if strings.TrimSpace(content) == "" {
		fmt.Println("Empty command")
		return
	}
	cmd := exec.Command("sh", "-c", content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to execute command:", err)
	}
}

func getOpenAIOptions(sentence string) ([]string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Generate 3 command line options based on the user's input. Respond with three shell commands only, no explanations. No bullet, no line numbers, no numbered list. Purely the command line",
			},
			{
				"role":    "user",
				"content": sentence,
			},
		},
		"temperature": 0.7,
		"max_tokens":  150,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	content := ""
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		choice := choices[0].(map[string]interface{})
		msg := choice["message"].(map[string]interface{})
		content = strings.TrimSpace(msg["content"].(string))
	} else {
		content = "No response"
	}

	lines := strings.Split(content, "\n")
	options := []string{}
	for i, line := range lines {
		if i >= 3 {
			break
		}
		cmd := strings.TrimLeft(line, "-*> ")
		options = append(options, fmt.Sprintf("Option %d: %s", i+1, cmd))
	}

	// Fallbacks
	for len(options) < 3 {
		switch len(options) {
		case 0:
			options = append(options, fmt.Sprintf("Option 1: echo \"%s\"", sentence))
		case 1:
			options = append(options, fmt.Sprintf("Option 2: %s", strings.ToLower(sentence)))
		case 2:
			options = append(options, fmt.Sprintf("Option 3: %s", reverseString(sentence)))
		}
	}

	return options, nil
}

func selectOption(options []string) (string, bool) {
	selected := 0
	printOptions(options, selected)
	if err := keyboard.Open(); err != nil {
		fmt.Println("Failed to open keyboard:", err)
		return "", false
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return "", false
		}
		switch key {
		case keyboard.KeyArrowUp:
			if selected > 0 {
				selected--
				printOptions(options, selected)
			}
		case keyboard.KeyArrowDown:
			if selected < len(options)-1 {
				selected++
				printOptions(options, selected)
			}
		case keyboard.KeyEnter:
			return options[selected], true
		case keyboard.KeyCtrlC:
			return "", false
		default:
			if char == 'q' || char == 'Q' {
				return "", false
			}
		}
	}
}

func printOptions(options []string, selected int) {
	fmt.Print("\033[2J\033[H") // Clear screen
	for i, opt := range options {
		if i == selected {
			fmt.Printf("> %s\n", opt)
		} else {
			fmt.Printf("  %s\n", opt)
		}
	}
	fmt.Println("\nUse up/down arrows to navigate, Enter to select")
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
