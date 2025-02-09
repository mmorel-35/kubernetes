/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
	v1 "k8s.io/code-generator/examples/crd/apis/conflicting/v1"
	examplev1 "k8s.io/code-generator/examples/crd/apis/example/v1"
	example2v1 "k8s.io/code-generator/examples/crd/apis/example2/v1"
	extensionsv1 "k8s.io/code-generator/examples/crd/apis/extensions/v1"
	conflictingv1 "k8s.io/code-generator/examples/crd/applyconfiguration/conflicting/v1"
	applyconfigurationexamplev1 "k8s.io/code-generator/examples/crd/applyconfiguration/example/v1"
	applyconfigurationexample2v1 "k8s.io/code-generator/examples/crd/applyconfiguration/example2/v1"
	applyconfigurationextensionsv1 "k8s.io/code-generator/examples/crd/applyconfiguration/extensions/v1"
	internal "k8s.io/code-generator/examples/crd/applyconfiguration/internal"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) any {
	switch kind {
	// Group=conflicting.test.crd.code-generator.k8s.io, Version=v1
	case v1.SchemeGroupVersion.WithKind("TestEmbeddedType"):
		return &conflictingv1.TestEmbeddedTypeApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("TestType"):
		return &conflictingv1.TestTypeApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("TestTypeStatus"):
		return &conflictingv1.TestTypeStatusApplyConfiguration{}

		// Group=example.crd.code-generator.k8s.io, Version=v1
	case examplev1.SchemeGroupVersion.WithKind("ClusterTestType"):
		return &applyconfigurationexamplev1.ClusterTestTypeApplyConfiguration{}
	case examplev1.SchemeGroupVersion.WithKind("ClusterTestTypeStatus"):
		return &applyconfigurationexamplev1.ClusterTestTypeStatusApplyConfiguration{}
	case examplev1.SchemeGroupVersion.WithKind("TestType"):
		return &applyconfigurationexamplev1.TestTypeApplyConfiguration{}
	case examplev1.SchemeGroupVersion.WithKind("TestTypeStatus"):
		return &applyconfigurationexamplev1.TestTypeStatusApplyConfiguration{}

		// Group=example.test.crd.code-generator.k8s.io, Version=v1
	case example2v1.SchemeGroupVersion.WithKind("TestType"):
		return &applyconfigurationexample2v1.TestTypeApplyConfiguration{}
	case example2v1.SchemeGroupVersion.WithKind("TestTypeStatus"):
		return &applyconfigurationexample2v1.TestTypeStatusApplyConfiguration{}

		// Group=extensions.test.crd.code-generator.k8s.io, Version=v1
	case extensionsv1.SchemeGroupVersion.WithKind("TestSubresource"):
		return &applyconfigurationextensionsv1.TestSubresourceApplyConfiguration{}
	case extensionsv1.SchemeGroupVersion.WithKind("TestType"):
		return &applyconfigurationextensionsv1.TestTypeApplyConfiguration{}
	case extensionsv1.SchemeGroupVersion.WithKind("TestTypeStatus"):
		return &applyconfigurationextensionsv1.TestTypeStatusApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
