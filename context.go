package gogo

import (
	stdContext "context"
)

type context struct {
	stdContext.Context
	Description string
	Help        string
	Arguments   map[any]argument
}

type argument struct {
	Description      string
	AllowedValues    []any
	RestrictedValues []any
	Default          any
}

type Context interface {
	stdContext.Context
	SetDescription(string) Context
	SetHelp(string) Context
	Argument(any) Argument
}

type Argument interface {
	Description(string) Argument
	AllowedValues(...any) Argument
	RestrictedValues(...any) Argument
	Default(any) Argument
}
