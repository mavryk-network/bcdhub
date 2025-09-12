package ast

import "github.com/mavryk-network/bcdhub/internal/bcd/consts"

//
//  MUMAV
//

// Mumav -
type Mumav struct {
	Default
}

// NewMumav -
func NewMumav(depth int) *Mumav {
	return &Mumav{
		Default: NewDefault(consts.MUMAV, 0, depth),
	}
}

// ToJSONSchema -
func (m *Mumav) ToJSONSchema() (*JSONSchema, error) {
	return getIntJSONSchema(m.Default), nil
}

// Compare -
func (m *Mumav) Compare(second Comparable) (int, error) {
	secondItem, ok := second.(*Mumav)
	if !ok {
		return 0, consts.ErrTypeIsNotComparable
	}
	return compareBigInt(m.Default, secondItem.Default), nil
}

// Distinguish -
func (m *Mumav) Distinguish(x Distinguishable) (*MiguelNode, error) {
	second, ok := x.(*Mumav)
	if !ok {
		return nil, nil
	}
	return m.Default.Distinguish(&second.Default)
}

// FromJSONSchema -
func (m *Mumav) FromJSONSchema(data map[string]interface{}) error {
	setIntJSONSchema(&m.Default, data)
	return nil
}

// FindByName -
func (m *Mumav) FindByName(name string, isEntrypoint bool) Node {
	if m.GetName() == name {
		return m
	}
	return nil
}
