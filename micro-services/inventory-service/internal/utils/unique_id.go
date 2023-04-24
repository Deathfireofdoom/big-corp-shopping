package utils

import 	"github.com/google/uuid"

func GetUniqueID() (string, error) {
	return uuid.New().String(), nil
}
