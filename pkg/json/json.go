package json

import "encoding/json"

type RawMessage = json.RawMessage

// 根据性能要求切换到其他的编解码器.
var (
	Marshal = json.Marshal
	// 转换JSON的过程中，Unmarshal会按字段名来转换，碰到同名但大小写不一样的字段会当做同一字段处理.
	// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys used by Marshal (either the
	// struct field name or its tag), preferring an exact match but also accepting a case-insensitive match. By default,
	// object keys which don't have a corresponding struct field are ignored (see Decoder.DisallowUnknownFields for an
	// alternative).
	// https://pkg.go.dev/encoding/json#Unmarshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)
