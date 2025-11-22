package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrTokenInvalid = errors.Unauthorized("UNAUTHORIZED", "Token is invalid")
	ErrTokenExpired = errors.Unauthorized("UNAUTHORIZED", "Token is expired")
)

type JWT struct {
	secret    []byte
	issuer    string
	expiresIn time.Duration
}

type Claims struct {
	UserId int64  `json:"user_id"`
	Mobile string `json:"mobile"`
	jwt.RegisteredClaims
}

func NewJWT(secret, issuer string, expiresIn int64) *JWT {
	return &JWT{
		secret:    []byte(secret),
		issuer:    issuer,
		expiresIn: time.Duration(expiresIn) * time.Second,
	}
}

// GenerateToken 生成 JWT Token
func (j *JWT) GenerateToken(userId int64, mobile string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(j.expiresIn)

	claims := &Claims{
		UserId: userId,
		Mobile: mobile,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// ParseToken 解析 JWT Token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return j.secret, nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新 Token
func (j *JWT) RefreshToken(tokenString string) (string, int64, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", 0, err
	}

	return j.GenerateToken(claims.UserId, claims.Mobile)
}
