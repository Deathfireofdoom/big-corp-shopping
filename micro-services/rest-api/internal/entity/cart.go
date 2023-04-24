package entity
import "time"

type Cart struct {
	UserID 			string
	ProductEntries 	map[string]ProductEntry
	LastActivity 	time.Time
}