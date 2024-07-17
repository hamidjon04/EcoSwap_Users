package token

import (
	"ecoswap/config"
	pb "ecoswap/genproto/users"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(user *pb.GenerateUserJWT) *pb.Token {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	accessClaim := accessToken.Claims.(jwt.MapClaims)
	accessClaim["id"] = user.Id
	accessClaim["username"] = user.Username
	accessClaim["email"] = user.Email
	accessClaim["full_name"] = user.FullName
	accessClaim["iat"] = time.Now().Unix()
	accessClaim["exp"] = time.Now().Add(time.Hour).Unix()

	cfg := config.Load()
	access, err := accessToken.SignedString([]byte(cfg.SIGNING_KEY))
	if err != nil {
		log.Fatalf("Access tokenni generatsiya qilishda xato: %v", err)
	}

	refreshClaim := refreshToken.Claims.(jwt.MapClaims)
	refreshClaim["id"] = user.Id
	refreshClaim["username"] = user.Username
	refreshClaim["email"] = user.Email
	refreshClaim["full_name"] = user.FullName
	refreshClaim["iat"] = time.Now().Unix()
	refreshClaim["exp"] = time.Now().Add(24 * time.Hour).Unix()

	refresh, err := refreshToken.SignedString([]byte(cfg.SIGNING_KEY))
	if err != nil {
		log.Fatalf("Refresh tokenni generatsiya qilishda xato: %v", err)
	}

	return &pb.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func ExtraClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Load().SIGNING_KEY), nil
	})

	if err != nil {
		return nil, fmt.Errorf("faieled to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
