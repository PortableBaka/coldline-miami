package core

import "github.com/google/uuid"

type Entity struct {
	ID   uuid.UUID
	name string
}
