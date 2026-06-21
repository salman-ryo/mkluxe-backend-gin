package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims defines the structure of our JWT payload.
// It combines:
// - custom fields (UserID, Role)
// - standard JWT fields (via jwt.RegisteredClaims)
type JWTClaims struct {
	UserID string `json:"user_id"` // custom claim: identifies the user in our system
	Role   string `json:"role"`    // custom claim: defines user permissions (admin, user, etc.)

	// Embedded struct from jwt library.
	// This adds standard JWT fields like:
	// - exp (expiration time)
	// - iat (issued at time)
	// - nbf (not before)
	jwt.RegisteredClaims
}

// getJWTSecret reads the secret key used to sign JWT tokens.
func getJWTSecret() []byte {
	// Read secret from environment variable (recommended for production)
	secret := os.Getenv("JWT_SECRET")

	// Fallback secret for local development only (NEVER use in production)
	if secret == "" {
		secret = "local_dev_secret_key_only"
	}

	// JWT library expects secret as raw bytes, not string
	return []byte(secret)
}

// GenerateTokens creates both access and refresh tokens for a user.
// Access token = short-lived (15 min)
// Refresh token = long-lived (7 days)
func GenerateTokens(userID, role string) (string, string, error) {

	// Get signing key used for HMAC SHA256 signing
	secret := getJWTSecret()

	// Current timestamp used for issuing tokens
	now := time.Now()

	// ----------------------------
	// ACCESS TOKEN CREATION
	// ----------------------------

	accessClaims := JWTClaims{
		UserID: userID, // embed user identity
		Role:   role,   // embed user role

		RegisteredClaims: jwt.RegisteredClaims{
			// exp = token expiration time (15 minutes from now)
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),

			// iat = when token was issued
			IssuedAt: jwt.NewNumericDate(now),
		},
	}

	// Create JWT object and sign it using HS256 algorithm + secret key
	accessToken, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString(secret)

	// If signing fails, return error immediately
	if err != nil {
		return "", "", err
	}

	// ----------------------------
	// REFRESH TOKEN CREATION
	// ----------------------------

	refreshClaims := JWTClaims{
		UserID: userID,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			// Refresh token lives longer (7 days)
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),

			IssuedAt: jwt.NewNumericDate(now),
		},
	}

	// Sign refresh token using same algorithm and secret
	refreshToken, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString(secret)

	// Handle signing failure
	if err != nil {
		return "", "", err
	}

	// Return both tokens to caller (usually login handler)
	return accessToken, refreshToken, nil
}

// ValidateToken verifies a JWT string and extracts claims if valid
func ValidateToken(tokenString string) (*JWTClaims, error) {

	// Parse token and attach expected claim structure (JWTClaims)
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{}, // empty struct where decoded claims will be stored
		func(token *jwt.Token) (interface{}, error) {

			// Ensure token was signed using HMAC (security check)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			// Provide secret key for signature verification
			return getJWTSecret(), nil
		},
	)

	// If parsing fails (invalid format, bad signature, etc.)
	if err != nil {
		return nil, err
	}

	// Type assertion: extract claims from parsed token
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {

		// Token is valid → return extracted user data
		return claims, nil
	}

	// If token is invalid or claims mismatch
	return nil, errors.New("invalid token")
}
