package utils

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SCRECT_KEY = os.Getenv("SCRECT_KEY")

type SignedDetails struct {
	Email    string
	Password string
	uid      string

	jwt.StandardClaims
}

func TokenGenrator(email string, password string, uid string) (singnedToken string, signedRefreshToken string, err error) {

	tokenCliam := &SignedDetails{Email: email, Password: password, uid: uid, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()}}

	refreshCalims := &SignedDetails{StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
	}}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenCliam).SignedString([]byte(SCRECT_KEY))

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshCalims).SignedString([]byte(SCRECT_KEY))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SCRECT_KEY), nil
		},
	)

	if err != nil {

		msg = "Invalid token"

		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		// msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		// msg = err.Error()
		return
	}

	return claims, msg
}
