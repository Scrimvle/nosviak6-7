package packages

import "encoding/base64"

// ENCODING is the standard encoding package
var ENCODING map[string]any = map[string]any{
	"base64": map[string]any{
		"encode": func(c string) string {
			return base64.StdEncoding.EncodeToString([]byte(c))
		},
	},
}