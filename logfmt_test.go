package ezdbg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogfmt(t *testing.T) {
	tests := []map[string]interface{}{
		{
			"value": 123,
			"want":  "123",
		},
		{
			"value": (*string)(nil),
			"want":  "null",
		},
		{
			"value": comptyp{
				I32:      32,
				I32_p:    i32ptr(32),
				I64:      64,
				I64_p:    nil,
				Str:      "str",
				Str_p:    strptr("str with space"),
				Simple:   simple{A: "simple.A"},
				Simple_p: nil,
			},
			"want": `i32=32 i32_p=32 i64=64 str=str str_p="str with space"`,
		},
		{
			"value": map[string]interface{}{
				"a": 1234,
				"b": "bcde",
				"c": 123.456,
				"d": simple{A: "simple.A"},
				"e": nil,
				"f": []byte("I'm bytes"),
			},
			"want": `a=1234 b=bcde c=123.456 f="I'm bytes"`,
		},
	}
	for _, test := range tests {
		got := Logfmt(test["value"])
		assert.Equal(t, test["want"], got)
	}
}

func i32ptr(x int32) *int32   { return &x }
func strptr(x string) *string { return &x }
