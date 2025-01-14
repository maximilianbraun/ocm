// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package rpi

import (
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

type (
	Option                 = optionutils.Option[*Options]
	GeneralOptionsProvider = optionutils.NestedOptionsProvider[*Options]
)

type Options struct {
	Global cpi.AccessSpec
	Hint   string
}

var (
	_ optionutils.NestedOptionsProvider[*Options] = (*Options)(nil)
	_ optionutils.Option[*Options]                = (*Options)(nil)
)

func (w *Options) NestedOptions() *Options {
	return w
}

func (o *Options) ApplyTo(opts *Options) {
	if o.Global != nil {
		opts.Global = o.Global
	}
	if o.Hint != "" {
		opts.Hint = o.Hint
	}
}

type hint string

func (o hint) ApplyTo(opts *Options) {
	opts.Hint = string(o)
}

func WithHint(h string) Option {
	return hint(h)
}

func WrapHint[O any, P optionutils.OptionTargetProvider[*Options, O]](h string) optionutils.Option[P] {
	return optionutils.OptionWrapper[*Options, O, P](WithHint(h))
}

////////////////////////////////////////////////////////////////////////////////

type global struct {
	cpi.AccessSpec
}

func (o global) ApplyTo(opts *Options) {
	opts.Global = o.AccessSpec
}

func WithGlobalAccess(a cpi.AccessSpec) Option {
	return global{a}
}

func WrapGlobalAccess[O any, P optionutils.OptionTargetProvider[*Options, O]](a cpi.AccessSpec) optionutils.Option[P] {
	return optionutils.OptionWrapper[*Options, O, P](WithGlobalAccess(a))
}
