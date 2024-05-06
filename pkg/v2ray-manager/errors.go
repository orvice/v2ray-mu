package v2raymanager

import (
	"fmt"
	"strings"
)

var (
	TODOErr = fmt.Errorf("TODO error ")
)

func IsNotFoundError(e error) bool {
	if e == nil {
		return false
	}
	return strings.Contains(e.Error(), "not found")
}

func IsAlreadyExistsError(e error) bool {
	if e == nil {
		return false
	}
	return strings.Contains(e.Error(), "already exists")
}
