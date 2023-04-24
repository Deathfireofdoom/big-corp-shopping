package migration

import (
    "net/http"
)

func InitializeInventoryHandler(w http.ResponseWriter, r *http.Request) {
    initializeInventory()
	w.Write([]byte("Inventory is initialized"))
}
