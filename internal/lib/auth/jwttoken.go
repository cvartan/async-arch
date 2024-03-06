package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"

	model "async-arch/internal/domain/auth"

	jwt "github.com/golang-jwt/jwt/v5"
)

// JwtTokenChecker - проверка jwt токена
type JwtTokenChecker struct {
	publicKey *rsa.PublicKey
}

func CreateJwtTokenChecker(httpProtocol, authServerAddr, publicKeyHttpMethod, publicKeyMethodPath string) (*JwtTokenChecker, error) {
	client := &http.Client{}

	req, err := http.NewRequest(publicKeyHttpMethod, fmt.Sprintf("%s://%s%s", httpProtocol, authServerAddr, publicKeyMethodPath), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(body)
	if block == nil {
		return nil, errors.New("некоретнный формат ключа")
	}

	pk, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &JwtTokenChecker{
		publicKey: pk,
	}, nil
}

func (c *JwtTokenChecker) Check(tokenString string) (*model.AuthClaims, error) {
	claims := &model.AuthClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) { return c.publicKey, nil })
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token invalid")
	}

	return claims, nil
}
