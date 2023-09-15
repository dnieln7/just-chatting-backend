package helpers

import "github.com/google/uuid"

func RemoveUUID(items []uuid.UUID, target uuid.UUID) []uuid.UUID {
	var index = -1
	var last = len(items) - 1

	for i, item := range items {
		if item == target {
			index = i
			break
		}
	}

	switch index {
	case -1:
		break
	case 0:
		items = items[1:]
	case last:
		items = items[:last]
	default:
		items[index] = items[last]
		items = items[:last]
	}

	return items
}
