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
	"fmt"
	"io"
	"strconv"
	"strings"

	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/types"
)

// genGoMarshal generates the generated.pb.go file directly from Go type
// information, without going through protoc-gen-gogo. This removes the
// gogo/protobuf dependency from the code-generator toolchain.
//
// The generated code is wire-format compatible with the previous gogo output:
// it implements the same Marshal/MarshalTo/MarshalToSizedBuffer/Size/Unmarshal
// interface consumed by k8s.io/apimachinery's runtime.ProtobufMarshaller.
type genGoMarshal struct {
	generator.GoGenerator
	localPackage   types.Name
	localGoPackage types.Name
	protoImports   *ImportTracker

	generateAll    bool
	omitFieldTypes map[types.Name]struct{}

	// neededImports tracks Go import path -> alias for the generated file.
	// Populated during GenerateType, consumed by Imports().
	neededImports map[string]string

	// typesWithValueStringMethod holds types that already define a String()
	// method with a value (non-pointer) receiver, so we can emit a separate
	// pointer-receiver String() without conflict.
	typesWithValueStringMethod map[string]bool
}

func (g *genGoMarshal) Name() string     { return "go-marshal" }
func (g *genGoMarshal) Filename() string { return "generated.pb.go" }
func (g *genGoMarshal) FileType() string { return "go" }

func (g *genGoMarshal) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"local": localNamer{g.localPackage},
	}
}

// Filter accepts the same types as genProtoIDL.
func (g *genGoMarshal) Filter(c *generator.Context, t *types.Type) bool {
	return (&genProtoIDL{
		localPackage:   g.localPackage,
		localGoPackage: g.localGoPackage,
		imports:        g.protoImports,
		generateAll:    g.generateAll,
		omitFieldTypes: g.omitFieldTypes,
	}).Filter(c, t)
}

// Init pre-scans types to discover which ones already have a String() method
// with a value receiver (so we don't create duplicate definitions).
func (g *genGoMarshal) Init(c *generator.Context, w io.Writer) error {
	g.neededImports = map[string]string{
		"fmt":        "fmt",
		"io":         "io",
		"math/bits":  "math_bits",
		"strings":    "strings",
		"reflect":    "reflect",
	}
	g.typesWithValueStringMethod = make(map[string]bool)

	pkg, ok := c.Universe[g.localGoPackage.Package]
	if !ok {
		return nil
	}
	for _, t := range pkg.Types {
		if m, ok := t.Methods["String"]; ok {
			// Value receiver means the method signature is `func (T) String() string`.
			if m.Signature != nil && len(m.Signature.Parameters) == 0 {
				// Check receiver kind – value vs pointer is indicated by the
				// receiver type being the same as the declaring type (not a pointer).
				// In gengo types, if the method is on a value receiver the
				// receiver type kind != types.Pointer.
				if m.Signature.Receiver != nil && m.Signature.Receiver.Kind != types.Pointer {
					g.typesWithValueStringMethod[t.Name.Name] = true
				}
			}
		}
	}
	return nil
}

// Imports returns the Go imports needed by the generated file.
// Called after all GenerateType calls, so neededImports is complete.
func (g *genGoMarshal) Imports(c *generator.Context) []string {
	out := make([]string, 0, len(g.neededImports))
	for path, alias := range g.neededImports {
		base := path
		if idx := strings.LastIndex(path, "/"); idx >= 0 {
			base = path[idx+1:]
		}
		if base == alias {
			out = append(out, strconv.Quote(path))
		} else {
			out = append(out, fmt.Sprintf("%s %s", alias, strconv.Quote(path)))
		}
	}
	return out
}

// GenerateType emits all marshal/unmarshal methods for a single type.
func (g *genGoMarshal) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	locator := &protobufLocator{
		namer:          c.Namers["proto"].(ProtobufFromGoNamer),
		tracker:        g.protoImports,
		universe:       c.Universe,
		localGoPackage: g.localGoPackage.Package,
	}

	switch t.Kind {
	case types.Struct:
		return g.generateForStruct(w, locator, t)
	case types.Alias:
		if isOptionalAlias(t) {
			return g.generateForOptionalAlias(w, locator, t)
		}
	}
	return nil
}

// Finalize writes the package-level shared helpers at the end of the file.
func (g *genGoMarshal) Finalize(c *generator.Context, w io.Writer) error {
	io.WriteString(w, sharedHelpersCode) //nolint:errcheck
	return nil
}

// ---------------------------------------------------------------------------
// Struct type generation
// ---------------------------------------------------------------------------

func (g *genGoMarshal) generateForStruct(w io.Writer, locator ProtobufLocator, t *types.Type) error {
	typeName := t.Name.Name
	fields, err := membersToFields(locator, t, g.localPackage, g.omitFieldTypes)
	if err != nil {
		return fmt.Errorf("unable to get fields for %s: %v", typeName, err)
	}

	g.emitReset(w, typeName, false)
	g.emitMarshal(w, typeName, false)
	g.emitMarshalTo(w, typeName, false)
	g.emitMarshalToSizedBuffer(w, typeName, fields, false)
	g.emitSize(w, typeName, fields, false)
	g.emitUnmarshal(w, typeName, fields, false)
	g.emitString(w, typeName, fields, false)
	return nil
}

// ---------------------------------------------------------------------------
// Optional alias type generation (e.g. type Verbs []string)
// ---------------------------------------------------------------------------

