package main

import (
	ws2 "fish-game/ws/ws"
	"log"
	"net/http"
)

func main() {
	hub := ws2.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws2.ServeWS(hub, w, r)
	})

	log.Println("WebSocket listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
