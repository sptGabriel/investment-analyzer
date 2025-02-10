package vos

import (
	"fmt"
	"strings"

	"github.com/sptGabriel/investment-analyzer/domain"
)

var (
	SideBuy  = Side{value: "buy"}
	SideSell = Side{value: "sell"}
	SideNil  = Side{}

	ErrUnknownSide = fmt.Errorf("%w:unknown side", domain.ErrFailedDependency)
)

type Side struct {
	value string
}

func (s Side) Value() string {
	return s.value
}

func (s Side) IsZero() bool {
	return s == Side{}
}

func ParseSide(s string) (Side, error) {
	switch strings.ToLower(s) {
	case "":
		return SideNil, nil
	case "buy":
		return SideBuy, nil
	case "sell":
		return SideSell, nil
	default:
		return Side{}, ErrUnknownSide
	}
}