func (g *genGoMarshal) generateForOptionalAlias(w io.Writer, locator ProtobufLocator, t *types.Type) error {
	typeName := t.Name.Name

	// Build a synthetic protoField representing the single repeated/map field
	// that wraps the alias's underlying type.
	field := protoField{
		LocalPackage: g.localPackage,
		Tag:          1,
		Name:         "items",
		Extras:       make(map[string]string),
	}
	if err := memberTypeToProtobufField(locator, &field, t.Underlying); err != nil {
		return fmt.Errorf("unable to determine proto type for alias %s: %v", typeName, err)
	}
	// For optional aliases the items field is repeated or map.
	// memberTypeToProtobufField will have set Repeated/Map already for slices/maps.

	g.emitReset(w, typeName, true)
	g.emitMarshalOptionalAlias(w, typeName, &field)
	g.emitSizeOptionalAlias(w, typeName, &field)
	g.emitUnmarshalOptionalAlias(w, typeName, &field)
	// Only emit String() if the type doesn't already define one. Check by
	// looking at both the genng-scanned methods and the Init pre-scan.
	_, typeDefinesStringMethod := t.Methods["String"]
	if !typeDefinesStringMethod && !g.typesWithValueStringMethod[typeName] {
		g.emitString(w, typeName, nil, true)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Reset
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitReset(w io.Writer, typeName string, _ bool) {
	fmt.Fprintf(w, "func (m *%s) Reset() { *m = %s{} }\n\n", typeName, typeName)
}

// ---------------------------------------------------------------------------
// Marshal / MarshalTo
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitMarshal(w io.Writer, typeName string, optional bool) {
	recv := receiverExpr(typeName, optional, false)
	fmt.Fprintf(w, `func (%s) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

`, recv)
}

func (g *genGoMarshal) emitMarshalTo(w io.Writer, typeName string, optional bool) {
	recv := receiverExpr(typeName, optional, false)
	fmt.Fprintf(w, `func (%s) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

`, recv)
}

// ---------------------------------------------------------------------------
// MarshalToSizedBuffer (struct)
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitMarshalToSizedBuffer(w io.Writer, typeName string, fields []protoField, optional bool) {
	recv := receiverExpr(typeName, optional, false)
	fmt.Fprintf(w, "func (%s) MarshalToSizedBuffer(dAtA []byte) (int, error) {\n", recv)
	fmt.Fprint(w, "\ti := len(dAtA)\n\t_ = i\n\tvar l int\n\t_ = l\n")

	// Fields are written in reverse order.
	for fi := len(fields) - 1; fi >= 0; fi-- {
		f := &fields[fi]
		g.emitMarshalField(w, f)
	}
	fmt.Fprint(w, "\treturn len(dAtA) - i, nil\n}\n\n")
}

// emitMarshalField generates the marshal code for a single field.
func (g *genGoMarshal) emitMarshalField(w io.Writer, f *protoField) {
	goName := g.goFieldName(f)
	fieldAccess := "m." + goName
	tagCode := writeTagCode(f.Tag, protoWireType(f))

	if f.Map {
		g.emitMarshalMapField(w, f, fieldAccess, tagCode)
		return
	}

	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package

	switch {
	case f.Repeated:
		g.emitMarshalRepeatedField(w, f, fieldAccess, tagCode)

	case isMessageType(f):
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n", fieldAccess)
			fmt.Fprintf(w, "\t\t{\n\t\t\tsize, err := %s.MarshalToSizedBuffer(dAtA[:i])\n\t\t\tif err != nil {\n\t\t\t\treturn 0, err\n\t\t\t}\n\t\t\ti -= size\n\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(size))\n\t\t}\n", fieldAccess)
			fmt.Fprintf(w, "%s\n\t}\n", tagCode)
		} else {
			fmt.Fprintf(w, "\t{\n\t\tsize, err := %s.MarshalToSizedBuffer(dAtA[:i])\n\t\tif err != nil {\n\t\t\treturn 0, err\n\t\t}\n\t\ti -= size\n\t\ti = encodeVarintGenerated(dAtA, i, uint64(size))\n\t}\n", fieldAccess)
			fmt.Fprintf(w, "%s\n", tagCode)
		}

	case protoName == "bool" && pkg == "":
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n", fieldAccess)
			fmt.Fprintf(w, "\t\ti--\n\t\tif *%s {\n\t\t\tdAtA[i] = 1\n\t\t} else {\n\t\t\tdAtA[i] = 0\n\t\t}\n", fieldAccess)
			fmt.Fprintf(w, "%s\n\t}\n", tagCode)
		} else {
			fmt.Fprintf(w, "\ti--\n\tif %s {\n\t\tdAtA[i] = 1\n\t} else {\n\t\tdAtA[i] = 0\n\t}\n", fieldAccess)
			fmt.Fprintf(w, "%s\n", tagCode)
		}

	case isVarintType(protoName) && pkg == "":
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n", fieldAccess)
			fmt.Fprintf(w, "\t\ti = encodeVarintGenerated(dAtA, i, uint64(*%s))\n", fieldAccess)
			fmt.Fprintf(w, "%s\n\t}\n", tagCode)
		} else {
			fmt.Fprintf(w, "\ti = encodeVarintGenerated(dAtA, i, uint64(%s))\n", fieldAccess)
			fmt.Fprintf(w, "%s\n", tagCode)
		}

	case protoName == "double" && pkg == "":
		g.addImport("encoding/binary", "encoding_binary")
		g.addImport("math", "math")
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n", fieldAccess)
			fmt.Fprintf(w, "\t\ti -= 8\n\t\tencoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(*%s))))\n", fieldAccess)
			fmt.Fprintf(w, "%s\n\t}\n", tagCode)
		} else {
			fmt.Fprintf(w, "\ti -= 8\n\tencoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(%s))))\n", fieldAccess)
			fmt.Fprintf(w, "%s\n", tagCode)
		}

	case protoName == "float" && pkg == "":
		g.addImport("encoding/binary", "encoding_binary")
		g.addImport("math", "math")
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n", fieldAccess)
			fmt.Fprintf(w, "\t\ti -= 4\n\t\tencoding_binary.LittleEndian.PutUint32(dAtA[i:], math.Float32bits(float32(*%s)))\n", fieldAccess)
			fmt.Fprintf(w, "%s\n\t}\n", tagCode)
		} else {
			fmt.Fprintf(w, "\ti -= 4\n\tencoding_binary.LittleEndian.PutUint32(dAtA[i:], math.Float32bits(float32(%s)))\n", fieldAccess)
			fmt.Fprintf(w, "%s\n", tagCode)
		}

	default:
		// string, bytes, or cast-type string field (wire type 2)
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n", fieldAccess)
			fmt.Fprintf(w, "\t\ti -= len(*%s)\n\t\tcopy(dAtA[i:], *%s)\n\t\ti = encodeVarintGenerated(dAtA, i, uint64(len(*%s)))\n", fieldAccess, fieldAccess, fieldAccess)
			fmt.Fprintf(w, "%s\n\t}\n", tagCode)
		} else {
			fmt.Fprintf(w, "\ti -= len(%s)\n\tcopy(dAtA[i:], %s)\n\ti = encodeVarintGenerated(dAtA, i, uint64(len(%s)))\n", fieldAccess, fieldAccess, fieldAccess)
			fmt.Fprintf(w, "%s\n", tagCode)
		}
	}
}

func (g *genGoMarshal) emitMarshalRepeatedField(w io.Writer, f *protoField, fieldAccess, tagCode string) {
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package

	fmt.Fprintf(w, "\tif len(%s) > 0 {\n", fieldAccess)
	fmt.Fprintf(w, "\t\tfor iNdEx := len(%s) - 1; iNdEx >= 0; iNdEx-- {\n", fieldAccess)

	elem := fieldAccess + "[iNdEx]"
	switch {
	case isMessageType(f):
		fmt.Fprintf(w, "\t\t\t{\n\t\t\t\tsize, err := %s.MarshalToSizedBuffer(dAtA[:i])\n\t\t\t\tif err != nil {\n\t\t\t\t\treturn 0, err\n\t\t\t\t}\n\t\t\t\ti -= size\n\t\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(size))\n\t\t\t}\n", elem)
	case isVarintType(protoName) && pkg == "":
		fmt.Fprintf(w, "\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(%s))\n", elem)
	case protoName == "bool" && pkg == "":
		fmt.Fprintf(w, "\t\t\ti--\n\t\t\tif %s {\n\t\t\t\tdAtA[i] = 1\n\t\t\t} else {\n\t\t\t\tdAtA[i] = 0\n\t\t\t}\n", elem)
	default:
		// string / bytes
		fmt.Fprintf(w, "\t\t\ti -= len(%s)\n\t\t\tcopy(dAtA[i:], %s)\n\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(len(%s)))\n", elem, elem, elem)
	}
	fmt.Fprintf(w, "%s\n", indentTagCode(tagCode, "\t\t\t"))
	fmt.Fprint(w, "\t\t}\n\t}\n")
}

