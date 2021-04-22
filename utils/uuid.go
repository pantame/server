package utils

import "github.com/google/uuid"

type UUID interface {
	String() string
}

func NewUUID() (UUID, error) {
	return uuid.NewUUID()
}
