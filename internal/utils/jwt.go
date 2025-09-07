package utils

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateToken(secret, email string, ttl time.Duration) (string, time.Time, error) {
    now := time.Now()
    exp := now.Add(ttl)
    claims := TokenClaims{
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(exp),
            IssuedAt:  jwt.NewNumericDate(now),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    s, err := token.SignedString([]byte(secret))
    return s, exp, err
}

func ParseToken(secret, tok string) (string, time.Time, error) {
    parsed, err := jwt.ParseWithClaims(tok, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return "", time.Time{}, err
    }
    claims, ok := parsed.Claims.(*TokenClaims)
    if !ok || !parsed.Valid {
        return "", time.Time{}, errors.New("invalid token claims")
    }
    var exp time.Time
    if claims.ExpiresAt != nil {
        exp = claims.ExpiresAt.Time
    }
    return claims.Email, exp, nil
}

