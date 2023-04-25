package cart

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"

	"big-corp-shopping/rest-api/internal/entity"
	"big-corp-shopping/rest-api/internal/cart_request_service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketMessage struct {
    Message string `json:"message"`
}


func UpdateProductToCart(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a ProductPayload struct.
	var payload entity.ProductPayload
	payloadHeader := []byte(r.Header.Get("payload"))
	err := json.Unmarshal(payloadHeader, &payload)

	// Check if decoding was successful.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the product and quantity values from the payload.
	product := payload.Product
	quantity := payload.Quantity

	// Fetches the user-id
	userID := r.Header.Get("User-ID")
	if userID == "" {
		http.Error(w, "Missing header User-ID", http.StatusBadRequest)
		return
	}
	

	// Make a new request
	var cartRequestAction entity.CartRequestAction
	if quantity == 0 {
		cartRequestAction = entity.CartRequestCheck
	} else if quantity < 0 {
		cartRequestAction = entity.CartRequestDelete
	} else {
		cartRequestAction = entity.CartRequestAdd
	}

	ch, err := cart_request_service.Service.NewCartRequest(userID, cartRequestAction, product, quantity)
	if err != nil {
		http.Error(w, "Serverside error.", http.StatusBadRequest)
		return
	}
	if cartRequestAction == entity.CartRequestDelete {
		return
	}
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Serverside error.", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Waiting for response to be handled
	log.Println("[info] waiting on response")
	cartRequestResponse :=<- ch

	// Marshaling response to it could be sent to client.
	data, err := json.Marshal(cartRequestResponse)
	if err != nil {
		http.Error(w, "Serverside error.", http.StatusBadRequest)
		return
	}

	// Sending message back to client.
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		http.Error(w, "Serverside error.", http.StatusBadRequest)
		return
	}

	log.Printf("[info] Successfully sent message over websocket.")
	
}

func RemoveProductFromCart(w http.ResponseWriter, r *http.Request) {
	panic("Implement this.")
}