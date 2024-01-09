package httpcli

import (
	"fmt"
	"mime"
)

// ManifestMediaTypes returns the supported media types for manifests.
func ManifestMediaTypes() (mediaTypes []string) {
	for t := range mappings {
		if t != "" {
			mediaTypes = append(mediaTypes, t)
		}
	}
	return
}

// UnmarshalFunc implements manifest unmarshalling a given MediaType
type UnmarshalFunc func([]byte) (interface{}, error)

var mappings = make(map[string]UnmarshalFunc)

// RegisterManifestSchema registers an UnmarshalFunc for a given schema type.  This
// should be called from specific
func RegisterManifestSchema(mediaType string, u UnmarshalFunc) error {
	if _, ok := mappings[mediaType]; ok {
		return fmt.Errorf("manifest media type registration would overwrite existing: %s", mediaType)
	}
	mappings[mediaType] = u
	return nil
}

// UnmarshalManifest looks up manifest unmarshal functions based on
// MediaType
func UnmarshalManifest(ctHeader string, p []byte) (interface{}, error) {
	// Need to look up by the actual media type, not the raw contents of
	// the header. Strip semicolons and anything following them.
	var mediaType string
	if ctHeader != "" {
		var err error
		mediaType, _, err = mime.ParseMediaType(ctHeader)
		if err != nil {
			return nil, err
		}
	}

	unmarshalFunc, ok := mappings[mediaType]
	if !ok {
		unmarshalFunc, ok = mappings[""]
		if !ok {
			return nil, fmt.Errorf("unsupported manifest media type and no default available: %s", mediaType)
		}
	}

	return unmarshalFunc(p)
}
