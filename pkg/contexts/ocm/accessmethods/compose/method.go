// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"fmt"
	"io"
	"sync/atomic"

	"github.com/open-component-model/ocm/pkg/common/accessio/blobaccess"
	cpi "github.com/open-component-model/ocm/pkg/contexts/ocm/internal" // avoid cycle
	"github.com/open-component-model/ocm/pkg/runtime"
)

// Type is the access type of GitHub registry.
const (
	Type   = "compose"
	TypeV1 = Type + runtime.VersionSeparator + "v1"
)

func Is(spec cpi.AccessSpec) bool {
	return spec != nil && spec.GetKind() == Type
}

// AccessSpec describes the access for a GitHub registry.
type AccessSpec struct {
	runtime.ObjectVersionedType `json:",inline"`

	// Id is the internal id to identify the content
	Id string `json:"id"`

	// MediaType is the media type of the object represented by the blob
	MediaType string `json:"mediaType"`

	// GlobalAccess is an optional field describing a possibility
	// for a global access. If given, it MUST describe a global access method.
	GlobalAccess *cpi.AccessSpecRef `json:"globalAccess,omitempty"`
	// ReferenceName is an optional static name the object should be
	// use in a local repository context. It is use by a repository
	// to optionally determine a globally referencable access according
	// to the OCI distribution spec. The result will be stored
	// by the repository in the field ImageReference.
	// The value is typically an OCI repository name optionally
	// followed by a colon ':' and a tag
	ReferenceName string `json:"referenceName,omitempty"`
}

var (
	_ cpi.AccessSpec           = (*AccessSpec)(nil)
	_ cpi.HintProvider         = (*AccessSpec)(nil)
	_ cpi.GlobalAccessProvider = (*AccessSpec)(nil)
)

// New creates a new GitHub registry access spec version v1.
func New(hint string, mediaType string, global cpi.AccessSpec) *AccessSpec {
	id := fmt.Sprintf("compose-%d", number.Add(1))
	s := &AccessSpec{
		ObjectVersionedType: runtime.NewVersionedTypedObject(Type),
		Id:                  id,
		ReferenceName:       hint,
		MediaType:           mediaType,
		GlobalAccess:        cpi.NewAccessSpecRef(global),
	}
	return s
}

var number atomic.Int64

func (a *AccessSpec) Describe(ctx cpi.Context) string {
	return fmt.Sprintf("Composition blob %s", a.Id)
}

func (_ *AccessSpec) IsLocal(cpi.Context) bool {
	return true
}

func (a *AccessSpec) GetReferenceHint(cv cpi.ComponentVersionAccess) string {
	return a.ReferenceName
}

func (a *AccessSpec) GlobalAccessSpec(ctx cpi.Context) cpi.AccessSpec {
	if g, err := ctx.AccessSpecForSpec(a.GlobalAccess); err == nil {
		return g
	}
	return a.GlobalAccess.Unwrap()
}

func (_ *AccessSpec) GetType() string {
	return Type
}

func (a *AccessSpec) AccessMethod(cv cpi.ComponentVersionAccess) (cpi.AccessMethod, error) {
	return cv.AccessMethod(a)
}

func (a *AccessSpec) GetInexpensiveContentVersionIdentity(access cpi.ComponentVersionAccess) string {
	return a.Id
}

type accessMethod struct {
	access blobaccess.BlobAccess

	spec *AccessSpec
}

var _ cpi.AccessMethod = (*accessMethod)(nil)

func NewMethod(spec *AccessSpec, blob blobaccess.BlobAccess) (cpi.AccessMethod, error) {
	if blob.MimeType() != spec.MediaType {
		return nil, fmt.Errorf("mimetype mismatch (spec=%s, blob=%s)", spec.MediaType, blob.MimeType())
	}
	b, err := blob.Dup()
	if err != nil {
		return nil, err
	}
	return &accessMethod{
		access: b,
		spec:   spec,
	}, nil
}

func (_ *accessMethod) IsLocal() bool {
	return true
}

func (m *accessMethod) GetKind() string {
	return Type
}

func (m *accessMethod) MimeType() string {
	return m.access.MimeType()
}

func (m *accessMethod) AccessSpec() cpi.AccessSpec {
	return m.spec
}

func (m *accessMethod) Get() ([]byte, error) {
	return m.access.Get()
}

func (m *accessMethod) Reader() (io.ReadCloser, error) {
	return m.access.Reader()
}

func (m *accessMethod) Close() error {
	if m.access == nil {
		return nil
	}
	return m.access.Close()
}
