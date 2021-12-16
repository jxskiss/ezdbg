package ezdbg

import (
	"bytes"
	"encoding/json"
	"fmt"
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
