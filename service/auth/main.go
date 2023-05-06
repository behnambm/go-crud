package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
)

type Service struct {
	SecretKey  string
	SignMethod *jwt.SigningMethodHMAC
}

func New(secret string) Service {
	return Service{
		SecretKey:  secret,
		SignMethod: jwt.SigningMethodHS256,
	}
}

func (s Service) GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(s.SignMethod, jwt.MapClaims{
		"uid": userID,
	})
	return token.SignedString([]byte(s.SecretKey))
}

func (s Service) GetClaim(tokenString string) (jwt.MapClaims, error) {
	token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("AUTH SERVICE GetClaim METHOD ERR - Unexpected signing method: %v\n", token.Header["alg"])
		}
		return []byte(s.SecretKey), nil
	})
	if parseErr != nil {
		return nil, fmt.Errorf("AUAT SERVICE GetClaim PARSE ERR %w", parseErr)
	}
	if !token.Valid {
		return nil, fmt.Errorf("AUTH SERVICE GetClaim INVALID TOKEN ERR")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claim type is not valid")
	}
	return claims, nil
}

func (s Service) IsValid(tokenString string) bool {
	_, claimErr := s.GetClaim(tokenString)
	return claimErr == nil
}

func (s Service) IsValidWithClaim(tokenString string) (jwt.MapClaims, bool) {
	claim, claimErr := s.GetClaim(tokenString)
	return claim, claimErr == nil
}
