package main

import (
	"bytes"
	"log"
	"net/http"
)

func sendCommand(command string) {
	url := "http://localhost:17000"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(command)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-OK response: %v", resp.Status)
	}
}

func main() {
	sendCommand("green")
	sendCommand("bgrect 0.25 0.25 0.75 0.75")
	sendCommand("update")
}
