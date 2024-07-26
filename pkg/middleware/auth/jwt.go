package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// var JwtKey = []byte(os.Getenv("JWT_KEY"))
var JwtKey = []byte("unodesecret")

type JWTClaim struct {
	Wallet string `json:"wallet"`
	jwt.StandardClaims
}

func GenerateJWT(address string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour).UTC()

	claims := &JWTClaim{
		Wallet: address,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(JwtKey)
	if err != nil {
		logrus.Errorln(err)
		return "", err
	}

	return tokenString, nil
}
