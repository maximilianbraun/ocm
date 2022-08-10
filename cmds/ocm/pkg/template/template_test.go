// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package template_test

import (
	"testing"

	"github.com/mandelsoft/vfs/pkg/osfs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/open-component-model/ocm/cmds/ocm/pkg/template"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Template Test Suite")
}

var _ = Describe("Template", func() {

	Context("Parse Arguments", func() {

		It("should parse one argument after a '--'", func() {
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			Expect(opts.FilterSettings("MY_VAR=test")).To(BeNil())
			Expect(opts.Vars).To(HaveKeyWithValue("MY_VAR", "test"))
		})

		It("should return non variable arguments", func() {
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			args := opts.FilterSettings("--", "MY_VAR=test", "my-arg")
			Expect(args).To(Equal([]string{
				"--", "my-arg",
			}))
			Expect(opts.Vars).To(HaveKeyWithValue("MY_VAR", "test"))
		})

		It("should parse multiple values", func() {
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			Expect(opts.FilterSettings("MY_VAR=test", "myOtherVar=true")).To(BeNil())
			Expect(opts.Vars).To(HaveKeyWithValue("MY_VAR", "test"))
			Expect(opts.Vars).To(HaveKeyWithValue("myOtherVar", "true"))
		})

		It("should filter multiple values", func() {
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			Expect(opts.FilterSettings("MY_VAR=test", "other")).To(Equal([]string{"other"}))
			Expect(opts.Vars).To(HaveKeyWithValue("MY_VAR", "test"))
		})
	})

	Context("Settings", func() {
		It("should filter multiple values", func() {
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			Expect(opts.ParseSettings(osfs.New(), "testdata/env.values")).To(Succeed())
			Expect(opts.Vars).To(HaveKeyWithValue("NAME", "test.de/x"))
			Expect(opts.Vars).To(HaveKeyWithValue("VERSION", "v1"))
		})
	})

	Context("Subst Template", func() {
		It("should template with a single value", func() {
			s := "my ${MY_VAR}"
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			opts.Vars = map[string]interface{}{
				"MY_VAR": "test",
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test"))
		})

		It("should template multiple value", func() {
			s := "my ${MY_VAR} ${my_second_var}"
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			opts.Vars = map[string]interface{}{
				"MY_VAR":        "test",
				"my_second_var": "testvalue",
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test testvalue"))
		})

		It("should use an empty string if no value is provided", func() {
			s := "my ${MY_VAR}"
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			opts.Vars = map[string]interface{}{}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my "))
		})

		It("should template with simple values", func() {
			s := "my ${MY_VAR}"
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			opts.Vars = map[string]interface{}{
				"MY_VAR": 5,
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my 5"))
		})

		It("should template with complex values", func() {
			s := "my ${MY_VAR}"
			opts := template.Options{}
			Expect(opts.Complete(nil)).To(Succeed())
			opts.Vars = map[string]interface{}{
				"MY_VAR": map[string]interface{}{
					"key": "value",
				},
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my {\"key\":\"value\"}"))
		})
	})

	Context("Go Template", func() {
		var opts *template.Options
		BeforeEach(func() {
			opts = &template.Options{Mode: "go"}
			Expect(opts.Complete(nil)).To(Succeed())
		})

		It("should template with a single value", func() {
			s := "my {{.MY_VAR}}"
			opts.Vars = map[string]interface{}{
				"MY_VAR": "test",
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test"))
		})

		It("should template multiple value", func() {
			s := "my {{.MY_VAR}} {{.my_second_var}}"
			opts.Vars = map[string]interface{}{
				"MY_VAR":        "test",
				"my_second_var": "testvalue",
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test testvalue"))
		})

		It("should not use an empty string if no value is provided", func() {
			s := "my {{.MY_VAR}}"
			opts.Vars = map[string]interface{}{}
			res, err := opts.Execute(s)
			_ = res
			Expect(err).To(HaveOccurred())
		})

		It("should template with simple values", func() {
			s := "my {{.MY_VAR}}"
			opts.Vars = map[string]interface{}{
				"MY_VAR": 5,
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my 5"))
		})

		It("should template with complex values", func() {
			s := "my {{.MY_VAR.key}}"
			opts.Vars = map[string]interface{}{
				"MY_VAR": map[string]interface{}{
					"key": "value",
				},
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my value"))
		})
	})

	Context("Spiff Template", func() {
		var opts *template.Options
		BeforeEach(func() {
			opts = &template.Options{Mode: "spiff"}
			Expect(opts.Complete(nil)).To(Succeed())
		})

		It("should template with a single value", func() {
			s := "my (( values.MY_VAR ))"
			opts.Vars = map[string]interface{}{
				"MY_VAR": "test",
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test\n"))
		})

		It("should template multiple value", func() {
			s := "my (( values.MY_VAR )) (( values.my_second_var ))"
			opts.Vars = map[string]interface{}{
				"MY_VAR":        "test",
				"my_second_var": "testvalue",
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my test testvalue\n"))
		})

		It("should not use an empty string if no value is provided", func() {
			s := "my (( values.MY_VAR ))"
			opts.Vars = map[string]interface{}{}
			res, err := opts.Execute(s)
			_ = res
			Expect(err).To(HaveOccurred())
		})

		It("should template with simple values", func() {
			s := "my (( values.MY_VAR ))"
			opts.Vars = map[string]interface{}{
				"MY_VAR": 5,
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my 5\n"))
		})

		It("should template with complex values", func() {
			s := "my (( values.MY_VAR.key ))"
			opts.Vars = map[string]interface{}{
				"MY_VAR": map[string]interface{}{
					"key": "value",
				},
			}
			res, err := opts.Execute(s)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(Equal("my value\n"))
		})
	})
})
