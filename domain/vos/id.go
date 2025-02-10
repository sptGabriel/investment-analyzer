package vos

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/sptGabriel/investment-analyzer/domain"
)

var ErrInvalidUUID = fmt.Errorf("%w:invalid uuid", domain.ErrMalformedParameters)

type ID struct {
	value uuid.UUID
}

func (id ID) IsZero() bool {
	return id.value == uuid.Nil
}

func (id ID) Value() string {
	if id.IsZero() {
		return ""
	}

	return id.value.String()
}

func (id ID) String() string {
	return id.value.String()
}

func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + id.String() + `"`), nil
}

func ParseID(s string) (ID, error) {
	uuidParsed, err := uuid.Parse(s)
	if err != nil {
		return ID{}, ErrInvalidUUID
	}

	return ID{uuidParsed}, nil
}

func MustID(s string) ID {
	return ID{value: uuid.MustParse(s)}
}

func NewID() ID {
	return ID{value: uuid.New()}
}
