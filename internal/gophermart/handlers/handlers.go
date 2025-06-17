package handlers

import "fmt"

// Helpers

func parseID(id string) uint {
	var uid uint
	_, _ = fmt.Sscanf(id, "%d", &uid)
	return uid
}
