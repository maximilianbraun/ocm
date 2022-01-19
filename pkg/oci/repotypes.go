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

package oci

import (
	"fmt"
	"strings"

	"github.com/gardener/ocm/pkg/common"
	"github.com/gardener/ocm/pkg/errors"
	"github.com/gardener/ocm/pkg/ocm/runtime"
)

type RepositoryType interface {
	runtime.TypedObjectDecoder
	common.VersionedElement
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

func NewRepositoryTypeScheme() RepositoryTypeScheme {
	var rt RepositorySpec
	scheme := runtime.MustNewDefaultScheme(&rt, &UnknownRepositorySpec{}, true)
	return &repositoryTypeScheme{scheme}
}

func (t *repositoryTypeScheme) GetRepositoryType(name string) RepositoryType {
	d := t.GetDecoder(name)
	if d == nil {
		return nil
	}
	return d.(RepositoryType)
}

func (t *repositoryTypeScheme) Register(name string, rtype RepositoryType) {
	t.RegisterByDecoder(name, rtype)
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
var DefaultRepositoryTypeScheme = NewRepositoryTypeScheme()

func RegisterRepositoryType(name string, atype RepositoryType) {
	DefaultRepositoryTypeScheme.Register(name, atype)
}

func CreateRepositorySpec(t runtime.TypedObject) (RepositorySpec, error) {
	return DefaultRepositoryTypeScheme.CreateRepositorySpec(t)
}

type UnknownRepositorySpec struct {
	runtime.UnstructuredTypedObject
}

var _ RepositorySpec = &UnknownRepositorySpec{}

func (r *UnknownRepositorySpec) Repository(Context) (Repository, error) {
	return nil, errors.ErrUnknown("respository type", r.GetType())
}

func (r *UnknownRepositorySpec) GetName() string {
	t := r.GetType()
	i := strings.LastIndex(t, "/")
	if i < 0 {
		return t
	}
	return t[:i]
}

func (r *UnknownRepositorySpec) GetVersion() string {
	t := r.GetType()
	i := strings.LastIndex(t, "/")
	if i < 0 {
		return "v1"
	}
	return t[i+1:]
}