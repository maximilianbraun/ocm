// Copyright 2022 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"fmt"
	"reflect"

	"github.com/gardener/ocm/pkg/common"
	"github.com/gardener/ocm/pkg/errors"
	"github.com/gardener/ocm/pkg/ocm/compdesc"
	"github.com/gardener/ocm/pkg/runtime"
)

type RepositoryType interface {
	runtime.TypedObjectDecoder
	common.VersionedElement

	// LocalSupportForAccessSpec checks whether a repository
	// provides a local version for the access spec.
	LocalSupportForAccessSpec(ctx Context, a compdesc.AccessSpec) bool
}

type RepositorySpec interface {
	runtime.TypedObject
	common.VersionedElement

	Repository(Context) (Repository, error)
}

type RepositoryTypeScheme interface {
	runtime.Scheme

	GetRepositoryType(name string) RepositoryType
	Register(name string, atype RepositoryType)

	DecodeRepositorySpec(data []byte, unmarshaler runtime.Unmarshaler) (RepositorySpec, error)
	CreateRepositorySpec(obj runtime.TypedObject) (RepositorySpec, error)
}

type repositoryTypeScheme struct {
	runtime.Scheme
}

func NewRepositoryTypeScheme(defaultRepoDecoder runtime.TypedObjectDecoder) RepositoryTypeScheme {
	var rt RepositorySpec
	scheme := runtime.MustNewDefaultScheme(&rt, &UnknownRepositorySpec{}, true, defaultRepoDecoder)
	return &repositoryTypeScheme{scheme}
}

func (t *repositoryTypeScheme) AddKnowntypes(s RepositoryTypeScheme) {
	t.Scheme.AddKnownTypes(s)
}

func (t *repositoryTypeScheme) GetRepositoryType(name string) RepositoryType {
	d := t.GetDecoder(name)
	if d == nil {
		return nil
	}
	return d.(RepositoryType)
}

func (t *repositoryTypeScheme) RegisterByDecoder(name string, decoder runtime.TypedObjectDecoder) error {
	if _, ok := decoder.(RepositoryType); !ok {
		errors.ErrInvalid("type", reflect.TypeOf(decoder).String())
	}
	return t.Scheme.RegisterByDecoder(name, decoder)
}

func (t *repositoryTypeScheme) AddKnownTypes(scheme runtime.Scheme) {
	if _, ok := scheme.(RepositoryTypeScheme); !ok {
		panic("can only add RepositoryTypeSchemes")
	}
	t.Scheme.AddKnownTypes(scheme)
}

func (t *repositoryTypeScheme) Register(name string, rtype RepositoryType) {
	t.Scheme.RegisterByDecoder(name, rtype)
}

func (t *repositoryTypeScheme) DecodeRepositorySpec(data []byte, unmarshaler runtime.Unmarshaler) (RepositorySpec, error) {
	obj, err := t.Decode(data, unmarshaler)
	if err != nil {
		return nil, err
	}
	if spec, ok := obj.(RepositorySpec); ok {
		return spec, nil
	}
	return nil, fmt.Errorf("invalid access spec type: yield %T instead of RepositorySpec")
}

func (t *repositoryTypeScheme) CreateRepositorySpec(obj runtime.TypedObject) (RepositorySpec, error) {
	if s, ok := obj.(RepositorySpec); ok {
		return s, nil
	}
	if u, ok := obj.(*runtime.UnstructuredTypedObject); ok {
		raw, err := u.GetRaw()
		if err != nil {
			return nil, err
		}
		return t.DecodeRepositorySpec(raw, runtime.DefaultJSONEncoding)
	}
	return nil, fmt.Errorf("invalid object type %T for repository specs", obj)
}

// DefaultRepositoryTypeScheme contains all globally known access serializer
var DefaultRepositoryTypeScheme = NewRepositoryTypeScheme(nil)

func RegisterRepositoryType(name string, atype RepositoryType) {
	DefaultRepositoryTypeScheme.Register(name, atype)
}

func CreateRepositorySpec(t runtime.TypedObject) (RepositorySpec, error) {
	return DefaultRepositoryTypeScheme.CreateRepositorySpec(t)
}

type UnknownRepositorySpec struct {
	runtime.UnstructuredVersionedTypedObject `json:",inline"`
}

var _ RepositorySpec = &UnknownRepositorySpec{}

func (r *UnknownRepositorySpec) Repository(Context) (Repository, error) {
	return nil, errors.ErrUnknown("respository type", r.GetType())
}
