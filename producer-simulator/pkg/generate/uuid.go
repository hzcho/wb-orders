package generate

import "github.com/google/uuid"

func NewUUID() uuid.UUID {
	return uuid.New()
}