package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stevanusy21/golang_sandbox/pkg/utils"
)

const LogLocation = "Token"

type Claims struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		utils.LogErrorNoValue(LogLocation, utils.ErrJwtSecretMissing.Error())
		return "", utils.ErrJwtSecretMissing
	}

	claims := Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		utils.LogError(LogLocation, utils.ErrTokenCreatedFailed.Error(), err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		utils.LogErrorNoValue(LogLocation, utils.ErrJwtSecretMissing.Error())
		return nil, utils.ErrJwtSecretMissing
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		utils.LogError(LogLocation, utils.ErrTokenParseFailed.Error(), err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, utils.ErrTokenInvalid
	}

	return claims, nil
}
