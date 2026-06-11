package ast

import (
	"github.com/mavryk-network/nexushub/internal/bcd/consts"
)

// Parameter -
type Parameter struct {
	*SectionType
}

// NewParameter -
func NewParameter(depth int) *Parameter {
	return &Parameter{
		SectionType: NewSectionType(consts.PARAMETER, depth),
	}
}