func (g *genGoMarshal) emitMarshalMapField(w io.Writer, f *protoField, fieldAccess, tagCode string) {
	g.addImport("sort", "sort")

	castKey := unquote(f.Extras["(gogoproto.castkey)"])
	castVal := unquote(f.Extras["(gogoproto.castvalue)"])

	keyProtoName := f.Type.Key.Name.Name
	valProtoName := f.Type.Elem.Name.Name
	valPkg := f.Type.Elem.Name.Package

	// Sort keys for deterministic output.
	keysVar := "keysFor" + g.goFieldName(f)
	fmt.Fprintf(w, "\tif len(%s) > 0 {\n", fieldAccess)
	fmt.Fprintf(w, "\t\t%s := make([]string, 0, len(%s))\n", keysVar, fieldAccess)
	fmt.Fprintf(w, "\t\tfor k := range %s {\n\t\t\t%s = append(%s, string(k))\n\t\t}\n", fieldAccess, keysVar, keysVar)
	fmt.Fprintf(w, "\t\tsort.Strings(%s)\n", keysVar)
	fmt.Fprintf(w, "\t\tfor iNdEx := len(%s) - 1; iNdEx >= 0; iNdEx-- {\n", keysVar)

	// Reconstruct the key access with possible cast.
	keyExpr := keysVar + "[iNdEx]"
	mapKeyExpr := keyExpr
	if castKey != "" {
		castKeyPkg, castKeyType, castKeyAlias := castTypeComponents(castKey)
		if castKeyPkg != "" {
			g.addImport(castKeyPkg, castKeyAlias)
			mapKeyExpr = castKeyAlias + "." + castKeyType + "(" + keyExpr + ")"
		} else {
			mapKeyExpr = castKeyType + "(" + keyExpr + ")"
		}
	} else if keyProtoName != "string" {
		mapKeyExpr = keyProtoName + "(" + keyExpr + ")"
	}

	valAccess := fieldAccess + "[" + mapKeyExpr + "]"

	// For message values, extract to a local variable because map index
	// expressions are not addressable (needed for pointer-receiver MarshalToSizedBuffer).
	if isProtoMessageType(f.Type.Elem) {
		fmt.Fprintf(w, "\t\t\tv := %s\n", valAccess)
	}
	fmt.Fprint(w, "\t\t\tbaseI := i\n")
	// Write value first (field 2), then key (field 1), then total length.
	switch {
	case isProtoMessageType(f.Type.Elem):
		g.emitMapValueMessage(w, "&v", castVal)
	case isVarintType(valProtoName) && valPkg == "":
		fmt.Fprintf(w, "\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(%s))\n", valAccess)
		fmt.Fprint(w, "\t\t\ti--\n\t\t\tdAtA[i] = 0x10\n") // field 2, wire type 0
	case valProtoName == "bool" && valPkg == "":
		fmt.Fprintf(w, "\t\t\ti--\n\t\t\tif %s {\n\t\t\t\tdAtA[i] = 1\n\t\t\t} else {\n\t\t\t\tdAtA[i] = 0\n\t\t\t}\n", valAccess)
		fmt.Fprint(w, "\t\t\ti--\n\t\t\tdAtA[i] = 0x10\n") // field 2, wire type 0
	default:
		// string/bytes value
		fmt.Fprintf(w, "\t\t\ti -= len(%s)\n\t\t\tcopy(dAtA[i:], %s)\n\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(len(%s)))\n", valAccess, valAccess, valAccess)
		fmt.Fprint(w, "\t\t\ti--\n\t\t\tdAtA[i] = 0x12\n") // field 2, wire type 2
	}
	// Write key (field 1).
	switch {
	case isVarintType(keyProtoName):
		fmt.Fprintf(w, "\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(%s))\n", keyExpr)
		fmt.Fprint(w, "\t\t\ti--\n\t\t\tdAtA[i] = 0x8\n") // field 1, wire type 0
	default:
		// string key
		fmt.Fprintf(w, "\t\t\ti -= len(%s)\n\t\t\tcopy(dAtA[i:], %s)\n\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(len(%s)))\n", keyExpr, keyExpr, keyExpr)
		fmt.Fprint(w, "\t\t\ti--\n\t\t\tdAtA[i] = 0xa\n") // field 1, wire type 2
	}
	// Write the map entry outer length.
	fmt.Fprint(w, "\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(baseI-i))\n")
	// Write the map entry tag.
	fmt.Fprintf(w, "%s\n", indentTagCode(tagCode, "\t\t\t"))
	fmt.Fprint(w, "\t\t}\n\t}\n")
}

func (g *genGoMarshal) emitMapValueMessage(w io.Writer, valAccess, castVal string) {
	fmt.Fprintf(w, "\t\t\t{\n\t\t\t\tsize, err := %s.MarshalToSizedBuffer(dAtA[:i])\n\t\t\t\tif err != nil {\n\t\t\t\t\treturn 0, err\n\t\t\t\t}\n\t\t\t\ti -= size\n\t\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(size))\n\t\t\t}\n", valAccess)
	fmt.Fprint(w, "\t\t\ti--\n\t\t\tdAtA[i] = 0x12\n") // field 2, wire type 2
}

// ---------------------------------------------------------------------------
// MarshalToSizedBuffer for optional alias
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitMarshalOptionalAlias(w io.Writer, typeName string, f *protoField) {
	tagCode := writeTagCode(1, protoWireType(f))
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package

	// Marshal uses a value receiver for optional aliases.
	fmt.Fprintf(w, "func (m %s) Marshal() (dAtA []byte, err error) {\n\tsize := m.Size()\n\tdAtA = make([]byte, size)\n\tn, err := m.MarshalToSizedBuffer(dAtA[:size])\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn dAtA[:n], nil\n}\n\n", typeName)
	fmt.Fprintf(w, "func (m %s) MarshalTo(dAtA []byte) (int, error) {\n\tsize := m.Size()\n\treturn m.MarshalToSizedBuffer(dAtA[:size])\n}\n\n", typeName)

	fmt.Fprintf(w, "func (m %s) MarshalToSizedBuffer(dAtA []byte) (int, error) {\n", typeName)
	fmt.Fprint(w, "\ti := len(dAtA)\n\t_ = i\n\tvar l int\n\t_ = l\n")
	fmt.Fprint(w, "\tif len(m) > 0 {\n")
	fmt.Fprint(w, "\t\tfor iNdEx := len(m) - 1; iNdEx >= 0; iNdEx-- {\n")

	elem := "m[iNdEx]"
	switch {
	case isMessageType(f):
		fmt.Fprintf(w, "\t\t\t{\n\t\t\t\tsize, err := %s.MarshalToSizedBuffer(dAtA[:i])\n\t\t\t\tif err != nil {\n\t\t\t\t\treturn 0, err\n\t\t\t\t}\n\t\t\t\ti -= size\n\t\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(size))\n\t\t\t}\n", elem)
	case isVarintType(protoName) && pkg == "":
		fmt.Fprintf(w, "\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(%s))\n", elem)
	default:
		// string
		fmt.Fprintf(w, "\t\t\ti -= len(%s)\n\t\t\tcopy(dAtA[i:], %s)\n\t\t\ti = encodeVarintGenerated(dAtA, i, uint64(len(%s)))\n", elem, elem, elem)
	}
	fmt.Fprintf(w, "%s\n", indentTagCode(tagCode, "\t\t\t"))
	fmt.Fprint(w, "\t\t}\n\t}\n")
	fmt.Fprint(w, "\treturn len(dAtA) - i, nil\n}\n\n")
}

// ---------------------------------------------------------------------------
// Size (struct)
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitSize(w io.Writer, typeName string, fields []protoField, optional bool) {
	recv := receiverExpr(typeName, optional, false)
	fmt.Fprintf(w, "func (%s) Size() (n int) {\n", recv)
	fmt.Fprint(w, "\tif m == nil {\n\t\treturn 0\n\t}\n")
	fmt.Fprint(w, "\tvar l int\n\t_ = l\n")
	for _, f := range fields {
		g.emitSizeField(w, &f)
	}
	fmt.Fprint(w, "\treturn n\n}\n\n")
}

