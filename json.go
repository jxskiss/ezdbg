package ezdbg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

func JSON(v interface{}) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	b := bytes.TrimSpace(buf.Bytes())
	return bytesToString(b)
}

// Pretty converts given object to a pretty formatted json string.
// If the input is a json string, it will be formatted using json.Indent
// with four space characters as indent.
func Pretty(v interface{}) string {
	return prettyIndent(v, "    ")
}

// Pretty2 is like Pretty, but it uses two space characters as indent,
// instead of four.
func Pretty2(v interface{}) string {
	return prettyIndent(v, "  ")
}

func prettyIndent(v interface{}, indent string) string {
	var src []byte
	switch v := v.(type) {
	case []byte:
		src = v
	case string:
		src = stringToBytes(v)
	}
	if src != nil {
		if json.Valid(src) {
			buf := bytes.NewBuffer(nil)
			_ = json.Indent(buf, src, "", indent)
			return bytesToString(buf.Bytes())
		}
		if utf8.Valid(src) {
			return string(src)
		}
		return "<pretty: non-printable bytes>"
	}
	buf, err := marshalJSONDisableHTMLEscape(v, "", indent)
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	buf = bytes.TrimSpace(buf)
	return bytesToString(buf)
}

func marshalJSONDisableHTMLEscape(v interface{}, prefix, indent string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetIndent(prefix, indent)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
