package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claim := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %w", err)
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	uuidStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	retVal, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invlaid UUID: %w", err)
	}

	return retVal, nil

}
