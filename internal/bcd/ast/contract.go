package ast

import (
	"fmt"
	"strings"

	"github.com/mavryk-network/bcdhub/internal/bcd/base"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/mavryk-network/bcdhub/internal/bcd/forge"
	"github.com/mavryk-network/bcdhub/internal/bcd/formatter"
)

// Contract -
type Contract struct {
	Address

	Type Node
}

// NewContract -
func NewContract(depth int) *Contract {
	return &Contract{
		Address: Address{
			Default: NewDefault(consts.CONTRACT, 1, depth),
		},
	}
}

// MarshalJSON -
func (c *Contract) MarshalJSON() ([]byte, error) {
	return marshalJSON(consts.CONTRACT, c.Annots, c.Type)
}

// String -
func (c *Contract) String() string {
	var s strings.Builder

	s.WriteString(c.Default.String())
	s.WriteString(strings.Repeat(consts.DefaultIndent, c.Depth))
	s.WriteString(c.Type.String())

	return s.String()
}

// ParseType -
func (c *Contract) ParseType(node *base.Node, id *int) error {
	if err := c.Default.ParseType(node, id); err != nil {
		return err
	}
	contractType, err := typeNode(node.Args[0], c.Depth, id)
	if err != nil {
		return err
	}
	c.Type = contractType
	return nil
}

// ToMiguel -
func (c *Contract) ToMiguel() (*MiguelNode, error) {
	name := c.GetTypeName()
	var value string
	if c.Value != nil {
		value = c.Value.(string)
		if c.ValueKind == valueKindBytes {
			v, err := forge.UnforgeContract(value)
			if err != nil {
				return nil, err
			}
			value = v
		}
	}
	return &MiguelNode{
		Prim:  c.Prim,
		Type:  strings.ToLower(c.Prim),
		Value: value,
		Name:  &name,
	}, nil
}

// ToJSONSchema -
func (c *Contract) ToJSONSchema() (*JSONSchema, error) {
	s := &JSONSchema{
		Prim:    c.Prim,
		Type:    JSONSchemaTypeString,
		Default: "",
		Title:   c.GetTypeName(),
	}

	tree := &TypedAst{Nodes: []Node{c.Type}}
	tags := findViewContractInterfaces(tree)
	if len(tags) == 1 {
		s.Tag = tags[0]
	}
	return s, nil
}

// Docs -
func (c *Contract) Docs(inferredName string) ([]Typedef, string, error) {
	name := getNameDocString(c, inferredName)
	typedef := Typedef{
		Name: name,
		Type: fmt.Sprintf("contract(%s)", c.Type.GetPrim()),
		Args: make([]TypedefArg, 0),
	}

	if isSimpleDocType(c.Type.GetPrim()) {
		return []Typedef{typedef}, typedef.Type, nil
	}

	str, err := json.MarshalToString(c.Type)
	if err != nil {
		return nil, "", err
	}
	paramName := fmt.Sprintf("%s_param", c.GetName())
	parameter, err := formatter.MichelineStringToMichelson(str, true, formatter.DefLineSize)
	if err != nil {
		return nil, "", err
	}

	typedef.Type = fmt.Sprintf("contract(%s)", makeVarDocString(paramName))
	paramTypedef := Typedef{
		Name: paramName,
		Type: parameter,
	}
	return []Typedef{typedef, paramTypedef}, typedef.Type, nil
}

// Distinguish -
func (c *Contract) Distinguish(x Distinguishable) (*MiguelNode, error) {
	second, ok := x.(*Contract)
	if !ok {
		return nil, nil
	}
	return c.Default.Distinguish(&second.Default)
}

// EqualType -
func (c *Contract) EqualType(node Node) bool {
	if !c.Default.EqualType(node) {
		return false
	}
	second, ok := node.(*Contract)
	if !ok {
		return false
	}
	return c.Type.EqualType(second.Type)
}

// FindByName -
func (c *Contract) FindByName(name string, isEntrypoint bool) Node {
	if c.GetName() == name {
		return c
	}
	return nil
}

// FromJSONSchema -
func (c *Contract) FromJSONSchema(data map[string]interface{}) error {
	return setOptimizedJSONSchema(&c.Default, data, forge.UnforgeContract, ContractValidator)
}

// ToParameters -
func (c *Contract) ToParameters() ([]byte, error) {
	return c.Default.ToParameters()
}
