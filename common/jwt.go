package common

import (
	"budget-app/internal/model"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"os"
	"time"
)

type CustomJwtClaims struct {
	Id uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwt(user model.User) (*string, *string, error) {
	userClaims := CustomJwtClaims{
		Id: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomJwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJwtSignedAccessToken(signedAccessToken string) (*CustomJwtClaims, error) {
	parsedJwtAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &CustomJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Error(err)
		return nil, err
	} else if claims, ok := parsedJwtAccessToken.Claims.(*CustomJwtClaims); ok {
		return claims, nil
	} else {
		log.Error("unknown claims type, cannot proceed")
		return nil, errors.New("unknown claims type")
	}
}

func IsExpired(claims *CustomJwtClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	fmt.Println("current time is", currentTime)
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}
