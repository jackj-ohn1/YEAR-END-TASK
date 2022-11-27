package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

type Claims struct {
	Account  string
	Nickname string
	jwt.StandardClaims
}

var (
	issuer string
	sign   string
)

func init() {
	issuer = viper.GetString("token.issuer")
	sign = viper.GetString("token.sign")
}

func GenerateToken(account string) (string, error) {
	expire := time.Now().Add(time.Hour * 2)
	claims := Claims{
		Account: account,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			Issuer:    issuer,
			Subject:   account,
			IssuedAt:  time.Now().Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(sign))
}

func ParseToken(token string) (*Claims, error) {
	tokenClaim, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaim != nil {
		if claims, ok := tokenClaim.Claims.(*Claims); ok {
			return claims, err
		}
	}
	return nil, err
}
