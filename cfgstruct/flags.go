// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package cfgstruct

import (
	"strings"
	"time"

	"github.com/spf13/pflag"
)

// FlagSet is an interface that matches *pflag.FlagSet.
type FlagSet interface {
	BoolVar(p *bool, name string, value bool, usage string)
	IntVar(p *int, name string, value int, usage string)
	Int64Var(p *int64, name string, value int64, usage string)
	UintVar(p *uint, name string, value uint, usage string)
	Uint64Var(p *uint64, name string, value uint64, usage string)
	DurationVar(p *time.Duration, name string, value time.Duration, usage string)
	Float64Var(p *float64, name string, value float64, usage string)
	StringVar(p *string, name string, value string, usage string)
	StringArrayVar(p *[]string, name string, value []string, usage string)
	Var(val pflag.Value, name string, usage string)
	MarkHidden(name string) error
}

var _ FlagSet = (*pflag.FlagSet)(nil)

// commaDelimitedStrings implements a flag value with comma delimited string list.
type commaDelimitedStrings struct {
	changed bool
	// list uses a pointer to slice, so we can bind it to an existing config field.
	list *[]string
}

// Type implements pflag.Value.
func (*commaDelimitedStrings) Type() string { return "[]string" }

// String returns the values as comma delimited.
func (xs *commaDelimitedStrings) String() string {
	return strings.Join(*xs.list, ",")
}

// Set implements flag.Value interface.
func (xs *commaDelimitedStrings) Set(s string) error {
	if s == "" {
		return nil
	}
	if !xs.changed {
		*xs.list = nil
		xs.changed = true
	}
	*xs.list = append(*xs.list, strings.Split(s, ",")...)
	return nil
}
