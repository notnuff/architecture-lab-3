package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
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
	posX := 0.5
	posY := 0.5
	velX := 0.0003
	velY := 0.0007
	sendCommand("white")
	log.Print(fmt.Sprintf("figure %f %f", posX, posY))
	sendCommand(fmt.Sprintf("figure %f %f", posX, posY))
	for {
		sendCommand(fmt.Sprintf("move %f %f", velX, velY))
		posX += velX
		posY += velY

		// Check for boundary collision and reverse velocity if necessary
		if posX <= 0 || posX >= 1 {
			velX = -velX
		}
		if posY <= 0 || posY >= 1 {
			velY = -velY
		}

		sendCommand("update")
		time.Sleep(2 * time.Millisecond)
	}
}
