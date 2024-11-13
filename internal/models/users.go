package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var CurrentHashedPassword []byte
	stmt := `SELECT hashed_password FROM users
    WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&CurrentHashedPassword)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(CurrentHashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	stmt1 := `UPDATE users SET hashed_password = $1 WHERE id = $2`
	_, err = m.DB.Exec(stmt1, string(newHashedPassword), id)
	if err != nil {
		return err
	}

	return nil
}
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
             VALUES ($1, $2, $3, NOW())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" && strings.Contains(pgError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = $1`

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
func (m *UserModel) Exists(id int) (bool, error) {
	var Exists bool

	stmt := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"

	err := m.DB.QueryRow(stmt, id).Scan(&Exists)
	return Exists, err
}
func (m *UserModel) Get(id int) (*User, error) {
	stmt := `SELECT name, email, created FROM users
    WHERE id = $1`
	row := m.DB.QueryRow(stmt, id)
	u := &User{}
	err := row.Scan(&u.Name, &u.Email, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}
