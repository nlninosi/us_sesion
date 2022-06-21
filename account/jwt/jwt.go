package jwt

import (
	"fmt"

	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userid string, username string, email string) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["username"] = username
	atClaims["user_id"] = userid
	atClaims["user_email"] = email
	atClaims["exp"] = time.Now().Add(time.Minute * 240).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func TokenValid(tokenString string) error {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return err
	}
	return nil
}
func ExtractTokenMetadata(tokenString string) (string, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, ok := claims["user_id"].(string)
		if !ok {
			return "", err
		}
		userName, ok := claims["username"].(string)
		if !ok {
			return "", err
		}
		userEmail, ok := claims["user_email"].(string)
		if !ok {
			return "", err
		}
		RefreshToken, err := CreateToken(userId, userName, userEmail)
		if err != nil {
			return "BadRequest", err
		}
		return RefreshToken, nil
	}
	return "", err
}
