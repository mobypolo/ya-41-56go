package customstrings

import "fmt"

func ParseID(id string) uint {
	var uid uint
	_, err := fmt.Sscanf(id, "%d", &uid)
	if err != nil {
		return 0
	}
	return uid
}
