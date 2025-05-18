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

	fmt.Println(command)
	fmt.Println("Press Enter to execute this as a command...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	if strings.TrimSpace(command) == "" {
		fmt.Println("Empty command")
		return
	}
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to execute command:", err)
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
