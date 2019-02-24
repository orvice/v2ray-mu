package server

import "strings"

func IsAlreadyExistsError(e error) bool {
	if strings.Contains(e.Error(), "already exists") {
		return true
	}
	return false
}
