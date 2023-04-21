package health 

import (
    "net/http"
	"github.com/gorilla/websocket"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Everything is ok!"))
}


// Websocket testing
type WebSocketMessage struct {
    Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
