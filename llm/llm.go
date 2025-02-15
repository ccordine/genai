package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

func Message(ollamaAPIURL string, modelName string, prompt string) (string, error) {
	response, err := generateResponse(prompt, ollamaAPIURL, modelName)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}
	return response, nil
}

func generateResponse(prompt string, ollamaAPIURL string, modelName string) (string, error) {

	fmt.Printf("Starting model %s...\n", modelName)

	// Prepare the request payload with adjusted temperature
	payload := map[string]interface{}{
		"model":  modelName,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"stream": false,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Make the HTTP POST request to Ollama API
	req, err := http.NewRequest("POST", ollamaAPIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("received non-OK response: %s", string(bodyBytes))
	}

	// Read the response body using json.Decoder
	var responseText string
	decoder := json.NewDecoder(resp.Body)
	for {
		fmt.Println("New response")
		var responseLine struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}

		// Decode the next JSON object in the stream
		err := decoder.Decode(&responseLine)
		if err != nil {
			if err == io.EOF {
				// End of response
				break
			}
			return "", fmt.Errorf("failed to decode response line: %v", err)
		}

		// Append the response text
		responseText += responseLine.Response

		// Optionally, check if generation is done
		if responseLine.Done {
			break
		}
	}

	// stop model
	// Construct the command
	cmd := exec.Command("ollama", "stop", modelName)

	// Run the command
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error stopping model %s: %v\n", modelName, err)
		return "", err
	}

	fmt.Printf("Successfully stopped model %s\n", modelName)

	return responseText, nil
}
