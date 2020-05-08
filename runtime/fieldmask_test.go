package runtime

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func newFieldMask(paths ...string) *field_mask.FieldMask {
	return &field_mask.FieldMask{Paths: paths}
}

func TestFieldMaskFromRequestBody(t *testing.T) {
	for _, tc := range []struct {
		name        string
		input       string
		msg         protoreflect.Message
		expected    *field_mask.FieldMask
		expectedErr error
	}{
		{
			name:     "empty",
			expected: newFieldMask(),
		},
		{
			name:     "simple",
			input:    `{"uuid":"1234", "floatValue":3.14}`,
			msg:      (&examplepb.ABitOfEverything{}).ProtoReflect(),
			expected: newFieldMask("uuid", "float_value"),
		},
		{
			name:     "nested",
			input:    `{"single_nested": {"name":"bob", "amount": 2}, "uuid":"1234"}`,
			msg:      (&examplepb.ABitOfEverything{}).ProtoReflect(),
			expected: newFieldMask("single_nested.name", "single_nested.amount", "uuid"),
		},
		{
			name:     "map",
			input:    `{"mapped_string_value": {"a": "x"}}`,
			msg:      (&examplepb.ABitOfEverything{}).ProtoReflect(),
			expected: newFieldMask("mapped_string_value"),
		},
		{
			name: "complex",
			input: `
			{
				"single_nested": {
					"name": "bar",
					"amount": 10,
					"ok": "TRUE"
				},
				"uuid": "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
				"nested": [
					{
						"name": "bar",
						"amount": 10
					},
					{
						"name": "baz",
						"amount": 20
					}
				],
				"float_value": 1.5,
				"double_value": 2.5,
				"int64_value": 4294967296,
				"int64_override_type": 12345,
				"int32_value": -2147483648,
				"uint64_value": 9223372036854775807,
				"uint32_value": 4294967295,
				"fixed64_value": 9223372036854775807,
				"fixed32_value": 4294967295,
				"sfixed64_value": -4611686018427387904,
				"sfixed32_value": 2147483647,
				"sint64_value": 4611686018427387903,
				"sint32_value": 2147483647,
				"bool_value": true,
				"string_value": "strprefix/foo",
				"bytes_value": "132456",
				"enum_value": "ONE",
				"oneof_string": "x",
				"nonConventionalNameValue": "camelCase",
				"timestamp_value": "2016-05-10T10:19:13.123Z",
				"enum_value_annotation": "ONE",
				"nested_annotation": {
					"name": "hoge",
					"amount": 10
				}
			}
`,
			msg: (&examplepb.ABitOfEverything{}).ProtoReflect(),
			expected: newFieldMask(
				"single_nested.name",
				"single_nested.amount",
				"single_nested.ok",
				"uuid",
				"float_value",
				"double_value",
				"int64_value",
				"int64_override_type",
				"int32_value",
				"uint64_value",
				"uint32_value",
				"fixed64_value",
				"fixed32_value",
				"sfixed64_value",
				"sfixed32_value",
				"sint64_value",
				"sint32_value",
				"bool_value",
				"string_value",
				"bytes_value",
				"enum_value",
				"oneof_string",
				"nonConventionalNameValue",
				"timestamp_value",
				"enum_value_annotation",
				"nested_annotation.name",
				"nested_annotation.amount",
				"nested",
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := FieldMaskFromRequestBody(bytes.NewReader([]byte(tc.input)), tc.msg)
			if diff := cmp.Diff(tc.expected, actual, cmpopts.SortSlices(func(x, y string) bool {
				return x < y
			})); diff != "" {
				t.Errorf("field masks differed:\n%s", diff)
			}
			if err != tc.expectedErr {
				t.Errorf("want %v; got %v", tc.expectedErr, err)
			}
		})
	}
}

// avoid compiler optimising benchmark away
var result *field_mask.FieldMask

func BenchmarkABEFieldMaskFromRequestBody(b *testing.B) {
	input := `{` +
		`"single_nested":				{"name": "bar",` +
		`                           	 "amount": 10,` +
		`                           	 "ok": "TRUE"},` +
		`"uuid":						"6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",` +
		`"nested": 						[{"name": "bar",` +
		`								  "amount": 10},` +
		`								 {"name": "baz",` +
		`								  "amount": 20}],` +
		`"float_value":             	1.5,` +
		`"double_value":            	2.5,` +
		`"int64_value":             	4294967296,` +
		`"uint64_value":            	9223372036854775807,` +
		`"int32_value":             	-2147483648,` +
		`"fixed64_value":           	9223372036854775807,` +
		`"fixed32_value":           	4294967295,` +
		`"bool_value":              	true,` +
		`"string_value":            	"strprefix/foo",` +
		`"bytes_value":					"132456",` +
		`"uint32_value":            	4294967295,` +
		`"enum_value":     		        "ONE",` +
		`"path_enum_value":	    	    "DEF",` +
		`"nested_path_enum_value":  	"JKL",` +
		`"sfixed32_value":          	2147483647,` +
		`"sfixed64_value":          	-4611686018427387904,` +
		`"sint32_value":            	2147483647,` +
		`"sint64_value":            	4611686018427387903,` +
		`"repeated_string_value": 		["a", "b", "c"],` +
		`"oneof_value":					{"oneof_string":"x"},` +
		`"map_value": 					{"a": "ONE",` +
		`								 "b": "ZERO"},` +
		`"mapped_string_value": 		{"a": "x",` +
		`								 "b": "y"},` +
		`"mapped_nested_value": 		{"a": {"name": "x", "amount": 1},` +
		`								 "b": {"name": "y", "amount": 2}},` +
		`"nonConventionalNameValue":	"camelCase",` +
		`"timestamp_value":				"2016-05-10T10:19:13.123Z",` +
		`"repeated_enum_value":			["ONE", "ZERO"],` +
		`"repeated_enum_annotation":	 ["ONE", "ZERO"],` +
		`"enum_value_annotation": 		"ONE",` +
		`"repeated_string_annotation":	["a", "b"],` +
		`"repeated_nested_annotation": 	[{"name": "hoge",` +
		`								  "amount": 10},` +
		`								 {"name": "fuga",` +
		`								  "amount": 20}],` +
		`"nested_annotation": 			{"name": "hoge",` +
		`								 "amount": 10},` +
		`"int64_override_type":			12345` +
		`}`
	var r *field_mask.FieldMask
	var err error
	for i := 0; i < b.N; i++ {
		r, err = FieldMaskFromRequestBody(bytes.NewReader([]byte(input)), nil)
	}
	if err != nil {
		b.Error(err)
	}
	result = r
}

func BenchmarkNonStandardFieldMaskFromRequestBody(b *testing.B) {
	input := `{` +
		`"id":			"foo",` +
		`"Num": 		2,` +
		`"line_num": 	3,` +
		`"langIdent":	"bar",` +
		`"STATUS": 		"baz"` +
		`}`
	var r *field_mask.FieldMask
	var err error
	for i := 0; i < b.N; i++ {
		r, err = FieldMaskFromRequestBody(bytes.NewReader([]byte(input)), nil)
	}
	if err != nil {
		b.Error(err)
	}
	result = r
}
