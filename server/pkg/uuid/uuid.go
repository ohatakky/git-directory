package uuid

import (
	"log"

	"github.com/google/uuid"
)

func String() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return uuid.String()
}
