package entity

import (
	"sync"
	"time"
)


type Cart struct {
	UserID 			string
	LastActivity 	time.Time
	ProductEntries 	map[string]ProductEntry
	PendingRequests map[string]InventoryRequest
	Hold			bool
}

type CartHandle struct {
	UserID 	string
	Cart	Cart
	Mu		sync.Mutex
}

func NewCartHandle(userID string) *CartHandle {
	return &CartHandle{
		UserID: userID,
		Cart: Cart{
			UserID: userID,
			LastActivity: time.Now(),
			ProductEntries: make(map[string]ProductEntry),
			PendingRequests: make(map[string]InventoryRequest),
			Hold: true,
		},
		Mu: sync.Mutex{},
	}
}