package dockermulti

import (
	"golang.org/x/exp/slices"

	"github.com/open-component-model/ocm/pkg/common"
	"github.com/open-component-model/ocm/pkg/contexts/oci"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

type Option = optionutils.Option[*Options]

type Options struct {
	Context  oci.Context
	Version  string
	Variants []string
	Origin   *common.NameVersion
	Printer  common.Printer
}

func (o *Options) ApplyTo(opts *Options) {
	if opts == nil {
		return
	}
	if o.Version != "" {
		opts.Version = o.Version
	}
}

////////////////////////////////////////////////////////////////////////////////

type context struct {
	oci.Context
}

func (o context) ApplyTo(opts *Options) {
	opts.Context = o
}

func WithContext(ctx oci.ContextProvider) Option {
	return context{ctx.OCIContext()}
}

////////////////////////////////////////////////////////////////////////////////

type version string

func (o version) ApplyTo(opts *Options) {
	opts.Version = string(o)
}

func WithVersion(v string) Option {
	return version(v)
}

////////////////////////////////////////////////////////////////////////////////

type compvers common.NameVersion

func (o compvers) ApplyTo(opts *Options) {
	n := common.NameVersion(o)
	opts.Origin = &n
}

func WithOrigin(o common.NameVersion) Option {
	return compvers(o)
}

////////////////////////////////////////////////////////////////////////////////

type variants []string

func (o variants) ApplyTo(opts *Options) {
	opts.Variants = append(opts.Variants, []string(o)...)
}

func WithVariants(v ...string) Option {
	return variants(slices.Clone(v))
}

////////////////////////////////////////////////////////////////////////////////

type printer struct {
	common.Printer
}

func (o printer) ApplyTo(opts *Options) {
	opts.Printer = o
}

func WithPrinter(p common.Printer) Option {
	return printer{p}
}
