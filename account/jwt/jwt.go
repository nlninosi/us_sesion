package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AccessDetails struct {
	UserId   string
	UserName string
}

// Function for generating the tokens.
func GenerateToken(header string, payload map[string]string, secret string) (string, error) {
	// create a new hash of type sha256. We pass the secret key to it
	// sha256 is a symmetric cryptographic algorithm
	h := hmac.New(sha256.New, []byte(secret))

	// We base encode the header which is a normal string
	header64 := base64.StdEncoding.EncodeToString([]byte(header))
	// We then Marshal the payload which is a map. This converts it to a string of JSON.
	// Now we base encode this string
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error generating Token")
		return string(payloadstr), err
	}
	payload64 := base64.StdEncoding.EncodeToString(payloadstr)

	// Now add the encoded string.
	message := header64 + "." + payload64

	// We have the unsigned message ready. This is simply concat of header and payload
	unsignedStr := header + string(payloadstr)

	// we write this to the SHA256 to hash it. We can use this to generate the signature now
	h.Write([]byte(unsignedStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//Finally we have the token
	tokenStr := message + "." + signature
	return tokenStr, nil
}

// This helps in validating the token
func ValidateToken(token string, secret string) (bool, error) {
	// JWT has 3 parts separated by '.'
	splitToken := strings.Split(token, ".")
	// if length is not 3, we know that the token is corrupt
	if len(splitToken) != 3 {
		return false, nil
	}

	// decode the header and payload back to strings
	header, err := base64.StdEncoding.DecodeString(splitToken[0])
	if err != nil {
		fmt.Println("1")
		return false, err
	}
	payload, err := base64.StdEncoding.DecodeString(splitToken[1])
	if err != nil {
		fmt.Println("2")
		return false, err
	}

	//again create the signature
	unsignedStr := string(header) + string(payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedStr))

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(signature)

	// if both the signature dont match, this means token is wrong
	if signature != splitToken[2] {
		fmt.Println("3")
		return false, nil
	}
	// This means the token matches
	return true, nil
}

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
