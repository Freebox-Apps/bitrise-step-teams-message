package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"encoding/json"
)

const (
	TeamsTokenEnv   = "teams_token"
	MessageTitleEnv = "message_title"
	MessageDescEnv  = "message_desc"
	DebugEnv        = "debug"
	DebugKeyOk      = "yes"
)

func main() {

	token := getToken()

	if len(token) == 0 {
		fmt.Printf("\n!!! Teams weebhook not configured !!!\n\n")
		os.Exit(1)
	}

	title := getTitle()
	message := getMessage()

	card, errCard := createCard(title, message)

	if errCard != nil {
		fmt.Printf("Failed to create message, error: %#v", errCard)
	}

	errSend := sendCard(token, card)
	if errSend != nil {
		fmt.Printf("Failed to send message, error: %#v", errSend)
	}

	if errSend != nil || errCard != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func createCard(title string, message string) ([]byte, error) {
	card := map[string]interface{}{
		"type": "message",
		"attachments": []map[string]interface{}{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]interface{}{
					"title":   title,
					"message": message,
				},
			},
		},
	}

	if isDebug() {
		fmt.Printf("\n    -------- Debug card output --------\n\n")
		prettyPrintCard(card)
		fmt.Printf("\n    -----------------------------------\n\n")
	}

	return json.Marshal(card)
}

func sendCard(token string, payload []byte) error {
	resp, err := http.Post(token, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Failed to POST message error: %#v | payload: %s", err, payload)
	}
	defer resp.Body.Close()
	return err
}

func getToken() string {
	return os.Getenv(TeamsTokenEnv)
}

func getTitle() string {
	title := os.Getenv(MessageTitleEnv)
	if len(title) != 0 {
		return title
	} else {
		return "Empty title"
	}
}

func getMessage() string {
	message := os.Getenv(MessageDescEnv)
	if len(message) != 0 {
		return message
	} else {
		return "No message found"
	}
}

func isDebug() bool {
	return os.Getenv(DebugEnv) == DebugKeyOk
}

func prettyPrintCard(card map[string]interface{}) {
	jsonBytes, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		fmt.Println("JSON Error :", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
