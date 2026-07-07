package goroute

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type JWTClaims struct {
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

// GenerateJWT creates a secure Header.Payload.Signature string
func GenerateJWT(username string, secret string, duration time.Duration) (string, error) {
	// 1. Create header and payload strings
	headerJSON, _ := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	payloadJSON, _ := json.Marshal(JWTClaims{
		Username: username,
		Exp:      time.Now().Add(duration).Unix(),
	})
	// 2. Encode to Base64URL
	encodeHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodePayload := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// 3. generate cryptographic signature
	signingInput := encodeHeader + "." + encodePayload
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature, nil
}

// VerifyJWT parse and validates the token signature and expiration
func VerifyJWT(tokenStr string, secret string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("malformed token signature")
	}

	// recompute signature to verify data integrity
	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if parts[2] != expectedSignature {
		return nil, errors.New("cryptographic signature mismatch")
	}

	// decode the payload json
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	// verify expiration bounds
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token has expired")
	}
	return &claims, nil
}

// JWTAuth returns a protection gateway middleware
func JWTAuth(secret string) RouteHandler {
	return func(c *Context) {
		authHeader := c.Req.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Header must follow the standard convention: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization format must be 'Bearer <token>'"})
			c.Abort()
			return
		}

		claims, err := VerifyJWT(parts[1], secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			c.Abort()
			return
		}
		// Inject identity parameters cleanly into context state storage
		c.Set("username", claims.Username)
		c.Next()
	}
}
