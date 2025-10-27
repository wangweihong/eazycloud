package pkce

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

func VerifyCodeChallenge(codeVerifier, codeChallenge, method string) bool {
	switch method {
	case "plain":
		return codeVerifier == codeChallenge
	case "S256":
		hash := sha256.Sum256([]byte(codeVerifier))
		expected := base64.RawURLEncoding.EncodeToString(hash[:])
		return expected == codeChallenge
	default:
		return false
	}
}

func ValidateCodeChallengeMethod(method string) error {
	if method != "plain" && method != "S256" {
		return errors.Errorf("unsupported code challenge method")
	}
	return nil
}

func NormalizeBase64(s string) string {
	return strings.TrimRight(s, "=")
}
