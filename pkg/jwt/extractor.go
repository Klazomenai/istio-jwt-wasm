package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// Claims represents the JWT claims we care about
type Claims struct {
	JTI string `json:"jti"` // JWT ID
}

// ExtractJTI extracts the JTI (JWT ID) claim from a JWT token without signature verification.
// We skip verification here because Istio's RequestAuthentication already validated the token.
func ExtractJTI(token string) (string, error) {
	// Split JWT into parts: header.payload.signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid JWT format: expected 3 parts")
	}

	// Decode the payload (part 1, zero-indexed)
	payload, err := base64DecodeSegment(parts[1])
	if err != nil {
		return "", errors.New("failed to decode JWT payload: " + err.Error())
	}

	// Parse JSON payload
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", errors.New("failed to parse JWT claims: " + err.Error())
	}

	// Validate JTI is present
	if claims.JTI == "" {
		return "", errors.New("JWT missing 'jti' claim")
	}

	return claims.JTI, nil
}

// base64DecodeSegment decodes a base64url-encoded JWT segment
func base64DecodeSegment(seg string) ([]byte, error) {
	// JWT uses base64url encoding (RFC 4648), which Go's base64.RawURLEncoding handles
	return base64.RawURLEncoding.DecodeString(seg)
}
