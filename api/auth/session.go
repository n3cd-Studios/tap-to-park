package auth

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"tap-to-park/database"
	"time"

	"github.com/golang-jwt/jwt"
)

func Generate(user_id uint, device string, host string) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))

	if err != nil {
		return "", err
	}

	expires := time.Now().Add(time.Hour * time.Duration(token_lifespan))
	session := database.Session{
		Device:   device,
		IP:       host,
		Expires:  expires,
		LastUsed: time.Now(),
		UserID:   user_id,
	}
	if err := database.Db.Create(&session).Error; err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["session"] = session.Guid
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

}

func Get(bearerToken string) (string, error) {
	if len(strings.Split(bearerToken, " ")) == 2 {
		token, err := jwt.Parse(strings.Split(bearerToken, " ")[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})
		if err != nil {
			return "", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			return string(fmt.Sprint(claims["session"])), nil
		}
	}
	return "", nil
}
