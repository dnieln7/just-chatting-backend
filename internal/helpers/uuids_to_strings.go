package helpers

import "github.com/google/uuid"

func UUIDsToStrings(uuids []uuid.UUID) []string {
	strings := []string{}

	for _, u := range uuids {
		strings = append(strings, u.String())
	}

	return strings
}
