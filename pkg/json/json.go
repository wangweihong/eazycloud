package json

import "encoding/json"

type RawMessage = json.RawMessage

// 根据性能要求切换到其他的编解码器.
var (
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)
