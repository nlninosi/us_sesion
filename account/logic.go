// Logica de negocio
package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-ldap/ldap/v3"
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
	//number, err := strconv.ParseUint(string(id), 10, 64)
	token, err := jwt.CreateToken(id, username, email)

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
	var PasswordErr = errors.New("contraseña Incorrecta")
	//Accion del Repo
	id, username, upassword, err := s.repository.ValidateUser(ctx, email)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	//Primero contrastamos con el LDAP
	ldapURL := "ldap://host.docker.internal:389/"
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	defer l.Close()
	err = l.Bind("cn=admin,dc=unstream,dc=com", "admin")
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	baseDN := "DC=unstream,DC=com"
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(username))
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})
	result, err := l.Search(searchReq)
	if err != nil {
		level.Error(logger).Log("failed to query LDAP: %w", err)
	}
	fmt.Println("Got", len(result.Entries), "search results")
	userdn := result.Entries[0].DN
	fmt.Println(userdn)
	err = l.Bind(userdn, password)
	if err != nil {
		level.Error(logger).Log("err", err)
	} else {
		fmt.Println("User Authenticaded")
	}
	err = l.Bind("cn=admin,dc=unstream,dc=com", "admin")
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	//Luego con el repo
	if upassword == password {
		token, err := jwt.CreateToken(id, username, email)
		if err != nil {
			level.Error(logger).Log("err", err)
			return "", err
		}
		ok, err := s.repository.UpdateToken(ctx, email, token)
		if err != nil {
			level.Error(logger).Log("err", err)
			return ok, err
		}
		return token, nil
	} else {
		return "", PasswordErr
	}
}

func (s service) ValidateToken(ctx context.Context, email string, token string) (string, error) {
	logger := log.With(s.logger, "method", "ValidateToken")
	logger.Log("token", token)
	check, err := jwt.ExtractTokenMetadata(token)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "badrequest", err
	}
	return check, nil
}
func (s service) NewPassword(ctx context.Context, email string, password string, repassword string) (string, error) {
	logger := log.With(s.logger, "method", "GetUser")
	var PasswordErr = errors.New("las contraseñas ingresadas no son iguales")
	if repassword != password {
		return "", PasswordErr
	}
	ok, err := s.repository.NewPassword(ctx, email, password)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "BadRequest", err
	}
	return ok, nil
}
func (s service) CloseSesion(ctx context.Context, ok string) (string, error) {
	logger := log.With(s.logger, "method", "CloseSesion")

	logger.Log("Closed", ok)
	return " ", nil
}
