package lib

import "encoding/json"

// PrettyPrint pretty prints maps and structs.
func PrettyPrint(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
