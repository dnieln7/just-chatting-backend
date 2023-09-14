package helpers

import "github.com/google/uuid"

func ContainsUUID(items []uuid.UUID, target uuid.UUID) bool {
	var contains bool = false

	for _, item := range items {
		if item == target {
			contains = true
			break
		}
	}

	return contains
}
