package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type JwtClaim struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateToken(id int, email string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtClaim{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	})

	ss, err := token.SignedString(secret)
	if err != nil {
		return "", errors.Wrap(err, "failed to create auth token")
	}

	return ss, nil
}

// jwt claim token
func ClaimToken(token string) (*JwtClaim, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	parsed, err := jwt.ParseWithClaims(token, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if parsed.Method != jwt.SigningMethodHS256 {
		return nil, errors.New("Invalid token")
	}

	if claim, ok := parsed.Claims.(*JwtClaim); ok {
		return claim, nil
	} else {
		return nil, errors.New("Invalid token")
	}
}
