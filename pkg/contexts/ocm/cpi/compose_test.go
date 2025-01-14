// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package cpi_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/open-component-model/ocm/pkg/contexts/ocm/testhelper"
	. "github.com/open-component-model/ocm/pkg/finalizer"
	. "github.com/open-component-model/ocm/pkg/testutils"

	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/common/accessio/blobaccess"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/datacontext"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/attrs/compositionmodeattr"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/ctf"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
	"github.com/open-component-model/ocm/pkg/mime"
)

const COMPONENT = "github.com/mandelsoft/ocm"
const VERSION = "1.0.0"

var _ = Describe("access method", func() {
	var fs vfs.FileSystem
	var ctx ocm.Context

	BeforeEach(func() {
		ctx = ocm.New(datacontext.MODE_EXTENDED)
		fs = memoryfs.New()
	})

	DescribeTable("composes cv in one repo", func(mode bool) {
		final := Finalizer{}
		defer Defer(final.Finalize)

		compositionmodeattr.Set(ctx, mode)
		a := Must(ctf.Create(ctx, accessobj.ACC_WRITABLE|accessobj.ACC_CREATE, "ctf", 0o700, accessio.PathFileSystem(fs)))
		final.Close(a)
		c := Must(a.LookupComponent(COMPONENT))
		final.Close(c)

		cv := Must(c.NewVersion(VERSION))
		final.Close(cv)

		// add resource
		MustBeSuccessful(cv.SetResourceBlob(compdesc.NewResourceMeta("text1", resourcetypes.PLAIN_TEXT, metav1.LocalRelation), blobaccess.ForString(mime.MIME_TEXT, S_TESTDATA), "", nil))
		Expect(Must(cv.GetResource(compdesc.NewIdentity("text1"))).Meta().Digest).To(Equal(DS_TESTDATA))

		MustBeSuccessful(c.AddVersion(cv))
		MustBeSuccessful(final.Finalize())

		a = Must(ctf.Open(ctx, accessobj.ACC_READONLY, "ctf", 0o700, accessio.PathFileSystem(fs)))
		final.Close(a)

		cv = Must(a.LookupComponentVersion(COMPONENT, VERSION))
		final.Close(cv)

		Expect(Must(cv.GetResourcesByName("text1"))[0].Meta().Digest).To(Equal(DS_TESTDATA))
	},
		Entry("direct", false),
		Entry("compose", true),
	)

	DescribeTable("composes cv in one repo and add it to another", func(mode bool) {
		final := Finalizer{}
		defer Defer(final.Finalize)

		compositionmodeattr.Set(ctx, mode)
		a := Must(ctf.Create(ctx, accessobj.ACC_WRITABLE|accessobj.ACC_CREATE, "ctf", 0o700, accessio.PathFileSystem(fs)))
		final.Close(a)
		c := Must(a.LookupComponent(COMPONENT))
		final.Close(c)

		cv := Must(c.NewVersion(VERSION))
		final.Close(cv)

		// add resource
		MustBeSuccessful(cv.SetResourceBlob(compdesc.NewResourceMeta("text1", resourcetypes.PLAIN_TEXT, metav1.LocalRelation), blobaccess.ForString(mime.MIME_TEXT, S_TESTDATA), "", nil))
		Expect(Must(cv.GetResource(compdesc.NewIdentity("text1"))).Meta().Digest).To(Equal(DS_TESTDATA))

		a2 := Must(ctf.Create(ctx, accessobj.ACC_WRITABLE|accessobj.ACC_CREATE, "ctf2", 0o700, accessio.PathFileSystem(fs)))
		final.Close(a2)
		c2 := Must(a2.LookupComponent(COMPONENT))
		final.Close(c2)

		MustBeSuccessful(c2.AddVersion(cv))
		MustBeSuccessful(final.Finalize())

		a = Must(ctf.Open(ctx, accessobj.ACC_READONLY, "ctf", 0o700, accessio.PathFileSystem(fs)))
		final.Close(a)
		ExpectError(a.LookupComponentVersion(COMPONENT, VERSION)).To(MatchError(`component version "github.com/mandelsoft/ocm:1.0.0" not found: oci artifact "1.0.0" not found in component-descriptors/github.com/mandelsoft/ocm`))

		a2 = Must(ctf.Open(ctx, accessobj.ACC_READONLY, "ctf2", 0o700, accessio.PathFileSystem(fs)))
		final.Close(a2)

		cv = Must(a2.LookupComponentVersion(COMPONENT, VERSION))
		final.Close(cv)

		Expect(Must(cv.GetResourcesByName("text1"))[0].Meta().Digest).To(Equal(DS_TESTDATA))
	},
		Entry("direct", false),
		Entry("compose", true),
	)
})