func (g *genGoMarshal) emitSizeField(w io.Writer, f *protoField) {
	goName := g.goFieldName(f)
	fieldAccess := "m." + goName
	ts := tagSizeExpr(f.Tag, protoWireType(f))
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package

	if f.Map {
		g.emitSizeMapField(w, f, fieldAccess, ts)
		return
	}

	switch {
	case f.Repeated && isMessageType(f):
		fmt.Fprintf(w, "\tif len(%s) > 0 {\n", fieldAccess)
		fmt.Fprintf(w, "\t\tfor _, e := range %s {\n", fieldAccess)
		fmt.Fprintf(w, "\t\t\tl = e.Size()\n\t\t\tn += %s + l + sovGenerated(uint64(l))\n\t\t}\n\t}\n", ts)

	case f.Repeated && (protoName == "string" || protoName == "bytes") && pkg == "":
		fmt.Fprintf(w, "\tif len(%s) > 0 {\n", fieldAccess)
		fmt.Fprintf(w, "\t\tfor _, s := range %s {\n", fieldAccess)
		fmt.Fprintf(w, "\t\t\tl = len(s)\n\t\t\tn += %s + l + sovGenerated(uint64(l))\n\t\t}\n\t}\n", ts)

	case f.Repeated && isVarintType(protoName) && pkg == "":
		fmt.Fprintf(w, "\tif len(%s) > 0 {\n", fieldAccess)
		fmt.Fprintf(w, "\t\tfor _, v := range %s {\n", fieldAccess)
		fmt.Fprintf(w, "\t\t\tn += %s + sovGenerated(uint64(v))\n\t\t}\n\t}\n", ts)

	case f.Repeated && protoName == "bool" && pkg == "":
		fmt.Fprintf(w, "\tn += %s * len(%s)\n", ts, fieldAccess)

	case f.Repeated:
		// repeated casttype string
		fmt.Fprintf(w, "\tif len(%s) > 0 {\n", fieldAccess)
		fmt.Fprintf(w, "\t\tfor _, s := range %s {\n", fieldAccess)
		fmt.Fprintf(w, "\t\t\tl = len(s)\n\t\t\tn += %s + l + sovGenerated(uint64(l))\n\t\t}\n\t}\n", ts)

	case isMessageType(f):
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n", fieldAccess)
			fmt.Fprintf(w, "\t\tl = %s.Size()\n\t\tn += %s + l + sovGenerated(uint64(l))\n\t}\n", fieldAccess, ts)
		} else {
			fmt.Fprintf(w, "\tl = %s.Size()\n\tn += %s + l + sovGenerated(uint64(l))\n", fieldAccess, ts)
		}

	case protoName == "bool" && pkg == "":
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n\t\tn += 1 + %s\n\t}\n", fieldAccess, ts)
		} else {
			fmt.Fprintf(w, "\tn += 1 + %s\n", ts)
		}

	case isVarintType(protoName) && pkg == "":
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n\t\tn += %s + sovGenerated(uint64(*%s))\n\t}\n", fieldAccess, ts, fieldAccess)
		} else {
			fmt.Fprintf(w, "\tn += %s + sovGenerated(uint64(%s))\n", ts, fieldAccess)
		}

	case protoName == "double" && pkg == "":
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n\t\tn += 8 + %s\n\t}\n", fieldAccess, ts)
		} else {
			fmt.Fprintf(w, "\tn += 8 + %s\n", ts)
		}

	case protoName == "float" && pkg == "":
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n\t\tn += 4 + %s\n\t}\n", fieldAccess, ts)
		} else {
			fmt.Fprintf(w, "\tn += 4 + %s\n", ts)
		}

	default:
		// string, bytes, or casttype string
		if f.Nullable {
			fmt.Fprintf(w, "\tif %s != nil {\n\t\tl = len(*%s)\n\t\tn += %s + l + sovGenerated(uint64(l))\n\t}\n", fieldAccess, fieldAccess, ts)
		} else {
			fmt.Fprintf(w, "\tl = len(%s)\n\tn += %s + l + sovGenerated(uint64(l))\n", fieldAccess, ts)
		}
	}
}

func (g *genGoMarshal) emitSizeMapField(w io.Writer, f *protoField, fieldAccess, ts string) {
	valProtoName := f.Type.Elem.Name.Name
	valPkg := f.Type.Elem.Name.Package

	fmt.Fprintf(w, "\tif len(%s) > 0 {\n", fieldAccess)
	fmt.Fprintf(w, "\t\tfor k, v := range %s {\n", fieldAccess)
	// Each map entry is a length-delimited message with key (field1) and value (field2).
	// Compute entry size:  1 (key tag) + sovGenerated(keyLen) + keyLen
	//                    + 1 (val tag) + sovGenerated(valLen) + valLen
	//                    + outer tag size + sovGenerated(total)
	fmt.Fprint(w, "\t\t\tmapEntrySize := 1 + len(k) + sovGenerated(uint64(len(k)))\n")
	fmt.Fprint(w, "\t\t\t_ = v\n")
	switch {
	case isProtoMessageType(f.Type.Elem):
		// Compute v.Size() once; reuse l for both the length value and sovGenerated.
		fmt.Fprint(w, "\t\t\tl = v.Size()\n")
		fmt.Fprint(w, "\t\t\tmapEntrySize += 1 + l + sovGenerated(uint64(l))\n")
	case isVarintType(valProtoName) && valPkg == "":
		fmt.Fprintf(w, "\t\t\tmapEntrySize += 1 + sovGenerated(uint64(v))\n")
	case valProtoName == "bool" && valPkg == "":
		fmt.Fprint(w, "\t\t\tmapEntrySize += 1 + 1\n")
	default:
		// string/bytes value
		fmt.Fprintf(w, "\t\t\tmapEntrySize += 1 + len(v) + sovGenerated(uint64(len(v)))\n")
	}
	fmt.Fprintf(w, "\t\t\tn += mapEntrySize + %s + sovGenerated(uint64(mapEntrySize))\n", ts)
	fmt.Fprint(w, "\t\t}\n\t}\n")
}

// ---------------------------------------------------------------------------
// Size for optional alias
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitSizeOptionalAlias(w io.Writer, typeName string, f *protoField) {
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package

	fmt.Fprintf(w, "func (m %s) Size() (n int) {\n", typeName)
	fmt.Fprint(w, "\tif m == nil {\n\t\treturn 0\n\t}\n")
	fmt.Fprint(w, "\tvar l int\n\t_ = l\n")
	fmt.Fprint(w, "\tif len(m) > 0 {\n")
	switch {
	case isMessageType(f):
		fmt.Fprint(w, "\t\tfor _, e := range m {\n\t\t\tl = e.Size()\n\t\t\tn += 1 + l + sovGenerated(uint64(l))\n\t\t}\n")
	case isVarintType(protoName) && pkg == "":
		fmt.Fprint(w, "\t\tfor _, v := range m {\n\t\t\tn += 1 + sovGenerated(uint64(v))\n\t\t}\n")
	default:
		// string
		fmt.Fprint(w, "\t\tfor _, s := range m {\n\t\t\tl = len(s)\n\t\t\tn += 1 + l + sovGenerated(uint64(l))\n\t\t}\n")
	}
	fmt.Fprint(w, "\t}\n\treturn n\n}\n\n")
}

// ---------------------------------------------------------------------------
// Unmarshal (struct)
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitUnmarshal(w io.Writer, typeName string, fields []protoField, optional bool) {
	fmt.Fprintf(w, "func (m *%s) Unmarshal(dAtA []byte) error {\n", typeName)
	fmt.Fprint(w, unmarshalPreamble)

	fmt.Fprintf(w, "\t\tif wireType == 4 {\n\t\t\treturn fmt.Errorf(\"proto: %s: wiretype end group for non-group\")\n\t\t}\n", typeName)
	fmt.Fprintf(w, "\t\tif fieldNum <= 0 {\n\t\t\treturn fmt.Errorf(\"proto: %s: illegal tag %%d (wire type %%d)\", fieldNum, wireType)\n\t\t}\n", typeName)
	fmt.Fprint(w, "\t\tswitch fieldNum {\n")

	// Build a map of tag -> field for the switch.
	for _, f := range fields {
		fmt.Fprintf(w, "\t\tcase %d:\n", f.Tag)
		g.emitUnmarshalField(w, &f, typeName)
	}

	fmt.Fprint(w, "\t\tdefault:\n")
	fmt.Fprint(w, "\t\t\tiNdEx = preIndex\n")
	fmt.Fprint(w, "\t\t\tskippy, err := skipGenerated(dAtA[iNdEx:])\n")
	fmt.Fprint(w, "\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n")
	fmt.Fprint(w, "\t\t\tif (skippy < 0) || (iNdEx+skippy) < 0 {\n\t\t\t\treturn ErrInvalidLengthGenerated\n\t\t\t}\n")
	fmt.Fprint(w, "\t\t\tif (iNdEx + skippy) > l {\n\t\t\t\treturn io.ErrUnexpectedEOF\n\t\t\t}\n")
	fmt.Fprint(w, "\t\t\tiNdEx += skippy\n")
	fmt.Fprint(w, "\t\t}\n\t}\n\n\tif iNdEx > l {\n\t\treturn io.ErrUnexpectedEOF\n\t}\n\treturn nil\n}\n\n")
}

