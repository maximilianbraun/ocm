// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package dockerdaemon

import (
	"fmt"

	"github.com/open-component-model/ocm/pkg/common/accessio/blobaccess"
	"github.com/open-component-model/ocm/pkg/common/accessio/blobaccess/spi"
	"github.com/open-component-model/ocm/pkg/contexts/oci"
	"github.com/open-component-model/ocm/pkg/contexts/oci/annotations"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/artifactset"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/docker"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

func (o *Options) OCIContext() oci.Context {
	if o.Context == nil {
		return oci.DefaultContext()
	}
	return o.Context
}

func ImageInfoFor(name string, opts ...Option) (locator string, version string, err error) {
	eff := optionutils.EvalOptions(opts...)

	locator, version, err = docker.ParseGenericRef(name)
	if err != nil {
		return "", "", err
	}

	if version == "" || version == "latest" || optionutils.AsValue(eff.OverrideVersion) {
		version = eff.Version
	}
	if version == "" {
		return "", "", fmt.Errorf("no version specified")
	}
	return locator, version, nil
}

func BlobAccessProviderForImageFromDockerDaemon(name string, opts ...Option) spi.BlobAccessProvider {
	return spi.BlobAccessProviderFunction(func() (spi.BlobAccess, error) {
		b, _, err := BlobAccessForImageFromDockerDaemon(name, opts...)
		return b, err
	})
}

func BlobAccessForImageFromDockerDaemon(name string, opts ...Option) (blobaccess.BlobAccess, string, error) {
	eff := optionutils.EvalOptions(opts...)
	ctx := eff.OCIContext()

	locator, version, err := ImageInfoFor(name, eff)
	if err != nil {
		return nil, "", err
	}
	spec := docker.NewRepositorySpec()
	repo, err := ctx.RepositoryForSpec(spec)
	if err != nil {
		return nil, "", err
	}
	ns, err := repo.LookupNamespace(locator)
	if err != nil {
		return nil, "", err
	}
	blob, err := artifactset.SynthesizeArtifactBlob(ns, version,
		func(art oci.ArtifactAccess) error {
			if eff.Origin != nil {
				art.Artifact().SetAnnotation(annotations.COMPVERS_ANNOTATION, eff.Origin.String())
			}
			return nil
		},
	)
	if err != nil {
		return nil, "", err
	}
	return blob, version, nil
}
