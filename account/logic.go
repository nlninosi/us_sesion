// Logica de negocio
package account

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"

	"us_sesion_ms/account/jwt"
)

type service struct {
	repository Repository
	logger     log.Logger
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

//Estas funciones son un wrapper
//Crea un usuario
func (s service) CreateUser(ctx context.Context, email string, password string, status string, username string) (string, error) {
	logger := log.With(s.logger, "method", "CreateUser")
	// Creacion de id del usuario
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	//forja del token
	token, err := getSignedToken()
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	//
	user := User{
		ID:       id,
		Email:    email,
		Password: password,
		Status:   status,
		UserName: username,
		Token:    token,
	}
	// logger Arroja los errores
	if err := s.repository.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("user created", id)

	return token, nil
}

//Busca un usuario con su correo

func (s service) GetUser(ctx context.Context, id string) (string, string, string, error) {
	logger := log.With(s.logger, "method", "GetUser")

	email, status, username, err := s.repository.GetUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", "", "", err
	}

	logger.Log("Get user", id)

	return email, status, username, nil
}

func (s service) GetId(ctx context.Context, username string) (string, error) {
	logger := log.With(s.logger, "method", "GetId")

	id, err := s.repository.GetId(ctx, username)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Get id", id)

	return id, nil
}

func (s service) GetUsers(ctx context.Context) (string, string, string, string, string, error) {
	logger := log.With(s.logger, "method", "GetUsers")

	usernames, err := s.repository.GetUsers(ctx)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", "", "", "", "", err
	}

	return usernames[0], usernames[1], usernames[2], usernames[3], usernames[4], nil
}

func (s service) ValidateUser(ctx context.Context, email string, password string) (string, error) {
	logger := log.With(s.logger, "method", "ValidateUser")

	upassword, err := s.repository.ValidateUser(ctx, email)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	if upassword == password {
		tokenString, err := getSignedToken()
		if err != nil {
			level.Error(logger).Log("err", err)
			return "", err
		}
		ok, err := s.repository.UpdateToken(ctx, email, tokenString)
		if err != nil {
			level.Error(logger).Log("err", err)
			return ok, err
		}
		return tokenString, nil
	} else {
		return "Incorrect Password", nil
	}
}

func (s service) ValidateToken(ctx context.Context, email string, token string) (string, error) {
	logger := log.With(s.logger, "method", "ValidateToken")
	logger.Log("token", token)
	check, err := jwt.ValidateToken(token, "S0m3_R4n90m_sss")
	if err != nil {
		level.Error(logger).Log("err", err)
		return "badrequest", err
	}
	if check {
		oldToken, err := s.repository.ValidateToken(ctx, email)
		if err != nil {
			level.Error(logger).Log("err", err)
			return "Invalid Token", err
		}
		if token == oldToken {
			newToken, err := getSignedToken()
			if err != nil {
				level.Error(logger).Log("err", err)
				return "not valid token", err
			}
			ok, err := s.repository.UpdateToken(ctx, email, newToken)
			if err != nil {
				level.Error(logger).Log("err", err)
				return ok, err
			}
			return newToken, nil
		} else {
			return "bad Token", nil
		}
	} else {
		return "bad Token", nil
	}
}
func (s service) NewPassword(ctx context.Context, email string, password string, repassword string) (string, error) {
	logger := log.With(s.logger, "method", "GetUser")
	if repassword != password {
		return "Passwords dont match", nil
	}
	ok, err := s.repository.NewPassword(ctx, email, password)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "BadRequest", err
	}

	//logger.Log("Get user", id)

	return ok, nil
}
func (s service) CloseSesion(ctx context.Context, ok string) (string, error) {
	logger := log.With(s.logger, "method", "CloseSesion")

	logger.Log("Closed", ok)
	return " ", nil
}

func getSignedToken() (string, error) {
	// we make a JWT Token here with signing method of ES256 and claims.
	// claims are attributes.
	// aud - audience
	// iss - issuer
	// exp - expiration of the Token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"aud": "frontend.knowsearch.ml",
	// 	"iss": "knowsearch.ml",
	// 	"exp": string(time.Now().Add(time.Minute * 1).Unix()),
	// })
	claimsMap := map[string]string{
		"aud": "frontend.knowsearch.ml",
		"iss": "knowsearch.ml",
		"exp": fmt.Sprint(time.Now().Add(time.Minute * 1).Unix()),
	}
	// here we provide the shared secret. It should be very complex.\
	// Aslo, it should be passed as a System Environment variable

	//	Esta key debe ser secreta
	secret := "S0m3_R4n90m_sss"
	//Implementar como variable de ambiente en docker

	header := "HS256"
	tokenString, err := jwt.GenerateToken(header, claimsMap, secret)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}
