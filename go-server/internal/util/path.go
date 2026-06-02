package util

import (
	"errors"
	"strconv"
	"strings"
)

// ParseIDFromPath извлекает целочисленный id из пути вида /prefix/{id}.
func ParseIDFromPath(path, prefix string) (int, error) {
	rawID := strings.TrimPrefix(path, prefix)
	if rawID == "" || strings.Contains(rawID, "/") {
		return 0, errors.New("invalid id")
	}
	uid, err := strconv.Atoi(rawID)
	if err != nil {
		return 0, err
	}
	return uid, nil
}
