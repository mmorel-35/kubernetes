/*
Copyright 2024 The Kubernetes Authors.

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

package protobuf

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"k8s.io/gengo/v2/types"
)

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

func TestWriteTagCode(t *testing.T) {
	tests := []struct {
		name     string
		fieldNum int
		wireType int
		wantByte string // the hex byte that should appear in the generated code
	}{
		{name: "field1_varint", fieldNum: 1, wireType: 0, wantByte: "0x08"},
		{name: "field1_ld", fieldNum: 1, wireType: 2, wantByte: "0x0a"},
		{name: "field2_varint", fieldNum: 2, wireType: 0, wantByte: "0x10"},
		{name: "field2_ld", fieldNum: 2, wireType: 2, wantByte: "0x12"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code := writeTagCode(tc.fieldNum, tc.wireType)
			if !strings.Contains(code, tc.wantByte) {
				t.Errorf("writeTagCode(%d,%d) = %q, want to contain %q", tc.fieldNum, tc.wireType, code, tc.wantByte)
			}
		})
	}
}

func TestCastTypeComponents(t *testing.T) {
	tests := []struct {
		input     string
		wantPkg   string
		wantType  string
		wantAlias string
	}{
		{
			input:     "k8s.io/apimachinery/pkg/types.UID",
			wantPkg:   "k8s.io/apimachinery/pkg/types",
			wantType:  "UID",
			wantAlias: "k8s_io_apimachinery_pkg_types",
		},
		{
			input:     "StatusReason",
			wantPkg:   "",
			wantType:  "StatusReason",
			wantAlias: "",
		},
		{
			input:     "k8s.io/api/core/v1.ResourceName",
			wantPkg:   "k8s.io/api/core/v1",
			wantType:  "ResourceName",
			wantAlias: "k8s_io_api_core_v1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			gotPkg, gotType, gotAlias := castTypeComponents(tc.input)
			if gotPkg != tc.wantPkg {
				t.Errorf("pkg: got %q, want %q", gotPkg, tc.wantPkg)
			}
			if gotType != tc.wantType {
				t.Errorf("type: got %q, want %q", gotType, tc.wantType)
			}
			if gotAlias != tc.wantAlias {
				t.Errorf("alias: got %q, want %q", gotAlias, tc.wantAlias)
			}
		})
	}
}

func TestProtoToGoType(t *testing.T) {
	tests := map[string]string{
		"string":   "string",
		"bytes":    "[]byte",
		"bool":     "bool",
		"int32":    "int32",
		"int64":    "int64",
		"uint32":   "uint32",
		"uint64":   "uint64",
		"sint32":   "int32",
		"sint64":   "int64",
		"fixed32":  "uint32",
		"fixed64":  "uint64",
		"sfixed32": "int32",
		"sfixed64": "int64",
		"double":   "float64",
		"float":    "float32",
	}
	for proto, want := range tests {
		got := protoToGoType(proto)
		if got != want {
			t.Errorf("protoToGoType(%q) = %q, want %q", proto, got, want)
		}
	}
}

func TestGoImportAlias(t *testing.T) {
	tests := map[string]string{
		"k8s.io/apimachinery/pkg/types": "k8s_io_apimachinery_pkg_types",
		"fmt": "fmt",
		"some/package-name": "some_package_name",
	}
	for input, want := range tests {
		got := goImportAlias(input)
		if got != want {
			t.Errorf("goImportAlias(%q) = %q, want %q", input, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Imports: sorted output
// ---------------------------------------------------------------------------

func TestImportsSorted(t *testing.T) {
	g := &genGoMarshal{
		neededImports: map[string]string{
			"fmt":       "fmt",
			"io":        "io",
			"sort":      "sort",
			"math/bits": "math_bits",
			"strings":   "strings",
		},
	}
	got := g.Imports(nil)
	if !sort.StringsAreSorted(got) {
		t.Errorf("Imports() returned unsorted list: %v", got)
	}
}

// ---------------------------------------------------------------------------
// emitMarshalMapField: key collection and encoding
// ---------------------------------------------------------------------------

// newMapField builds a minimal protoField representing a map<keyProtoType, string> field.
func newMapField(fieldName, keyProtoName string) *protoField {
	return &protoField{
		Tag:  3,
		Name: fieldName,
		Map:  true,
		Type: &types.Type{
			Key: &types.Type{
				Name: types.Name{Name: keyProtoName},
			},
			Elem: &types.Type{
				Name: types.Name{Name: "string"},
			},
		},
		Extras: map[string]string{},
	}
}

func newGenGoMarshal() *genGoMarshal {
	return &genGoMarshal{
		neededImports: map[string]string{},
	}
}

func TestEmitMarshalMapField_StringKey(t *testing.T) {
	g := newGenGoMarshal()
	f := newMapField("Labels", "string")

	var buf bytes.Buffer
	tagCode := writeTagCode(f.Tag, 2)
	g.emitMarshalMapField(&buf, f, "m.Labels", tagCode)
	code := buf.String()

	// Should use []string slice and sort.Strings
	if !strings.Contains(code, "make([]string, 0") {
		t.Errorf("expected []string slice for string key, got:\n%s", code)
	}
	if !strings.Contains(code, "sort.Strings(") {
		t.Errorf("expected sort.Strings for string key, got:\n%s", code)
	}
	// String key encoding: field 1, wire type 2 → 0xa
	if !strings.Contains(code, "0xa") {
		t.Errorf("expected 0xa tag byte for string key, got:\n%s", code)
	}
	// No integer conversion of keyExpr
	if strings.Contains(code, "sort.Slice") {
		t.Errorf("unexpected sort.Slice for string key, got:\n%s", code)
	}
}

func TestEmitMarshalMapField_Int32Key(t *testing.T) {
	g := newGenGoMarshal()
	f := newMapField("Counters", "int32")

	var buf bytes.Buffer
	tagCode := writeTagCode(f.Tag, 2)
	g.emitMarshalMapField(&buf, f, "m.Counters", tagCode)
	code := buf.String()

	// Should use []int32 slice and sort.Slice
	if !strings.Contains(code, "make([]int32, 0") {
		t.Errorf("expected []int32 slice for int32 key, got:\n%s", code)
	}
	if !strings.Contains(code, "sort.Slice(") {
		t.Errorf("expected sort.Slice for int32 key, got:\n%s", code)
	}
	// Integer key encoding: field 1, wire type 0 → 0x08
	if !strings.Contains(code, "0x8") {
		t.Errorf("expected 0x8 tag byte for varint key, got:\n%s", code)
	}
	// No sort.Strings
	if strings.Contains(code, "sort.Strings") {
		t.Errorf("unexpected sort.Strings for int32 key, got:\n%s", code)
	}
}

func TestEmitMarshalMapField_BoolKey(t *testing.T) {
	g := newGenGoMarshal()
	f := newMapField("Flags", "bool")

	var buf bytes.Buffer
	tagCode := writeTagCode(f.Tag, 2)
	g.emitMarshalMapField(&buf, f, "m.Flags", tagCode)
	code := buf.String()

	// Should use []bool slice and sort.Slice
	if !strings.Contains(code, "make([]bool, 0") {
		t.Errorf("expected []bool slice for bool key, got:\n%s", code)
	}
	if !strings.Contains(code, "sort.Slice(") {
		t.Errorf("expected sort.Slice for bool key, got:\n%s", code)
	}
	// Bool key encoding: field 1, wire type 0 → 0x08
	if !strings.Contains(code, "0x8") {
		t.Errorf("expected 0x8 tag byte for bool key, got:\n%s", code)
	}
}

func TestEmitMarshalMapField_CastKey(t *testing.T) {
	g := newGenGoMarshal()
	f := newMapField("Labels", "string")
	// Simulate (gogoproto.castkey) = "k8s.io/api/core/v1.ResourceName"
	f.Extras["(gogoproto.castkey)"] = `"k8s.io/api/core/v1.ResourceName"`

	var buf bytes.Buffer
	tagCode := writeTagCode(f.Tag, 2)
	g.emitMarshalMapField(&buf, f, "m.Labels", tagCode)
	code := buf.String()

	// Key collection still uses []string (proto key is still "string")
	if !strings.Contains(code, "make([]string, 0") {
		t.Errorf("expected []string slice for castkey with string proto, got:\n%s", code)
	}
	// Map lookup should apply the cast
	if !strings.Contains(code, "ResourceName(") {
		t.Errorf("expected ResourceName() cast in map lookup, got:\n%s", code)
	}
	// The castkey package should have been added to imports
	if _, ok := g.neededImports["k8s.io/api/core/v1"]; !ok {
		t.Errorf("expected castkey import to be added, got imports: %v", g.neededImports)
	}
}

// ---------------------------------------------------------------------------
// emitSizeMapField: key size computation
// ---------------------------------------------------------------------------

func TestEmitSizeMapField_StringKey(t *testing.T) {
	g := newGenGoMarshal()
	f := newMapField("Labels", "string")

	var buf bytes.Buffer
	g.emitSizeMapField(&buf, f, "m.Labels", "1")
	code := buf.String()

	// String key size: len(k) + sovGenerated(uint64(len(k)))
	if !strings.Contains(code, "len(k)") {
		t.Errorf("expected len(k) for string key size, got:\n%s", code)
	}
	if strings.Contains(code, "sovGenerated(uint64(k))") {
		t.Errorf("unexpected varint key size for string key, got:\n%s", code)
	}
}

func TestEmitSizeMapField_VarintKey(t *testing.T) {
	g := newGenGoMarshal()
	f := newMapField("Counters", "int32")

	var buf bytes.Buffer
	g.emitSizeMapField(&buf, f, "m.Counters", "1")
	code := buf.String()

	// Varint key size: sovGenerated(uint64(k))
	if !strings.Contains(code, "sovGenerated(uint64(k))") {
		t.Errorf("expected sovGenerated(uint64(k)) for varint key size, got:\n%s", code)
	}
	if strings.Contains(code, "len(k)") {
		t.Errorf("unexpected len(k) for varint key, got:\n%s", code)
	}
}

func TestEmitSizeMapField_BoolKey(t *testing.T) {
	g := newGenGoMarshal()
	f := newMapField("Flags", "bool")

	var buf bytes.Buffer
	g.emitSizeMapField(&buf, f, "m.Flags", "1")
	code := buf.String()

	// Bool key size: 1 + 1 (tag + value)
	if !strings.Contains(code, "mapEntrySize := 1 + 1") {
		t.Errorf("expected 'mapEntrySize := 1 + 1' for bool key size, got:\n%s", code)
	}
	if strings.Contains(code, "len(k)") {
		t.Errorf("unexpected len(k) for bool key, got:\n%s", code)
	}
}

func TestEmitSizeMapField_MessageValue(t *testing.T) {
	g := newGenGoMarshal()
	// Map with string key and message value (e.g., map[string]metav1.Time)
	f := &protoField{
		Tag:  5,
		Name: "DisruptedPods",
		Map:  true,
		Type: &types.Type{
			Key: &types.Type{
				Name: types.Name{Name: "string"},
			},
			Elem: &types.Type{
				// Non-primitive name with a package = message type
				Name: types.Name{
					Name:    "Time",
					Package: "metav1",
					Path:    "k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto",
				},
			},
		},
		Extras: map[string]string{},
	}

	var buf bytes.Buffer
	g.emitSizeMapField(&buf, f, "m.DisruptedPods", "1")
	code := buf.String()

	// v.Size() should be called once and stored in l
	if !strings.Contains(code, "l = v.Size()") {
		t.Errorf("expected 'l = v.Size()' for message value, got:\n%s", code)
	}
	// Should use l (not v.Size()) in mapEntrySize expression
	if strings.Contains(code, "v.Size() + sovGenerated") {
		t.Errorf("unexpected duplicate v.Size() call, got:\n%s", code)
	}
}

// ---------------------------------------------------------------------------
// emitMarshalMapField: message value (addressability fix)
// ---------------------------------------------------------------------------

func TestEmitMarshalMapField_MessageValue(t *testing.T) {
	g := newGenGoMarshal()
	// Map with string key and message value (e.g., map[string]metav1.Time)
	f := &protoField{
		Tag:  5,
		Name: "DisruptedPods",
		Map:  true,
		Type: &types.Type{
			Key: &types.Type{
				Name: types.Name{Name: "string"},
			},
			Elem: &types.Type{
				Name: types.Name{
					Name:    "Time",
					Package: "metav1",
					Path:    "k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto",
				},
			},
		},
		Extras: map[string]string{},
	}

	var buf bytes.Buffer
	tagCode := writeTagCode(f.Tag, 2)
	g.emitMarshalMapField(&buf, f, "m.DisruptedPods", tagCode)
	code := buf.String()

	// Must assign map value to a local variable before baseI to avoid
	// "cannot take address of map index expression".
	if !strings.Contains(code, "v :=") {
		t.Errorf("expected local variable assignment 'v :=' for map message value, got:\n%s", code)
	}
	// Must use (&v) for pointer receiver
	if !strings.Contains(code, "(&v).MarshalToSizedBuffer") {
		t.Errorf("expected '(&v).MarshalToSizedBuffer' for map message value, got:\n%s", code)
	}
	// The v := assignment must come before baseI
	vIdx := strings.Index(code, "v :=")
	baseIdx := strings.Index(code, "baseI := i")
	if vIdx < 0 || baseIdx < 0 || vIdx > baseIdx {
		t.Errorf("'v :=' must appear before 'baseI := i', got:\n%s", code)
	}
}

// ---------------------------------------------------------------------------
// emitUnmarshalMapField: cross-package message value type aliasing
// ---------------------------------------------------------------------------

func TestEmitUnmarshalMapField_CrossPackageMessageValue(t *testing.T) {
	g := newGenGoMarshal()
	// map[string]metav1.Time — value type is from another package.
	f := &protoField{
		Tag:  5,
		Name: "DisruptedPods",
		Map:  true,
		Type: &types.Type{
			Key: &types.Type{
				Name: types.Name{Name: "string"},
			},
			Elem: &types.Type{
				Name: types.Name{
					Name:    "Time",
					Package: "metav1",
					Path:    "k8s.io/apimachinery/pkg/apis/meta/v1/generated.proto",
				},
			},
		},
		Extras: map[string]string{},
	}

	var buf bytes.Buffer
	g.emitUnmarshalMapField(&buf, f, "m.DisruptedPods", "MyStruct", "")
	code := buf.String()

	// The generated code must qualify the value type with the import alias.
	// For path "k8s.io/apimachinery/pkg/apis/meta/v1" the alias is
	// "k8s_io_apimachinery_pkg_apis_meta_v1".
	wantType := "k8s_io_apimachinery_pkg_apis_meta_v1.Time{}"
	if !strings.Contains(code, wantType) {
		t.Errorf("expected qualified type %q in unmarshal map field code, got:\n%s", wantType, code)
	}

	// The import must also have been recorded.
	wantPkg := "k8s.io/apimachinery/pkg/apis/meta/v1"
	if _, ok := g.neededImports[wantPkg]; !ok {
		t.Errorf("expected import %q to be recorded, neededImports=%v", wantPkg, g.neededImports)
	}
}
