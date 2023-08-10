package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey = "ASDIUFUIS9W7ERBUI89ERB"
)

type Claims struct {
	UID      string
	Username string
	jwt.RegisteredClaims
}

// TODO: Generate token & refreshToken
func GenerateTokens(uid string, username string) (signedToken string, err error) {
	claims := &Claims{
		UID:      uid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour).UTC()),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

	// TODO: generate refreshToken

	// TODO: also return refreshToken
	return signedToken, err
}

func VerifyToken(signedToken string) (uid *string, username *string, err error) {
	var claims = new(Claims)

	token, err := jwt.ParseWithClaims(signedToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, nil, errors.New("invalid authorization header")
	}

	if !token.Valid {
		return nil, nil, errors.New("invalid authorization header")
	}

	// TODO: check for expiry date if required

	return &claims.UID, &claims.Username, nil
}
