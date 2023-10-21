package utils

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

// CreateJSONDecoder - creates a new json decoder with custom settings
func CreateJSONDecoder(data io.Reader) *jsoniter.Decoder {
	decoder := jsoniter.NewDecoder(data)
	decoder.DisallowUnknownFields()
	return decoder
}
