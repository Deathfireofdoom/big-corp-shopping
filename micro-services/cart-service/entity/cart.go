package entity

type Cart struct {
	userID 			string
	lastActivity 	time.Time
	productEntries 	map[string]productEntry
	pendingRequests map[string]InventoryRequest
}

type CartHandle struct {
	userID 	string
	cart	Cart
	mu		sync.Mutex
}