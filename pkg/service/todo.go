package service

import (
	"log"

	"github.com/satori/go.uuid"
)

type Todo struct {
	Id        uuid.UUID
	Title     string
	Completed bool
	Order     int
}

func TodoFrom(title string, order int) Todo {
	id, e := uuid.NewV4()
	if e != nil {
		log.Fatalf("Could not create UUID: %s", e)
	}

	return Todo{id, title, false, order}
}
