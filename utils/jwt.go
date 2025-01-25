package utils

import (
	"ticketink/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	*jwt.RegisteredClaims
}

func GenerateToken(user models.User) (string, error) {

	claims := Claims{
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {

		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, jwt.ErrTokenExpired
		}

		return nil, jwt.NewValidationError("Invalid token", jwt.ValidationErrorClaimsInvalid)
	}

	if !token.Valid {
		return nil, jwt.NewValidationError("Invalid token", jwt.ValidationErrorClaimsInvalid)
	}

	return claims, nil
}
