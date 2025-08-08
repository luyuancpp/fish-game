package main

import (
	"fish-game/ws/ws"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", ws.ServeWS)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	fmt.Println("✅ WebSocket server at ws://localhost:8082/ws")
	http.ListenAndServe(":8082", nil)
}