func (g *genGoMarshal) emitUnmarshalField(w io.Writer, f *protoField, typeName string) {
	goName := g.goFieldName(f)
	fieldAccess := "m." + goName
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package
	castType := unquote(f.Extras["(gogoproto.casttype)"])

	if f.Map {
		g.emitUnmarshalMapField(w, f, fieldAccess, typeName, castType)
		return
	}

	wireType := protoWireType(f)

	switch wireType {
	case 0: // varint
		g.emitUnmarshalVarint(w, f, fieldAccess, protoName, castType)

	case 1: // 64-bit fixed
		g.addImport("encoding/binary", "encoding_binary")
		g.addImport("math", "math")
		fmt.Fprintf(w, "\t\t\tif wireType != 1 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
		fmt.Fprint(w, "\t\t\tvar v uint64\n\t\t\tif (iNdEx + 8) > l {\n\t\t\t\treturn io.ErrUnexpectedEOF\n\t\t\t}\n")
		fmt.Fprint(w, "\t\t\tv = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))\n\t\t\tiNdEx += 8\n")
		if f.Nullable {
			fmt.Fprintf(w, "\t\t\tv2 := float64(math.Float64frombits(v))\n\t\t\t%s = &v2\n", fieldAccess)
		} else {
			fmt.Fprintf(w, "\t\t\t%s = float64(math.Float64frombits(v))\n", fieldAccess)
		}

	case 5: // 32-bit fixed
		g.addImport("encoding/binary", "encoding_binary")
		g.addImport("math", "math")
		fmt.Fprintf(w, "\t\t\tif wireType != 5 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
		fmt.Fprint(w, "\t\t\tvar v uint32\n\t\t\tif (iNdEx + 4) > l {\n\t\t\t\treturn io.ErrUnexpectedEOF\n\t\t\t}\n")
		fmt.Fprint(w, "\t\t\tv = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))\n\t\t\tiNdEx += 4\n")
		if f.Nullable {
			fmt.Fprintf(w, "\t\t\tv2 := float32(math.Float32frombits(v))\n\t\t\t%s = &v2\n", fieldAccess)
		} else {
			fmt.Fprintf(w, "\t\t\t%s = float32(math.Float32frombits(v))\n", fieldAccess)
		}

	default: // wire type 2: length-delimited
		if isMessageType(f) {
			if f.Repeated {
				g.emitUnmarshalRepeatedMessage(w, f, fieldAccess, goName)
			} else {
				g.emitUnmarshalMessage(w, f, fieldAccess, goName, pkg, protoName)
			}
		} else if protoName == "bytes" && pkg == "" {
			g.emitUnmarshalBytes(w, fieldAccess, goName)
		} else if f.Repeated {
			// repeated string/casttype string
			g.emitUnmarshalRepeatedString(w, f, fieldAccess, goName, castType)
		} else {
			// string or casttype string
			g.emitUnmarshalString(w, f, fieldAccess, goName, castType)
		}
	}
}

func (g *genGoMarshal) emitUnmarshalVarint(w io.Writer, f *protoField, fieldAccess, protoName, castType string) {
	goName := g.goFieldName(f)
	fmt.Fprintf(w, "\t\t\tif wireType != 0 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
	if protoName == "bool" {
		fmt.Fprint(w, "\t\t\tvar v int\n")
		fmt.Fprint(w, unmarshalVarintLoop("v", "int"))
		if f.Nullable {
			fmt.Fprintf(w, "\t\t\tb := bool(v != 0)\n\t\t\t%s = &b\n", fieldAccess)
		} else {
			fmt.Fprintf(w, "\t\t\t%s = bool(v != 0)\n", fieldAccess)
		}
		return
	}
	// Integer type: determine the Go type name.
	goTypeName := protoToGoType(protoName)
	if castType != "" {
		castPkg, castTyp, castAlias := castTypeComponents(castType)
		if castPkg != "" {
			g.addImport(castPkg, castAlias)
			if f.Nullable {
				fmt.Fprintf(w, "\t\t\tvar v %s\n", goTypeName)
				fmt.Fprint(w, unmarshalVarintLoop("v", goTypeName))
				fmt.Fprintf(w, "\t\t\tv2 := %s.%s(v)\n\t\t\t%s = &v2\n", castAlias, castTyp, fieldAccess)
			} else {
				fmt.Fprintf(w, "\t\t\tvar v %s\n", goTypeName)
				fmt.Fprint(w, unmarshalVarintLoop("v", goTypeName))
				fmt.Fprintf(w, "\t\t\t%s = %s.%s(v)\n", fieldAccess, castAlias, castTyp)
			}
		} else {
			if f.Nullable {
				fmt.Fprintf(w, "\t\t\tvar v %s\n", goTypeName)
				fmt.Fprint(w, unmarshalVarintLoop("v", goTypeName))
				fmt.Fprintf(w, "\t\t\tv2 := %s(v)\n\t\t\t%s = &v2\n", castTyp, fieldAccess)
			} else {
				fmt.Fprintf(w, "\t\t\tvar v %s\n", goTypeName)
				fmt.Fprint(w, unmarshalVarintLoop("v", goTypeName))
				fmt.Fprintf(w, "\t\t\t%s = %s(v)\n", fieldAccess, castTyp)
			}
		}
		return
	}
	if f.Nullable {
		fmt.Fprintf(w, "\t\t\tvar v %s\n", goTypeName)
		fmt.Fprint(w, unmarshalVarintLoop("v", goTypeName))
		fmt.Fprintf(w, "\t\t\t%s = &v\n", fieldAccess)
	} else {
		fmt.Fprintf(w, "\t\t\tvar v %s\n", goTypeName)
		fmt.Fprint(w, unmarshalVarintLoop("v", goTypeName))
		fmt.Fprintf(w, "\t\t\t%s = v\n", fieldAccess)
	}
}

func (g *genGoMarshal) emitUnmarshalString(w io.Writer, f *protoField, fieldAccess, goName, castType string) {
	fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
	fmt.Fprint(w, unmarshalStringLen)
	if castType != "" {
		castPkg, castTyp, castAlias := castTypeComponents(castType)
		if castPkg != "" {
			g.addImport(castPkg, castAlias)
			if f.Nullable {
				fmt.Fprintf(w, "\t\t\ts := %s.%s(dAtA[iNdEx:postIndex])\n\t\t\t%s = &s\n", castAlias, castTyp, fieldAccess)
			} else {
				fmt.Fprintf(w, "\t\t\t%s = %s.%s(dAtA[iNdEx:postIndex])\n", fieldAccess, castAlias, castTyp)
			}
		} else {
			if f.Nullable {
				fmt.Fprintf(w, "\t\t\ts := %s(dAtA[iNdEx:postIndex])\n\t\t\t%s = &s\n", castTyp, fieldAccess)
			} else {
				fmt.Fprintf(w, "\t\t\t%s = %s(dAtA[iNdEx:postIndex])\n", fieldAccess, castTyp)
			}
		}
	} else {
		if f.Nullable {
			fmt.Fprintf(w, "\t\t\ts := string(dAtA[iNdEx:postIndex])\n\t\t\t%s = &s\n", fieldAccess)
		} else {
			fmt.Fprintf(w, "\t\t\t%s = string(dAtA[iNdEx:postIndex])\n", fieldAccess)
		}
	}
	fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
}

func (g *genGoMarshal) emitUnmarshalBytes(w io.Writer, fieldAccess, goName string) {
	fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
	fmt.Fprint(w, unmarshalBytesLen)
	fmt.Fprintf(w, "\t\t\t%s = append(%s[:0], dAtA[iNdEx:postIndex]...)\n", fieldAccess, fieldAccess)
	fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
}

func (g *genGoMarshal) emitUnmarshalRepeatedString(w io.Writer, f *protoField, fieldAccess, goName, castType string) {
	fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
	fmt.Fprint(w, unmarshalStringLen)
	if castType != "" {
		castPkg, castTyp, castAlias := castTypeComponents(castType)
		if castPkg != "" {
			g.addImport(castPkg, castAlias)
			fmt.Fprintf(w, "\t\t\t%s = append(%s, %s.%s(dAtA[iNdEx:postIndex]))\n", fieldAccess, fieldAccess, castAlias, castTyp)
		} else {
			fmt.Fprintf(w, "\t\t\t%s = append(%s, %s(dAtA[iNdEx:postIndex]))\n", fieldAccess, fieldAccess, castTyp)
		}
	} else {
		fmt.Fprintf(w, "\t\t\t%s = append(%s, string(dAtA[iNdEx:postIndex]))\n", fieldAccess, fieldAccess)
	}
	fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
}

func (g *genGoMarshal) emitUnmarshalMessage(w io.Writer, f *protoField, fieldAccess, goName, pkg, protoName string) {
	fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
	fmt.Fprint(w, unmarshalMsgLen)
	msgType := g.messageGoType(f)
	if f.Nullable {
		fmt.Fprintf(w, "\t\t\tif %s == nil {\n\t\t\t\t%s = &%s{}\n\t\t\t}\n", fieldAccess, fieldAccess, msgType)
	}
	fmt.Fprintf(w, "\t\t\tif err := %s.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", fieldAccess)
	fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
}

func (g *genGoMarshal) emitUnmarshalRepeatedMessage(w io.Writer, f *protoField, fieldAccess, goName string) {
	fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
	fmt.Fprint(w, unmarshalMsgLen)
	msgType := g.messageGoType(f)
	fmt.Fprintf(w, "\t\t\t%s = append(%s, %s{})\n", fieldAccess, fieldAccess, msgType)
	fmt.Fprintf(w, "\t\t\tif err := %s[len(%s)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", fieldAccess, fieldAccess)
	fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
}

func (g *genGoMarshal) emitUnmarshalMapField(w io.Writer, f *protoField, fieldAccess, typeName, castType string) {
	goName := g.goFieldName(f)
	keyProto := f.Type.Key.Name.Name
	valProto := f.Type.Elem.Name.Name
	valPkg := f.Type.Elem.Name.Package

	castKey := unquote(f.Extras["(gogoproto.castkey)"])
	castVal := unquote(f.Extras["(gogoproto.castvalue)"])

	fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field %s\", wireType)\n\t\t\t}\n", goName)
	fmt.Fprint(w, unmarshalMsgLen)
	// Determine the map type: look at the original casttype for the field.
	mapGoType := g.mapGoType(f, castType, castKey, castVal)
	fmt.Fprintf(w, "\t\t\tif %s == nil {\n\t\t\t\t%s = make(%s)\n\t\t\t}\n", fieldAccess, fieldAccess, mapGoType)

	// Determine key and value Go type names.
	mapKeyGoType := g.mapKeyGoType(keyProto, castKey)
	mapValGoType := g.mapValGoType(valProto, valPkg, castVal)

	fmt.Fprint(w, "\t\t\tvar mapkey ", mapKeyGoType, "\n")
	if isProtoMessageType(f.Type.Elem) {
		fmt.Fprintf(w, "\t\t\tmapvalue := &%s{}\n", mapValGoType)
	} else {
		fmt.Fprint(w, "\t\t\tvar mapvalue ", mapValGoType, "\n")
	}
	fmt.Fprint(w, "\t\t\tfor iNdEx < postIndex {\n")
	fmt.Fprint(w, "\t\t\t\tentryPreIndex := iNdEx\n")
	fmt.Fprint(w, "\t\t\t\tvar wire uint64\n")
	fmt.Fprint(w, unmarshalMapEntryWire)
	fmt.Fprint(w, "\t\t\t\tif fieldNum == 1 {\n")
	// Key
	switch {
	case isVarintType(keyProto):
		fmt.Fprintf(w, "\t\t\t\t\tvar mapkeytemp %s\n", protoToGoType(keyProto))
		fmt.Fprint(w, unmarshalMapVarintKeyLoop("mapkeytemp", protoToGoType(keyProto)))
		if castKey != "" {
			castPkg, castTyp, castAlias := castTypeComponents(castKey)
			if castPkg != "" {
				g.addImport(castPkg, castAlias)
				fmt.Fprintf(w, "\t\t\t\t\tmapkey = %s.%s(mapkeytemp)\n", castAlias, castTyp)
			} else {
				fmt.Fprintf(w, "\t\t\t\t\tmapkey = %s(mapkeytemp)\n", castTyp)
			}
		} else {
			fmt.Fprint(w, "\t\t\t\t\tmapkey = mapkeytemp\n")
		}
	default:
		// string key
		fmt.Fprint(w, unmarshalMapStringKey)
		if castKey != "" {
			castPkg, castTyp, castAlias := castTypeComponents(castKey)
			if castPkg != "" {
				g.addImport(castPkg, castAlias)
				fmt.Fprintf(w, "\t\t\t\t\tmapkey = %s.%s(mapkeyBytes)\n", castAlias, castTyp)
			} else {
				fmt.Fprintf(w, "\t\t\t\t\tmapkey = %s(mapkeyBytes)\n", castTyp)
			}
		} else {
			fmt.Fprint(w, "\t\t\t\t\tmapkey = string(mapkeyBytes)\n")
		}
	}
	fmt.Fprint(w, "\t\t\t\t} else if fieldNum == 2 {\n")
	// Value
	switch {
	case isProtoMessageType(f.Type.Elem):
		fmt.Fprint(w, unmarshalMapMsgVal)
	case isVarintType(valProto) && valPkg == "":
		fmt.Fprintf(w, "\t\t\t\t\tvar mapvaltemp %s\n", protoToGoType(valProto))
		fmt.Fprint(w, unmarshalMapVarintKeyLoop("mapvaltemp", protoToGoType(valProto)))
		fmt.Fprint(w, "\t\t\t\t\tmapvalue = mapvaltemp\n")
	case valProto == "bool" && valPkg == "":
		fmt.Fprint(w, "\t\t\t\t\tvar mapvaltemp int\n")
		fmt.Fprint(w, unmarshalMapVarintKeyLoop("mapvaltemp", "int"))
		fmt.Fprint(w, "\t\t\t\t\tmapvalue = bool(mapvaltemp != 0)\n")
	default:
		fmt.Fprint(w, unmarshalMapStringVal)
		if castVal != "" {
			castPkg, castTyp, castAlias := castTypeComponents(castVal)
			if castPkg != "" {
				g.addImport(castPkg, castAlias)
				fmt.Fprintf(w, "\t\t\t\t\tmapvalue = %s.%s(mapvalBytes)\n", castAlias, castTyp)
			} else {
				fmt.Fprintf(w, "\t\t\t\t\tmapvalue = %s(mapvalBytes)\n", castTyp)
			}
		} else {
			fmt.Fprint(w, "\t\t\t\t\tmapvalue = string(mapvalBytes)\n")
		}
	}
	fmt.Fprint(w, "\t\t\t\t} else {\n\t\t\t\t\tiNdEx = entryPreIndex\n")
	fmt.Fprint(w, "\t\t\t\t\tskippy, err := skipGenerated(dAtA[iNdEx:])\n")
	fmt.Fprint(w, "\t\t\t\t\tif err != nil {\n\t\t\t\t\t\treturn err\n\t\t\t\t\t}\n")
	fmt.Fprint(w, "\t\t\t\t\tif (skippy < 0) || (iNdEx+skippy) < 0 {\n\t\t\t\t\t\treturn ErrInvalidLengthGenerated\n\t\t\t\t\t}\n")
	fmt.Fprint(w, "\t\t\t\t\tif (iNdEx + skippy) > postIndex {\n\t\t\t\t\t\treturn io.ErrUnexpectedEOF\n\t\t\t\t\t}\n")
	fmt.Fprint(w, "\t\t\t\t\tiNdEx += skippy\n\t\t\t\t}\n\t\t\t}\n")
	// Assign to map
	if isProtoMessageType(f.Type.Elem) {
		fmt.Fprintf(w, "\t\t\t%s[mapkey] = *mapvalue\n", fieldAccess)
	} else {
		fmt.Fprintf(w, "\t\t\t%s[mapkey] = mapvalue\n", fieldAccess)
	}
	fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
}

// ---------------------------------------------------------------------------
// Unmarshal for optional alias
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitUnmarshalOptionalAlias(w io.Writer, typeName string, f *protoField) {
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package

	fmt.Fprintf(w, "func (m *%s) Unmarshal(dAtA []byte) error {\n", typeName)
	fmt.Fprint(w, unmarshalPreamble)
	fmt.Fprintf(w, "\t\tif wireType == 4 {\n\t\t\treturn fmt.Errorf(\"proto: %s: wiretype end group for non-group\")\n\t\t}\n", typeName)
	fmt.Fprintf(w, "\t\tif fieldNum <= 0 {\n\t\t\treturn fmt.Errorf(\"proto: %s: illegal tag %%d (wire type %%d)\", fieldNum, wireType)\n\t\t}\n", typeName)
	fmt.Fprint(w, "\t\tswitch fieldNum {\n\t\tcase 1:\n")

	switch {
	case isMessageType(f):
		fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field Items\", wireType)\n\t\t\t}\n")
		fmt.Fprint(w, unmarshalMsgLen)
		msgType := g.messageGoType(f)
		fmt.Fprintf(w, "\t\t\t*m = append(*m, %s{})\n", msgType)
		fmt.Fprintf(w, "\t\t\tif err := (*m)[len(*m)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n")
		fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
	case isVarintType(protoName) && pkg == "":
		fmt.Fprintf(w, "\t\t\tif wireType != 0 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field Items\", wireType)\n\t\t\t}\n")
		goTypeName := protoToGoType(protoName)
		fmt.Fprintf(w, "\t\t\tvar v %s\n", goTypeName)
		fmt.Fprint(w, unmarshalVarintLoop("v", goTypeName))
		fmt.Fprintf(w, "\t\t\t*m = append(*m, v)\n")
	default:
		// string
		fmt.Fprintf(w, "\t\t\tif wireType != 2 {\n\t\t\t\treturn fmt.Errorf(\"proto: wrong wireType = %%d for field Items\", wireType)\n\t\t\t}\n")
		fmt.Fprint(w, unmarshalStringLen)
		fmt.Fprintf(w, "\t\t\t*m = append(*m, string(dAtA[iNdEx:postIndex]))\n")
		fmt.Fprint(w, "\t\t\tiNdEx = postIndex\n")
	}
	fmt.Fprint(w, "\t\tdefault:\n")
	fmt.Fprint(w, "\t\t\tiNdEx = preIndex\n")
	fmt.Fprint(w, "\t\t\tskippy, err := skipGenerated(dAtA[iNdEx:])\n")
	fmt.Fprint(w, "\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n")
	fmt.Fprint(w, "\t\t\tif (skippy < 0) || (iNdEx+skippy) < 0 {\n\t\t\t\treturn ErrInvalidLengthGenerated\n\t\t\t}\n")
	fmt.Fprint(w, "\t\t\tif (iNdEx + skippy) > l {\n\t\t\t\treturn io.ErrUnexpectedEOF\n\t\t\t}\n")
	fmt.Fprint(w, "\t\t\tiNdEx += skippy\n")
	fmt.Fprint(w, "\t\t}\n\t}\n\n\tif iNdEx > l {\n\t\treturn io.ErrUnexpectedEOF\n\t}\n\treturn nil\n}\n\n")
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

func (g *genGoMarshal) emitString(w io.Writer, typeName string, fields []protoField, optional bool) {
	fmt.Fprintf(w, "func (this *%s) String() string {\n", typeName)
	fmt.Fprint(w, "\tif this == nil {\n\t\treturn \"nil\"\n\t}\n")
	if optional || len(fields) == 0 {
		io.WriteString(w, "\treturn fmt.Sprintf(\"%v\", *this)\n}\n\n") //nolint:errcheck
		return
	}
	// Build gogo-compatible String for struct types.
	fmt.Fprint(w, "\ts := strings.Join([]string{`&", typeName, "{`,\n")
	for _, f := range fields {
		goName := g.goFieldName(&f)
		fieldExpr := "this." + goName
		protoName := f.Type.Name.Name
		if f.Map {
			fmt.Fprintf(w, "\t\t`%s:` + fmt.Sprintf(\"%%v\", %s) + `,`,\n", goName, fieldExpr)
			continue
		}
		if f.Repeated {
			fmt.Fprintf(w, "\t\t`%s:` + fmt.Sprintf(\"%%v\", %s) + `,`,\n", goName, fieldExpr)
			continue
		}
		if isMessageType(&f) {
			if f.Nullable {
				fmt.Fprintf(w, "\t\t`%s:` + strings.Replace(fmt.Sprintf(\"%%v\", %s), \"%s\", \"%s\", 1) + `,`,\n", goName, fieldExpr, protoName, protoName)
			} else {
				fmt.Fprintf(w, "\t\t`%s:` + strings.Replace(strings.Replace(%s.String(), \"%s\", \"%s\", 1), `&`, ``, 1) + `,`,\n", goName, fieldExpr, f.Type.Name.Name, f.Type.Name.Name)
			}
			continue
		}
		if f.Nullable {
			// Pointer to a primitive type: use valueToStringGenerated for nil-safe formatting.
			g.addImport("reflect", "reflect")
			fmt.Fprintf(w, "\t\t`%s:` + valueToStringGenerated(%s) + `,`,\n", goName, fieldExpr)
			continue
		}
		fmt.Fprintf(w, "\t\t`%s:` + fmt.Sprintf(\"%%v\", %s) + `,`,\n", goName, fieldExpr)
	}
	fmt.Fprint(w, "\t\t`}`,\n\t}, \"\")\n\treturn s\n}\n\n")
}

// ---------------------------------------------------------------------------
// Helper: Get the Go field name from a protoField
// ---------------------------------------------------------------------------

func (g *genGoMarshal) goFieldName(f *protoField) string {
	if cn, ok := f.Extras["(gogoproto.customname)"]; ok {
		return unquote(cn)
	}
	// No customname set means the proto name IS the Go field name (rare).
	return f.Name
}

// ---------------------------------------------------------------------------
// Helper: Get the Go type expression for a message field
// ---------------------------------------------------------------------------

func (g *genGoMarshal) messageGoType(f *protoField) string {
	if f.Type == nil {
		return "interface{}"
	}
	typName := f.Type.Name.Name
	typPath := f.Type.Name.Path // e.g. "k8s.io/apimachinery/pkg/runtime/generated.proto"

	// Strip the /generated.proto suffix to get the Go package path.
	goPkg := strings.TrimSuffix(typPath, "/generated.proto")
	if goPkg == "" || goPkg == g.localGoPackage.Package {
		return typName
	}
	// It's a cross-package type.
	alias := goImportAlias(goPkg)
	g.addImport(goPkg, alias)
	return alias + "." + typName
}

// ---------------------------------------------------------------------------
// Helpers: map type resolution
// ---------------------------------------------------------------------------

func (g *genGoMarshal) mapGoType(f *protoField, castType, castKey, castVal string) string {
	if castType != "" {
		castPkg, castTyp, castAlias := castTypeComponents(castType)
		if castPkg != "" {
			g.addImport(castPkg, castAlias)
			return castAlias + "." + castTyp
		}
		return castTyp
	}
	// Construct map[K]V from key and value proto types.
	keyType := g.mapKeyGoType(f.Type.Key.Name.Name, castKey)
	valType := g.mapValGoType(f.Type.Elem.Name.Name, f.Type.Elem.Name.Package, castVal)
	return "map[" + keyType + "]" + valType
}

func (g *genGoMarshal) mapKeyGoType(keyProto, castKey string) string {
	if castKey != "" {
		castPkg, castTyp, castAlias := castTypeComponents(castKey)
		if castPkg != "" {
			g.addImport(castPkg, castAlias)
			return castAlias + "." + castTyp
		}
		return castTyp
	}
	return protoToGoType(keyProto)
}

func (g *genGoMarshal) mapValGoType(valProto, valPkg, castVal string) string {
	if castVal != "" {
		castPkg, castTyp, castAlias := castTypeComponents(castVal)
		if castPkg != "" {
			g.addImport(castPkg, castAlias)
			return castAlias + "." + castTyp
		}
		return castTyp
	}
	if valPkg != "" {
		// message type from another proto package
		return valProto // caller will have set the import
	}
	return protoToGoType(valProto)
}

// ---------------------------------------------------------------------------
// Import tracking
// ---------------------------------------------------------------------------

func (g *genGoMarshal) addImport(pkgPath, alias string) {
	if _, ok := g.neededImports[pkgPath]; !ok {
		g.neededImports[pkgPath] = alias
	}
}

// ---------------------------------------------------------------------------
// Package-level helpers (non-method)
// ---------------------------------------------------------------------------

// receiverExpr returns the method receiver string.
// optional=true and valueRecv=true means value receiver (for optional alias marshal methods).
// Struct methods always use pointer receiver.
func receiverExpr(typeName string, optional, valueRecv bool) string {
	if optional && valueRecv {
		return "m " + typeName
	}
	return "m *" + typeName
}

// writeTagCode generates Go statements to write the proto tag bytes
// in reverse order (right-to-left) for MarshalToSizedBuffer.
func writeTagCode(fieldNum, wireType int) string {
	// Encode the tag as a protobuf varint and emit statements that write
	// the bytes in reverse order (right-to-left), as required by MarshalToSizedBuffer.
	tag := fieldNum<<3 | wireType
	var fwd []byte
	t := tag
	for t > 0 {
		b := byte(t & 0x7f)
		t >>= 7
		if t > 0 {
			b |= 0x80
		}
		fwd = append(fwd, b)
	}
	var lines []string
	for j := len(fwd) - 1; j >= 0; j-- {
		lines = append(lines, fmt.Sprintf("\ti--\n\tdAtA[i] = 0x%02x", fwd[j]))
	}
	return strings.Join(lines, "\n")
}

// tagSizeExpr returns the string representation of the number of bytes used
// by the proto tag for the given field number and wire type.
func tagSizeExpr(fieldNum, wireType int) string {
	tag := fieldNum<<3 | wireType
	n := 1
	for t := tag; t >= 0x80; t >>= 7 {
		n++
	}
	return fmt.Sprintf("%d", n)
}

// indentTagCode re-indents the tag code block (used inside nested loops).
func indentTagCode(code, indent string) string {
	lines := strings.Split(code, "\n")
	var out []string
	for _, l := range lines {
		if l == "" {
			out = append(out, l)
			continue
		}
		// Remove the leading single \t and add the desired indent.
		out = append(out, indent+strings.TrimPrefix(l, "\t"))
	}
	return strings.Join(out, "\n")
}

// protoWireType returns the protobuf wire type for a field.
func protoWireType(f *protoField) int {
	if f.Map {
		return 2 // length-delimited (map entries are messages)
	}
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package
	switch {
	case isMessageType(f):
		return 2
	case protoName == "bool" && pkg == "":
		return 0
	case isVarintType(protoName) && pkg == "":
		return 0
	case protoName == "double" && pkg == "":
		return 1
	case protoName == "float" && pkg == "":
		return 5
	default:
		return 2 // string, bytes
	}
}

// isMessageType returns true if the proto field represents an embedded message.
func isMessageType(f *protoField) bool {
	if f.Type == nil {
		return false
	}
	protoName := f.Type.Name.Name
	pkg := f.Type.Name.Package
	// Primitive protobuf types have empty package.
	switch protoName {
	case "string", "bytes", "bool", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "double", "float":
		if pkg == "" {
			return false
		}
	}
	// If there's a non-empty package or the type name doesn't match a primitive,
	// it's a message type.  Also check that it has Key/Elem for map (handled separately).
	if f.Map {
		return false
	}
	return pkg != "" || !isPrimitiveProtoType(protoName)
}

// isProtoMessageType returns true if a standalone *types.Type is a protobuf message.
func isProtoMessageType(t *types.Type) bool {
	if t == nil {
		return false
	}
	protoName := t.Name.Name
	pkg := t.Name.Package
	return pkg != "" || !isPrimitiveProtoType(protoName)
}

func isPrimitiveProtoType(name string) bool {
	switch name {
	case "string", "bytes", "bool",
		"int32", "int64", "uint32", "uint64",
		"sint32", "sint64",
		"fixed32", "fixed64", "sfixed32", "sfixed64",
		"double", "float":
		return true
	}
	return false
}

func isVarintType(protoName string) bool {
	switch protoName {
	case "int32", "int64", "uint32", "uint64", "sint32", "sint64":
		return true
	}
	return false
}

// protoToGoType maps a protobuf primitive type name to its Go equivalent.
func protoToGoType(protoName string) string {
	switch protoName {
	case "string":
		return "string"
	case "bytes":
		return "[]byte"
	case "bool":
		return "bool"
	case "int32", "sint32", "sfixed32":
		return "int32"
	case "int64", "sint64", "sfixed64":
		return "int64"
	case "uint32", "fixed32":
		return "uint32"
	case "uint64", "fixed64":
		return "uint64"
	case "double":
		return "float64"
	case "float":
		return "float32"
	}
	return protoName
}

// goImportAlias converts a Go package path to a valid import alias by
// replacing all `/`, `.`, and `-` with `_`.
func goImportAlias(pkgPath string) string {
	r := strings.NewReplacer("/", "_", ".", "_", "-", "_")
	return r.Replace(pkgPath)
}

// castTypeComponents splits a casttype string like "k8s.io/apimachinery/pkg/types.UID"
// into (pkgPath, typeName, importAlias).
// For local types like "StatusReason", returns ("", "StatusReason", "").
func castTypeComponents(castType string) (pkgPath, typeName, importAlias string) {
	lastDot := strings.LastIndex(castType, ".")
	if lastDot < 0 || !strings.Contains(castType[:lastDot], "/") {
		// Local type (no package path slash) or simple identifier.
		return "", castType, ""
	}
	pkgPath = castType[:lastDot]
	typeName = castType[lastDot+1:]
	importAlias = goImportAlias(pkgPath)
	return
}

// unquote removes surrounding double-quotes from a string (as stored in Extras).
func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// ---------------------------------------------------------------------------
// Reusable unmarshal code fragments
// ---------------------------------------------------------------------------

const unmarshalPreamble = `	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
`

const unmarshalStringLen = `			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
`

const unmarshalBytesLen = `			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
`

const unmarshalMsgLen = `			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
`

// unmarshalVarintLoop generates the varint read loop for a variable name and type.
func unmarshalVarintLoop(varName, typeName string) string {
	return fmt.Sprintf(`			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				%s |= %s(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
`, varName, typeName)
}

const unmarshalMapEntryWire = `			for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenerated
					}
					if iNdEx >= postIndex {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
`

// unmarshalMapVarintKeyLoop is the varint read loop used inside map entry parsing.
func unmarshalMapVarintKeyLoop(varName, typeName string) string {
	return fmt.Sprintf(`				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenerated
					}
					if iNdEx >= postIndex {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					%s |= %s(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
`, varName, typeName)
}

const unmarshalMapStringKey = `				var stringLenKey uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenerated
					}
					if iNdEx >= postIndex {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLenKey |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLenKey := int(stringLenKey)
				if intStringLenKey < 0 {
					return ErrInvalidLengthGenerated
				}
				postStringIndexKey := iNdEx + intStringLenKey
				if postStringIndexKey < 0 {
					return ErrInvalidLengthGenerated
				}
				if postStringIndexKey > l {
					return io.ErrUnexpectedEOF
				}
				mapkeyBytes := dAtA[iNdEx:postStringIndexKey]
				iNdEx = postStringIndexKey
`

const unmarshalMapStringVal = `				var stringLenVal uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenerated
					}
					if iNdEx >= postIndex {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLenVal |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLenVal := int(stringLenVal)
				if intStringLenVal < 0 {
					return ErrInvalidLengthGenerated
				}
				postStringIndexVal := iNdEx + intStringLenVal
				if postStringIndexVal < 0 {
					return ErrInvalidLengthGenerated
				}
				if postStringIndexVal > l {
					return io.ErrUnexpectedEOF
				}
				mapvalBytes := dAtA[iNdEx:postStringIndexVal]
				iNdEx = postStringIndexVal
`

const unmarshalMapMsgVal = `				var msglenVal int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenerated
					}
					if iNdEx >= postIndex {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglenVal |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglenVal < 0 {
					return ErrInvalidLengthGenerated
				}
				postMsgIndex := iNdEx + msglenVal
				if postMsgIndex < 0 {
					return ErrInvalidLengthGenerated
				}
				if postMsgIndex > l {
					return io.ErrUnexpectedEOF
				}
				if err := mapvalue.Unmarshal(dAtA[iNdEx:postMsgIndex]); err != nil {
					return err
				}
				iNdEx = postMsgIndex
`

// ---------------------------------------------------------------------------
// Shared helper code emitted once per file
// ---------------------------------------------------------------------------

const sharedHelpersCode = `func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenerated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func sovGenerated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenerated
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenerated
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenerated        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenerated = fmt.Errorf("proto: unexpected end of group")
)
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
`
