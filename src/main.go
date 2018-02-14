package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ----------------------------------------------------------------------------
//	Types
// ----------------------------------------------------------------------------

// Client - Данные юзера, находящегося в онлайне.
type Client struct {
	id   int64
	ship *Object
}

// Message - Сообщение.
type Message struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ----------------------------------------------------------------------------
//	Properties
// ----------------------------------------------------------------------------

var nextClientID int64
var clients = make(map[*websocket.Conn]*Client) // connected clients
var broadcast = make(chan Message)              // broadcast channel
var world *World
var upgrader websocket.Upgrader

// ----------------------------------------------------------------------------
//	Methods
// ----------------------------------------------------------------------------

func main() {

	world = CreateWorld()

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Объявляем роут по которому клиент получит сокет
	http.HandleFunc("/ws", handleConnection)

	// Слушаем пайп с сообщениями от воркеров и рассылаем их остальным воркерам
	go handleMessages()

	// Запускаем сервер
	log.Println("http server started on :8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func createClient(ws *websocket.Conn) {
	var client = new(Client)
	client.id = nextClientID
	client.ship = world.CreateObject()
	clients[ws] = client
	nextClientID++

	// Отправим клиенту мир
	sendWorld(ws)
}

func removeClient(ws *websocket.Conn) {
	world.RemoveObject(clients[ws].ship)
	delete(clients, ws)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register new client
	createClient(ws)

	var welcomeMessage Message
	welcomeMessage.Type = "world"
	welcomeMessage.Message = "Welcome!"
	ws.WriteJSON(welcomeMessage)

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			removeClient(ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		log.Printf("RECEIVED THIS SHITa: %v", msg)

		// Send it out to every client that is currently connected
		for ws := range clients {
			err := ws.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				ws.Close()
				removeClient(ws)
			}
		}
	}
}

func sendWorld(ws *websocket.Conn) {
	var msg Message
	msg.Type = "world"
	msg.Data = &world
	ws.WriteJSON(msg)
}
