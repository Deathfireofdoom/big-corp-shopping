package entity
import "time"

type Cart struct {
	userID 			string
	productEntries 	map[string]ProductEntry
	lastActivity 	time.Time
}