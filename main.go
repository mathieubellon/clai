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
	"time"
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

	// Spinner setup
	done := make(chan struct{})
	go spinner(done)

	command, err := getOpenAICommand(sentence)
	close(done) // Stop the spinner

	fmt.Print("\r") // Clear spinner line

	if err != nil {
		fmt.Println("Error getting command:", err)
		return
	}

	fmt.Printf("Command: %s\n", command)
	fmt.Print("Options: [Enter]=Execute, [e]=Explain, [c]=Copy, [q]=Quit: ")

	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "":
		if strings.TrimSpace(command) == "" {
			fmt.Println("Empty command")
			return
		}
		// Insert command into shell history
		shell := os.Getenv("SHELL")
		var historyFile string
		if strings.HasSuffix(shell, "zsh") {
			historyFile = os.Getenv("HOME") + "/.zsh_history"
		} else if strings.HasSuffix(shell, "bash") {
			historyFile = os.Getenv("HOME") + "/.bash_history"
		} else {
			// Default to zsh history if unknown
			historyFile = os.Getenv("HOME") + "/.zsh_history"
		}
		f, err := os.OpenFile(historyFile, os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			defer f.Close()
			// For zsh, prepend ": <epoch>:0;" for proper history format
			if strings.HasSuffix(historyFile, ".zsh_history") {
				fmt.Fprintf(f, ": %d:0;%s\n", time.Now().Unix(), command)
			} else {
				fmt.Fprintln(f, command)
			}
		}
		cmd := exec.Command("sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to execute command:", err)
		}
	case "e":
		explanation, err := getCommandExplanation(command)
		if err != nil {
			fmt.Println("Error getting explanation:", err)
			return
		}
		fmt.Println("\nStep-by-step explanation:")
		fmt.Println(explanation)
	case "c":
		// Copy to clipboard using pbcopy (macOS)
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(command)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Failed to copy to clipboard:", err)
		} else {
			fmt.Println("Command copied to clipboard.")
		}
	case "q":
		fmt.Println("Quitting without executing.")
		return
	default:
		fmt.Println("Unknown option. Exiting.")
	}
}

func getOpenAICommand(sentence string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Given the user's input, output the single best shell command to accomplish the task. Respond with only the command line, no explanations, no bullets, no list, no extra text. The User is on MAcOSX",
			},
			{
				"role":    "user",
				"content": sentence,
			},
		},
		"temperature": 0.7,
		"max_tokens":  100,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
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

	return content, nil
}

// spinner prints a simple animation to indicate waiting for OpenAI response.
func spinner(done chan struct{}) {
	chars := []rune{'|', '/', '-', '\\'}
	i := 0
	for {
		select {
		case <-done:
			return
		default:
			fmt.Printf("\rWaiting for OpenAI response... %c", chars[i%len(chars)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// getCommandExplanation queries OpenAI for a step-by-step explanation of the command.
func getCommandExplanation(command string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert in shell commands. Given a shell command, explain step by step what it does. Be concise and clear.",
			},
			{
				"role":    "user",
				"content": command,
			},
		},
		"temperature": 0.3,
		"max_tokens":  300,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
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
		content = "No explanation available"
	}

	return content, nil
}
