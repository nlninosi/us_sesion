// Encargado de interactuar con la base de datos
// Este ejemplo esta hecho con postgressql soo..
package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
)

// Error customizado
var ErrRepo = errors.New("unable to handle Repository Request")

type repo struct {
	db     *sql.DB
	logger log.Logger
}

// constructor
func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (repo *repo) CreateUser(ctx context.Context, user User) error {
	// Tabla
	sql := `
		INSERT INTO users (id, email, password, status, username, token)
		VALUES ($1, $2, $3, $4 , $5, $6)`
	// checkeamos si tienen empty strings
	var existemail string
	var existusername string
	err1 := repo.db.QueryRow("SELECT email FROM users WHERE email = $1", user.Email).Scan(&existemail)
	err2 := repo.db.QueryRow("SELECT username FROM users WHERE username = $1", user.UserName).Scan(&existusername)
	var RepoErr_empty = errors.New("existen campos sin rellenar")
	if user.Email == "" || user.Password == "" || user.UserName == "" {
		return RepoErr_empty

	}
	var RepoErr_email = errors.New("ya existe una cuenta con el email ingresado")
	var RepoErr_username = errors.New("ya existe una cuenta con el nombre de usuario ingresado")
	if existemail != "" {
		return RepoErr_email
	}
	if existusername != "" {
		return RepoErr_username
	}
	if err1 == nil && err2 == nil {
		return ErrRepo
	}

	//***************************************************
	// se inserta el usuario
	_, err := repo.db.ExecContext(ctx, sql, user.ID, user.Email, user.Password, user.Status, user.UserName, user.Token)
	if err != nil {
		return err
	}
	return nil
}

// query con el email
func (repo *repo) GetUser(ctx context.Context, id string) (string, string, string, error) {
	var email string
	var status string
	var username string
	err1 := repo.db.QueryRow("SELECT email FROM users WHERE id=$1", id).Scan(&email)
	err2 := repo.db.QueryRow("SELECT status FROM users WHERE id=$1", id).Scan(&status)
	err3 := repo.db.QueryRow("SELECT username FROM users WHERE id=$1", id).Scan(&username)
	if err1 != nil {
		return "", "", "", ErrRepo
	}
	if err2 != nil {
		return "", "", "", ErrRepo
	}
	if err3 != nil {
		return "", "", "", ErrRepo
	}

	return email, status, username, nil
}

func (repo *repo) GetId(ctx context.Context, username string) (string, error) {
	var userid string
	err1 := repo.db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userid)
	if err1 != nil {
		return "", ErrRepo
	}

	return userid, nil
}

func (repo *repo) GetUsers(ctx context.Context) ([5]string, error) {
	var t1Err = errors.New("unable to handle Repository Request at GetUsers")
	//var t2Err = errors.New("unable to handle Repository Request")
	var usernames [5]string
	err := repo.db.QueryRow("SELECT username FROM users order by random() limit 1").Scan(&usernames[0])
	if err != nil {
		return usernames, t1Err
	}
	err1 := repo.db.QueryRow("SELECT username FROM users order by random() limit 1").Scan(&usernames[1])
	if err1 != nil {
		return usernames, t1Err
	}
	err2 := repo.db.QueryRow("SELECT username FROM users order by random() limit 1").Scan(&usernames[2])
	if err2 != nil {
		return usernames, t1Err
	}
	err3 := repo.db.QueryRow("SELECT username FROM users order by random() limit 1").Scan(&usernames[3])
	if err3 != nil {
		return usernames, t1Err
	}
	err4 := repo.db.QueryRow("SELECT username FROM users order by random() limit 1").Scan(&usernames[4])
	if err4 != nil {
		return usernames, t1Err
	}

	return usernames, nil
}

// Aqui obtenemos la contrase;a
// metodo inseguro
// implementar hashes
func (repo *repo) ValidateUser(ctx context.Context, email string) (string, string, string, error) {
	var RepoErr_1 = errors.New("no se encuentra un usuario con el email ingresado")
	//var RepoErr = errors.New("unable to handle Repository Request")
	var userid string
	var password string
	var username string
	err := repo.db.QueryRow("SELECT password FROM users WHERE email=$1", email).Scan(&password)
	if err != nil {
		return "", "", "", RepoErr_1
	}
	err1 := repo.db.QueryRow("SELECT username FROM users WHERE email=$1", email).Scan(&username)
	if err1 != nil {
		return "", "", "", nil
	}
	err2 := repo.db.QueryRow("SELECT id FROM users WHERE email=$1", email).Scan(&userid)
	if err2 != nil {
		return "", "", "", nil
	}

	return userid, username, password, nil
}

func (repo *repo) NewPassword(ctx context.Context, email string, newpassword string) (string, error) {
	var RepoErr_2 = errors.New("no se encuentra un usuario con el email ingresado")
	_, err := repo.db.ExecContext(ctx, "UPDATE users SET password = $1 WHERE email=$2", newpassword, email)
	if err != nil {
		return "", RepoErr_2
	}

	return "password changed", nil
}

func (repo *repo) ValidateToken(ctx context.Context, email string) (string, error) {
	var token string
	err := repo.db.QueryRow("SELECT token FROM users WHERE email=$1", email).Scan(&token)
	if err != nil {
		return "", nil
	}

	return token, nil
}

func (repo *repo) UpdateToken(ctx context.Context, email string, token string) (string, error) {

	_, err := repo.db.ExecContext(ctx, "UPDATE users SET token = $1 WHERE email=$2", token, email)
	if err != nil {
		return "", nil
	}

	return "", nil
}
