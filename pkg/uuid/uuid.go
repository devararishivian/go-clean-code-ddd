package uuid

import (
	"fmt"
	"github.com/gofrs/uuid"
)

func NewUUID() (string, error) {
	uuidV4, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %v", err)
	}

	return uuidV4.String(), nil
}
