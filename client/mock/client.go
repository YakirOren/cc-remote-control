package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

const (
	uuidURL = "https://www.uuidtools.com/api/generate/v4"
	//wsURL   = "wss://cc-remote-control.fly.dev/ws"
	wsURL  = "ws://localhost:8080/ws"
	idFile = "id"
)

// Check if file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Get or generate UUID
func getOrCreateID() (string, error) {
	if !fileExists(idFile) {
		resp, err := http.Get(uuidURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var idArr []string
		if err := json.Unmarshal(content, &idArr); err != nil {
			return "", err
		}

		id := idArr[0]
		if err := ioutil.WriteFile(idFile, []byte(id), 0644); err != nil {
			return "", err
		}
	}

	id, err := ioutil.ReadFile(idFile)
	if err != nil {
		return "", err
	}

	return string(id), nil
}

// Main WebSocket loop
func websocketLoop() error {
	id, err := getOrCreateID()
	if err != nil {
		return err
	}

	headers := make(http.Header)
	headers.Set("User-Agent", id)

	dialer := websocket.DefaultDialer
	ws, _, err := dialer.Dial(wsURL, headers)
	if err != nil {
		return fmt.Errorf("error connecting to WebSocket: %v", err)
	}
	defer ws.Close()

	fmt.Printf("Connected as %s\n", id)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("error reading from WebSocket: %v", err)
		}
		fmt.Printf("Received: %s\n", message)

		var obj map[string]interface{}
		if err := json.Unmarshal(message, &obj); err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}

		action := obj["Action"].(string)
		switch action {
		case "eval":
			code := obj["Code"].(string)
			fmt.Println(code)
			err := ws.WriteMessage(websocket.TextMessage, []byte("Success"))
			if err != nil {
				return err
			}
		case "kill":
			ws.WriteMessage(websocket.TextMessage, []byte("bye bye"))
			return fmt.Errorf("kill action received")
		case "shell":
			code := obj["Code"].(string)
			cmd := exec.Command("sh", "-c", code)
			cmd.Run()
			ws.WriteMessage(websocket.TextMessage, []byte("No output currently, WIP"))
		}
	}
}

func main() {
	for {
		err := websocketLoop()
		if err != nil {
			fmt.Println(err)
			if err.Error() == "kill action received" {
				break
			}
			time.Sleep(30 * time.Second)
		}
	}
}
