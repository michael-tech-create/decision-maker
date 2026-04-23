package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Data structures for JSON communication
type UserRequest struct {
	Message string `json:"message"`
}

type BotResponse struct {
	Reply string `json:"reply"`
}

// Data structure for the Bored API response
type BoredResponse struct {
	Activity string `json:"activity"`
}

var UserDecision = []string{
	"Sleep: Should I hit snooze or get up immediately?",
	"Hydration: Drink water first or go straight for coffee?",
	"Apparel: What clothes should I wear today?",
	"Breakfast: What should I eat for my first meal?",
}

// Function to fetch a random activity from the Bored API
func fetchExternalActivity() string {
	// Set a timeout so the bot doesn't hang if the API is slow
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://appbrewery.com")
	if err != nil {
		return "I was going to suggest something new, but I'm having trouble connecting to my brain! 🧠"
	}
	defer resp.Body.Close()

	var data BoredResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "I found something, but I can't quite describe it. Try asking again!"
	}
	return "How about this: " + data.Activity
}

func decisionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	var req UserRequest
	json.NewDecoder(r.Body).Decode(&req)
	input := strings.ToLower(req.Message)

	var reply string
	// SIMULATED TYPING DELAY
	time.Sleep(800 * time.Millisecond)

	switch {
	case strings.Contains(input, "bored") || strings.Contains(input, "something to do"):
		// Use the External Bored API
		reply = "Micheal: " + fetchExternalActivity()

	case strings.Contains(input, "food") || strings.Contains(input, "breakfast"):
		reply = "Micheal: You should definitely decide what to eat for breakfast! 🍳"

	case strings.Contains(input, "clothes") || strings.Contains(input, "outfit"):
		reply = "Micheal: Focus on your outfit today! What matches your vibe? 👕"

	default:
		// Pick from your original local list
		randomChoice := UserDecision[rand.Intn(len(UserDecision))]
		reply = "Micheal: My gut says... " + randomChoice
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(BotResponse{Reply: reply})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Serve the frontend files from a folder named "static"
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// The API endpoint
	http.HandleFunc("/api/decision", decisionHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
