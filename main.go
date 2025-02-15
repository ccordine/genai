package main

import (
	"encoding/json"
	"fmt"
	"github.com/ccordine/genai/llm"
	"log"
	"net/http"
)

const (
	ollamaAPIURL = "http://localhost:11434/api/generate"
)

// Request represents the expected JSON structure for incoming requests.
type Request struct {
	ModelName string `json:"model_name"`
	Prompt    string `json:"prompt"`
}

// Response represents the JSON structure for responses.
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Optional: disallow unknown fields
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := llm.Message(ollamaAPIURL, req.ModelName, req.Prompt)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}
	// Create a response
	resp := Response{
		Status:  "success",
		Message: "Model " + req.ModelName + " received with prompt: " + req.Prompt,
		Result:  result,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Optional: set the status code

	// Encode the response as JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fmt.Println("Running...")
	http.HandleFunc("/api/ai", handler)
	http.ListenAndServe(":8089", nil)
}
