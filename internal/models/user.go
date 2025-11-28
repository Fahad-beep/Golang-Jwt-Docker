package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Age      int       `json:"age"`
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) FetchUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT id, email, password_hash FROM users WHERE email = $1"
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *User) error {
	query := `INSERT INTO users (email, password_hash, age) VALUES ($1, $2, $3) RETURNING ID`
	err := r.DB.QueryRow(query, user.Email, user.Password, user.Age).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil

}

func (r *UserRepository) StoreRefreshTokens(userID uuid.UUID, rawToken string, expiryAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, userID, rawToken, expiryAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) ValidateRefreshToken(rawToken string) (*User, error) {
	query := "SELECT u.id, u.email FROM refresh_tokens rt JOIN users u ON rt.user_id = u.id WHERE rt.is_invoked = FALSE AND rt.expires_at > NOW() AND token_hash = $1"
	user := &User{}
	err := r.DB.QueryRow(query, rawToken).Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
